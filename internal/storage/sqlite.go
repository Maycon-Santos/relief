package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Maycon-Santos/relief/pkg/fileutil"
	"github.com/Maycon-Santos/relief/pkg/logger"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	conn   *sql.DB
	logger *logger.Logger
}

func NewDB(log *logger.Logger) (*DB, error) {
	dataDir, err := fileutil.GetReliefSubDir("data")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de dados: %w", err)
	}

	dbPath := filepath.Join(dataDir, "orchestrator.db")

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	db := &DB{
		conn:   conn,
		logger: log,
	}

	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("erro ao executar migrations: %w", err)
	}

	log.Info("Banco de dados inicializado", map[string]interface{}{
		"path": dbPath,
	})

	return db, nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) GetConn() *sql.DB {
	return db.conn
}

func (db *DB) migrate() error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de migrations: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		db.logger.Info("Executando migration", map[string]interface{}{
			"file": filename,
		})

		content, err := migrationsFS.ReadFile(filepath.Join("migrations", filename))
		if err != nil {
			return fmt.Errorf("erro ao ler migration %s: %w", filename, err)
		}

		if _, err := db.conn.Exec(string(content)); err != nil {
			return fmt.Errorf("erro ao executar migration %s: %w", filename, err)
		}
	}

	db.logger.Info("Migrations executadas com sucesso", nil)
	return nil
}

func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.conn.Begin()
}

func (db *DB) ClearAllData() error {
	tables := []string{"dependencies", "logs", "projects", "settings"}

	for _, table := range tables {
		if _, err := db.conn.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("erro ao limpar tabela %s: %w", table, err)
		}
	}

	db.logger.Info("Todos os dados foram limpos", nil)
	return nil
}
