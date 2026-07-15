package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// TemoignageRepo implémente repository.TemoignageRepository (UC11).
type TemoignageRepo struct {
	pool *pgxpool.Pool
}

func NewTemoignageRepo(pool *pgxpool.Pool) *TemoignageRepo {
	return &TemoignageRepo{pool: pool}
}

func (r *TemoignageRepo) Create(ctx context.Context, t *models.Temoignage) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO temoignages (soiree_id, auteur_id, contenu) VALUES ($1, $2, $3)
		RETURNING id, created_at`, t.SoireeID, t.AuteurID, t.Contenu,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TemoignageRepo) GetByID(ctx context.Context, id int64) (*models.Temoignage, error) {
	var t models.Temoignage
	err := r.pool.QueryRow(ctx, `
		SELECT id, soiree_id, auteur_id, contenu, created_at FROM temoignages WHERE id = $1`, id,
	).Scan(&t.ID, &t.SoireeID, &t.AuteurID, &t.Contenu, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemoignageRepo) ListBySoiree(ctx context.Context, soireeID int64) ([]models.Temoignage, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, soiree_id, auteur_id, contenu, created_at FROM temoignages
		WHERE soiree_id = $1 ORDER BY created_at`, soireeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Temoignage{}
	for rows.Next() {
		var t models.Temoignage
		if err := rows.Scan(&t.ID, &t.SoireeID, &t.AuteurID, &t.Contenu, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// CountForOwner compte les témoignages reçus sur les soirées d'un
// utilisateur (utilisé par le calcul du score, UC16).
func (r *TemoignageRepo) CountForOwner(ctx context.Context, ownerID int64) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `
		SELECT count(*) FROM temoignages t
		JOIN soirees s ON s.id = t.soiree_id
		WHERE s.user_id = $1`, ownerID).Scan(&n)
	return n, err
}

func (r *TemoignageRepo) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM temoignages WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
