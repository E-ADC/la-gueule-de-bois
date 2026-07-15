package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// SoireeRepo implémente repository.SoireeRepository.
type SoireeRepo struct {
	pool *pgxpool.Pool
}

func NewSoireeRepo(pool *pgxpool.Pool) *SoireeRepo {
	return &SoireeRepo{pool: pool}
}

func (r *SoireeRepo) Create(ctx context.Context, s *models.Soiree) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO soirees (user_id, titre, date_soiree, lieu, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		s.UserID, s.Titre, s.DateSoiree, s.Lieu, s.Description,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SoireeRepo) GetByID(ctx context.Context, id int64) (*models.Soiree, error) {
	var s models.Soiree
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, titre, date_soiree, lieu, description, created_at, updated_at
		FROM soirees WHERE id = $1`, id,
	).Scan(&s.ID, &s.UserID, &s.Titre, &s.DateSoiree, &s.Lieu, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SoireeRepo) Update(ctx context.Context, s *models.Soiree) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE soirees SET titre = $1, date_soiree = $2, lieu = $3, description = $4, updated_at = now()
		WHERE id = $5`,
		s.Titre, s.DateSoiree, s.Lieu, s.Description, s.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *SoireeRepo) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM soirees WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *SoireeRepo) ListByUser(ctx context.Context, userID int64) ([]models.Soiree, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, titre, date_soiree, lieu, description, created_at, updated_at
		FROM soirees WHERE user_id = $1 ORDER BY date_soiree DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Soiree{}
	for rows.Next() {
		var s models.Soiree
		if err := rows.Scan(&s.ID, &s.UserID, &s.Titre, &s.DateSoiree, &s.Lieu, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SoireeRepo) CountByUser(ctx context.Context, userID int64) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT count(*) FROM soirees WHERE user_id = $1`, userID).Scan(&n)
	return n, err
}
