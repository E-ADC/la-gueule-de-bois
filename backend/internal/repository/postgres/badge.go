package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
)

// BadgeRepo implémente repository.BadgeRepository (UC14/UC15).
type BadgeRepo struct {
	pool *pgxpool.Pool
}

func NewBadgeRepo(pool *pgxpool.Pool) *BadgeRepo {
	return &BadgeRepo{pool: pool}
}

func (r *BadgeRepo) ListAll(ctx context.Context) ([]models.Badge, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, code, nom, description, seuil_score FROM badges ORDER BY seuil_score`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Badge
	for rows.Next() {
		var b models.Badge
		if err := rows.Scan(&b.ID, &b.Code, &b.Nom, &b.Description, &b.SeuilScore); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BadgeRepo) ListForUser(ctx context.Context, userID int64) ([]models.Badge, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.code, b.nom, b.description, b.seuil_score
		FROM badges b
		JOIN user_badges ub ON ub.badge_id = b.id
		WHERE ub.user_id = $1
		ORDER BY b.seuil_score`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Badge
	for rows.Next() {
		var b models.Badge
		if err := rows.Scan(&b.ID, &b.Code, &b.Nom, &b.Description, &b.SeuilScore); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BadgeRepo) AttachToUser(ctx context.Context, userID, badgeID int64) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_badges (user_id, badge_id) VALUES ($1, $2)
		ON CONFLICT (user_id, badge_id) DO NOTHING`, userID, badgeID)
	return err
}
