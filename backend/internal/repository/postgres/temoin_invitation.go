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

// ListSoireesForInvite liste les soirées où l'utilisateur a été invité
// comme témoin — sans quoi il n'a aucun moyen de retrouver la soirée
// dans l'app une fois invité.
func (r *TemoinInvitationRepo) ListSoireesForInvite(ctx context.Context, userID int64) ([]models.Soiree, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.user_id, s.titre, s.date_soiree, s.lieu, s.description, s.created_at, s.updated_at
		FROM soirees s
		JOIN temoin_invitations i ON i.soiree_id = s.id
		WHERE i.invite_id = $1
		ORDER BY s.date_soiree DESC`, userID)
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
