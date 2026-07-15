package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// PhotoRepo implémente repository.PhotoRepository.
type PhotoRepo struct {
	pool *pgxpool.Pool
}

func NewPhotoRepo(pool *pgxpool.Pool) *PhotoRepo {
	return &PhotoRepo{pool: pool}
}

func (r *PhotoRepo) Create(ctx context.Context, p *models.Photo) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO photos (soiree_id, chemin) VALUES ($1, $2)
		RETURNING id, created_at`, p.SoireeID, p.Chemin,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *PhotoRepo) ListBySoiree(ctx context.Context, soireeID int64) ([]models.Photo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, soiree_id, chemin, created_at FROM photos WHERE soiree_id = $1 ORDER BY created_at`, soireeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Photo{}
	for rows.Next() {
		var p models.Photo
		if err := rows.Scan(&p.ID, &p.SoireeID, &p.Chemin, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PhotoRepo) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM photos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
