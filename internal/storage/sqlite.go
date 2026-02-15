// Package storage fornece a camada de persistência da aplicação.
package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/omelete/relief/pkg/fileutil"
	"github.com/omelete/relief/pkg/logger"
)

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
	// Ler arquivo de migration
	migrationPath := "internal/storage/migrations/001_initial.sql"
	
	// Como estamos em um executável compilado, vamos embutir o SQL diretamente
	migrationSQL := `
		CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			path TEXT NOT NULL,
			domain TEXT,
			type TEXT NOT NULL,
			status TEXT NOT NULL,
			port INTEGER,
			pid INTEGER,
			last_error TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
		CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
		CREATE INDEX IF NOT EXISTS idx_projects_domain ON projects(domain);

		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT NOT NULL,
			level TEXT NOT NULL,
			message TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_logs_project_id ON logs(project_id);
		CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
		CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);

		CREATE TABLE IF NOT EXISTS dependencies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			version TEXT,
			required_version TEXT NOT NULL,
			managed BOOLEAN NOT NULL DEFAULT 0,
			satisfied BOOLEAN NOT NULL DEFAULT 0,
			message TEXT,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_dependencies_project_id ON dependencies(project_id);
		CREATE INDEX IF NOT EXISTS idx_dependencies_name ON dependencies(name);

		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME NOT NULL
		);
	`

	// Executar migration
	if _, err := db.conn.Exec(migrationSQL); err != nil {
		return fmt.Errorf("erro ao executar migration: %w", err)
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
