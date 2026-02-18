package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Maycon-Santos/relief/internal/domain"
)

type ProjectRepository struct {
	db *DB
}

func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *domain.Project) error {
	query := `
		INSERT INTO projects (id, name, path, domain, type, status, port, pid, last_error, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.conn.Exec(query,
		project.ID,
		project.Name,
		project.Path,
		project.Domain,
		project.Type,
		project.Status,
		project.Port,
		project.PID,
		project.LastError,
		project.CreatedAt,
		project.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar projeto: %w", err)
	}

	if err := r.saveDependencies(project); err != nil {
		return fmt.Errorf("erro ao salvar dependências: %w", err)
	}

	return nil
}

func (r *ProjectRepository) Update(project *domain.Project) error {
	query := `
		UPDATE projects 
		SET name = ?, path = ?, domain = ?, type = ?, status = ?, port = ?, 
		    pid = ?, last_error = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.conn.Exec(query,
		project.Name,
		project.Path,
		project.Domain,
		project.Type,
		project.Status,
		project.Port,
		project.PID,
		project.LastError,
		time.Now().Format(time.RFC3339),
		project.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar projeto: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("projeto não encontrado")
	}

	if err := r.saveDependencies(project); err != nil {
		return fmt.Errorf("erro ao atualizar dependências: %w", err)
	}

	return nil
}

func (r *ProjectRepository) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := r.db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar projeto: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("projeto não encontrado")
	}

	return nil
}

func (r *ProjectRepository) GetByID(id string) (*domain.Project, error) {
	query := `
		SELECT id, name, path, domain, type, status, port, pid, last_error, created_at, updated_at
		FROM projects WHERE id = ?
	`

	var project domain.Project
	err := r.db.conn.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Path,
		&project.Domain,
		&project.Type,
		&project.Status,
		&project.Port,
		&project.PID,
		&project.LastError,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("projeto não encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar projeto: %w", err)
	}

	if project.Scripts == nil {
		project.Scripts = make(map[string]string)
	}
	if project.Env == nil {
		project.Env = make(map[string]string)
	}
	if project.Dependencies == nil {
		project.Dependencies = []domain.Dependency{}
	}

	if err := r.loadDependencies(&project); err != nil {
		return nil, fmt.Errorf("erro ao carregar dependências: %w", err)
	}

	if manifest, err := domain.ParseManifest(project.Path); err == nil {
		project.Manifest = manifest
		if portStr, ok := manifest.Env["PORT"]; ok {
			if port, err := strconv.Atoi(portStr); err == nil && project.Port == 0 {
				project.Port = port
			}
		} else if manifest.Ports != nil {
			if mainPort, ok := manifest.Ports["main"]; ok && project.Port == 0 {
				project.Port = mainPort
			}
		}
	}

	return &project, nil
}

func (r *ProjectRepository) GetByName(name string) (*domain.Project, error) {
	query := `
		SELECT id, name, path, domain, type, status, port, pid, last_error, created_at, updated_at
		FROM projects WHERE name = ?
	`

	var project domain.Project
	err := r.db.conn.QueryRow(query, name).Scan(
		&project.ID,
		&project.Name,
		&project.Path,
		&project.Domain,
		&project.Type,
		&project.Status,
		&project.Port,
		&project.PID,
		&project.LastError,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("projeto não encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar projeto: %w", err)
	}

	if project.Scripts == nil {
		project.Scripts = make(map[string]string)
	}
	if project.Env == nil {
		project.Env = make(map[string]string)
	}
	if project.Dependencies == nil {
		project.Dependencies = []domain.Dependency{}
	}

	if err := r.loadDependencies(&project); err != nil {
		return nil, fmt.Errorf("erro ao carregar dependências: %w", err)
	}

	if manifest, err := domain.ParseManifest(project.Path); err == nil {
		project.Manifest = manifest
		if portStr, ok := manifest.Env["PORT"]; ok {
			if port, err := strconv.Atoi(portStr); err == nil && project.Port == 0 {
				project.Port = port
			}
		} else if manifest.Ports != nil {
			if mainPort, ok := manifest.Ports["main"]; ok && project.Port == 0 {
				project.Port = mainPort
			}
		}
	}

	return &project, nil
}

