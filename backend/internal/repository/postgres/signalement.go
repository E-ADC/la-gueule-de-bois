package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// SignalementRepo implémente repository.SignalementRepository (UC13/UC22).
type SignalementRepo struct {
	pool *pgxpool.Pool
}

func NewSignalementRepo(pool *pgxpool.Pool) *SignalementRepo {
	return &SignalementRepo{pool: pool}
}

func (r *SignalementRepo) Create(ctx context.Context, s *models.Signalement) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO signalements (temoignage_id, auteur_id, motif, statut)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		s.TemoignageID, s.AuteurID, s.Motif, models.SignalementEnAttente,
	).Scan(&s.ID, &s.CreatedAt)
	if isUniqueViolation(err) {
		return repository.ErrConflict
	}
	return err
}

func (r *SignalementRepo) Exists(ctx context.Context, temoignageID, auteurID int64) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM signalements WHERE temoignage_id = $1 AND auteur_id = $2)`,
		temoignageID, auteurID,
	).Scan(&exists)
	return exists, err
}

func (r *SignalementRepo) GetByID(ctx context.Context, id int64) (*models.Signalement, error) {
	var s models.Signalement
	err := r.pool.QueryRow(ctx, `
		SELECT id, temoignage_id, auteur_id, motif, statut, traite_par_id, created_at, traite_le
		FROM signalements WHERE id = $1`, id,
	).Scan(&s.ID, &s.TemoignageID, &s.AuteurID, &s.Motif, &s.Statut, &s.TraiteParID, &s.CreatedAt, &s.TraiteLe)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SignalementRepo) ListEnAttente(ctx context.Context) ([]models.Signalement, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, temoignage_id, auteur_id, motif, statut, traite_par_id, created_at, traite_le
		FROM signalements WHERE statut = $1 ORDER BY created_at`, models.SignalementEnAttente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Signalement
	for rows.Next() {
		var s models.Signalement
		if err := rows.Scan(&s.ID, &s.TemoignageID, &s.AuteurID, &s.Motif, &s.Statut, &s.TraiteParID, &s.CreatedAt, &s.TraiteLe); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SignalementRepo) MarkTraite(ctx context.Context, id int64, statut models.StatutSignalement, moderateurID int64, traiteLe time.Time) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE signalements SET statut = $1, traite_par_id = $2, traite_le = $3
		WHERE id = $4`, statut, moderateurID, traiteLe, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
