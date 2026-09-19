package repository

import (
	"JobTracker/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositoy interface {
	Create(ctx context.Context, cargo, empresa, status string) (models.Jobs, error)
	GetById(ctx context.Context, id int) (models.Jobs, error)
	List(ctx context.Context) ([]models.Jobs, error)
	Update(ctx context.Context, id int, cargo, empresa, status string) error
	DeleteById(ctx context.Context, id int) error
}

var ErrNotFound = errors.New("Job not Found")

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repositoy {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, cargo, empresa, status string) (models.Jobs, error) {
	const q = `INSERT INTO jobs (cargo, empresa, status) VALUES ($1,$2,$3) RETURNING id, cargo, empresa, status, created_at`

	var j models.Jobs
	err := r.db.QueryRow(ctx, q, cargo, empresa, status).
		Scan(&j.Id, &j.Cargo, &j.Empresa, &j.Status, &j.Data)
	return j, err
}

func (r *postgresRepository) GetById(ctx context.Context, id int) (models.Jobs, error) {
	const q = `SELECT id, cargo, empresa, status, created_at FROM jobs WHERE id = $1`

	var j models.Jobs
	err := r.db.QueryRow(ctx, q, id).Scan(&j.Id, &j.Cargo, &j.Empresa, &j.Status, &j.Data)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Jobs{}, ErrNotFound
	}
	return j, err
}

func (r *postgresRepository) List(ctx context.Context) ([]models.Jobs, error) {
	const q = `SELECT id, cargo, empresa, status, created_at FROM jobs ORDER BY id`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Jobs
	for rows.Next() {
		var j models.Jobs
		if err := rows.Scan(&j.Id, &j.Cargo, &j.Empresa, &j.Status, &j.Data); err != nil {
			return nil, err
		}
		out = append(out, j)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error")
	}
	return out, rows.Err()

}

func (r *postgresRepository) Update(ctx context.Context, id int, cargo, empresa, status string) error {
	tag, err := r.db.Exec(ctx, `UPDATE jobs 
SET cargo = $2, empresa = $3, status = $4 WHERE id = $1 `, id, cargo, empresa, status)

	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepository) DeleteById(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM jobs WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil

}
