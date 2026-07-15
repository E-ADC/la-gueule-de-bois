package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// VoteRepo implémente repository.VoteRepository (UC12).
type VoteRepo struct {
	pool *pgxpool.Pool
}

func NewVoteRepo(pool *pgxpool.Pool) *VoteRepo {
	return &VoteRepo{pool: pool}
}

func (r *VoteRepo) Create(ctx context.Context, v *models.Vote) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO votes (temoignage_id, user_id, valeur) VALUES ($1, $2, $3)
		RETURNING id, created_at`, v.TemoignageID, v.UserID, v.Valeur,
	).Scan(&v.ID, &v.CreatedAt)
	if isUniqueViolation(err) {
		return repository.ErrConflict
	}
	return err
}

func (r *VoteRepo) Exists(ctx context.Context, temoignageID, userID int64) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM votes WHERE temoignage_id = $1 AND user_id = $2)`,
		temoignageID, userID,
	).Scan(&exists)
	return exists, err
}

// SommeVotesForOwner agrège les votes positifs/négatifs reçus sur tous les
// témoignages des soirées d'un utilisateur (calcul du score, UC16).
func (r *VoteRepo) SommeVotesForOwner(ctx context.Context, ownerID int64) (positifs int, negatifs int, err error) {
	err = r.pool.QueryRow(ctx, `
		SELECT
			count(*) FILTER (WHERE v.valeur = 1),
			count(*) FILTER (WHERE v.valeur = -1)
		FROM votes v
		JOIN temoignages t ON t.id = v.temoignage_id
		JOIN soirees s ON s.id = t.soiree_id
		WHERE s.user_id = $1`, ownerID,
	).Scan(&positifs, &negatifs)
	return positifs, negatifs, err
}
