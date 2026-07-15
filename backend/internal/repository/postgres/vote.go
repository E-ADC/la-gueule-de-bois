package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

// CountForTemoignage compte les votes positifs/négatifs reçus par un
// témoignage donné (affichage enrichi, cf. SoireeDetailPage).
func (r *VoteRepo) CountForTemoignage(ctx context.Context, temoignageID int64) (positifs int, negatifs int, err error) {
	err = r.pool.QueryRow(ctx, `
		SELECT count(*) FILTER (WHERE valeur = 1), count(*) FILTER (WHERE valeur = -1)
		FROM votes WHERE temoignage_id = $1`, temoignageID,
	).Scan(&positifs, &negatifs)
	return positifs, negatifs, err
}

// GetByUserAndTemoignage renvoie le vote de cet utilisateur sur ce
// témoignage, ou repository.ErrNotFound s'il n'a pas encore voté.
func (r *VoteRepo) GetByUserAndTemoignage(ctx context.Context, userID, temoignageID int64) (*models.Vote, error) {
	var v models.Vote
	err := r.pool.QueryRow(ctx, `
		SELECT id, temoignage_id, user_id, valeur, created_at
		FROM votes WHERE user_id = $1 AND temoignage_id = $2`, userID, temoignageID,
	).Scan(&v.ID, &v.TemoignageID, &v.UserID, &v.Valeur, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}
