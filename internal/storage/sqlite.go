// Package storage fornece a camada de persistência da aplicação.
package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"github.com/relief-org/relief/pkg/fileutil"
	"github.com/relief-org/relief/pkg/logger"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB é o wrapper para o banco de dados SQLite
type DB struct {
	conn   *sql.DB
	logger *logger.Logger
}

// NewDB cria uma nova conexão com o banco de dados
func NewDB(log *logger.Logger) (*DB, error) {
	// Database path: ~/.relief/data/orchestrator.db
	dataDir, err := fileutil.GetReliefSubDir("data")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de dados: %w", err)
	}

	dbPath := filepath.Join(dataDir, "orchestrator.db")

	// Abrir conexão
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}

	// Configurar pool de conexões
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)

	// Testar conexão
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	db := &DB{
		conn:   conn,
		logger: log,
	}

	// Executar migrations
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("erro ao executar migrations: %w", err)
	}

	log.Info("Banco de dados inicializado", map[string]interface{}{
		"path": dbPath,
	})

	return db, nil
}

// Close fecha a conexão com o banco de dados
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// GetConn retorna a conexão subjacente
func (db *DB) GetConn() *sql.DB {
	return db.conn
}

// migrate executa as migrations do banco de dados
func (db *DB) migrate() error {
	// Ler todos os arquivos de migration do diretório embedado
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de migrations: %w", err)
	}

	// Ordenar migrations por nome (001, 002, etc.)
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Executar cada migration em ordem
	for _, filename := range migrationFiles {
		db.logger.Info("Executando migration", map[string]interface{}{
			"file": filename,
		})

		// Ler conteúdo do arquivo
		content, err := migrationsFS.ReadFile(filepath.Join("migrations", filename))
		if err != nil {
			return fmt.Errorf("erro ao ler migration %s: %w", filename, err)
		}

		// Executar SQL
		if _, err := db.conn.Exec(string(content)); err != nil {
			return fmt.Errorf("erro ao executar migration %s: %w", filename, err)
		}
	}

	db.logger.Info("Migrations executadas com sucesso", nil)
	return nil
}

// BeginTx inicia uma transação
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.conn.Begin()
}

// ClearAllData limpa todos os dados do banco (útil para testes)
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
