package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// TemoinInvitationRepo implémente repository.TemoinInvitationRepository (UC09).
type TemoinInvitationRepo struct {
	pool *pgxpool.Pool
}

func NewTemoinInvitationRepo(pool *pgxpool.Pool) *TemoinInvitationRepo {
	return &TemoinInvitationRepo{pool: pool}
}

func (r *TemoinInvitationRepo) Create(ctx context.Context, inv *models.TemoinInvitation) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO temoin_invitations (soiree_id, invite_id) VALUES ($1, $2)
		RETURNING id, created_at`, inv.SoireeID, inv.InviteID,
	).Scan(&inv.ID, &inv.CreatedAt)
	if isUniqueViolation(err) {
		return repository.ErrConflict
	}
	return err
}

func (r *TemoinInvitationRepo) IsInvited(ctx context.Context, soireeID, userID int64) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM temoin_invitations WHERE soiree_id = $1 AND invite_id = $2)`,
		soireeID, userID,
	).Scan(&exists)
	return exists, err
}