func (r *ProjectRepository) List() ([]*domain.Project, error) {
	query := `
		SELECT id, name, path, domain, type, status, port, pid, last_error, created_at, updated_at
		FROM projects
		ORDER BY name ASC
	`

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar projetos: %w", err)
	}
	defer rows.Close()

	projects := []*domain.Project{}
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Path,
			&project.Domain,
			&project.Type,
			&project.Status,
			&project.Port,
			&project.PID,
			&project.LastError,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear projeto: %w", err)
		}

		if project.Scripts == nil {
			project.Scripts = make(map[string]string)
		}
		if project.Env == nil {
			project.Env = make(map[string]string)
		}
		if project.Dependencies == nil {
			project.Dependencies = []domain.Dependency{}
		}

		if err := r.loadDependencies(&project); err != nil {
			return nil, fmt.Errorf("erro ao carregar dependências: %w", err)
		}

		if manifest, err := domain.ParseManifest(project.Path); err == nil {
			project.Manifest = manifest
			if portStr, ok := manifest.Env["PORT"]; ok {
				if port, err := strconv.Atoi(portStr); err == nil && project.Port == 0 {
					project.Port = port
				}
			} else if manifest.Ports != nil {
				if mainPort, ok := manifest.Ports["main"]; ok && project.Port == 0 {
					project.Port = mainPort
				}
			}
		}

		projects = append(projects, &project)
	}

	return projects, nil
}

func (r *ProjectRepository) saveDependencies(project *domain.Project) error {
	deleteQuery := `DELETE FROM dependencies WHERE project_id = ?`
	if _, err := r.db.conn.Exec(deleteQuery, project.ID); err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO dependencies (project_id, name, version, required_version, managed, satisfied, message)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	for _, dep := range project.Dependencies {
		_, err := r.db.conn.Exec(insertQuery,
			project.ID,
			dep.Name,
			dep.Version,
			dep.RequiredVersion,
			dep.Managed,
			dep.Satisfied,
			dep.Message,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectRepository) loadDependencies(project *domain.Project) error {
	query := `
		SELECT name, version, required_version, managed, satisfied, message
		FROM dependencies WHERE project_id = ?
	`

	rows, err := r.db.conn.Query(query, project.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	project.Dependencies = []domain.Dependency{}
	for rows.Next() {
		var dep domain.Dependency
		err := rows.Scan(
			&dep.Name,
			&dep.Version,
			&dep.RequiredVersion,
			&dep.Managed,
			&dep.Satisfied,
			&dep.Message,
		)
		if err != nil {
			return err
		}
		project.Dependencies = append(project.Dependencies, dep)
	}

	return nil
}

type LogRepository struct {
	db *DB
}

func NewLogRepository(db *DB) *LogRepository {
	return &LogRepository{db: db}
}

func (r *LogRepository) Create(log *domain.LogEntry) error {
	query := `INSERT INTO logs (project_id, level, message, timestamp) VALUES (?, ?, ?, ?)`

	result, err := r.db.conn.Exec(query, log.ProjectID, log.Level, log.Message, log.Timestamp)
	if err != nil {
		return fmt.Errorf("erro ao criar log: %w", err)
	}

	id, _ := result.LastInsertId()
	log.ID = id

	return nil
}

func (r *LogRepository) GetByProjectID(projectID string, limit int) ([]domain.LogEntry, error) {
	query := `
		SELECT id, project_id, level, message, timestamp
		FROM logs WHERE project_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.conn.Query(query, projectID, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar logs: %w", err)
	}
	defer rows.Close()

	logs := []domain.LogEntry{}
	for rows.Next() {
		var log domain.LogEntry
		err := rows.Scan(&log.ID, &log.ProjectID, &log.Level, &log.Message, &log.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	for i := len(logs)/2 - 1; i >= 0; i-- {
		opp := len(logs) - 1 - i
		logs[i], logs[opp] = logs[opp], logs[i]
	}

	return logs, nil
}

func (r *LogRepository) DeleteOldLogs(olderThan time.Time) error {
	query := `DELETE FROM logs WHERE timestamp < ?`
	_, err := r.db.conn.Exec(query, olderThan)
	return err
}
