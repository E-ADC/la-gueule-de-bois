package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// GroupeRepo implémente repository.GroupeRepository (UC18/19/20).
type GroupeRepo struct {
	pool *pgxpool.Pool
}

func NewGroupeRepo(pool *pgxpool.Pool) *GroupeRepo {
	return &GroupeRepo{pool: pool}
}

func (r *GroupeRepo) Create(ctx context.Context, g *models.Groupe) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO groupes (nom, createur_id) VALUES ($1, $2)
		RETURNING id, created_at`, g.Nom, g.CreateurID,
	).Scan(&g.ID, &g.CreatedAt)
	if isUniqueViolation(err) {
		return repository.ErrConflict
	}
	return err
}

func (r *GroupeRepo) GetByID(ctx context.Context, id int64) (*models.Groupe, error) {
	var g models.Groupe
	err := r.pool.QueryRow(ctx, `
		SELECT id, nom, createur_id, created_at FROM groupes WHERE id = $1`, id,
	).Scan(&g.ID, &g.Nom, &g.CreateurID, &g.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GroupeRepo) GetByNom(ctx context.Context, nom string) (*models.Groupe, error) {
	var g models.Groupe
	err := r.pool.QueryRow(ctx, `
		SELECT id, nom, createur_id, created_at FROM groupes WHERE nom = $1`, nom,
	).Scan(&g.ID, &g.Nom, &g.CreateurID, &g.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GroupeRepo) AddMember(ctx context.Context, groupeID, userID int64) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO groupe_membres (groupe_id, user_id) VALUES ($1, $2)
		ON CONFLICT (groupe_id, user_id) DO NOTHING`, groupeID, userID)
	return err
}

func (r *GroupeRepo) IsMember(ctx context.Context, groupeID, userID int64) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM groupe_membres WHERE groupe_id = $1 AND user_id = $2)`,
		groupeID, userID,
	).Scan(&exists)
	return exists, err
}

// ListForUser liste les groupes dont l'utilisateur est membre.
func (r *GroupeRepo) ListForUser(ctx context.Context, userID int64) ([]models.Groupe, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.nom, g.createur_id, g.created_at
		FROM groupes g
		JOIN groupe_membres m ON m.groupe_id = g.id
		WHERE m.user_id = $1
		ORDER BY g.nom`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Groupe{}
	for rows.Next() {
		var g models.Groupe
		if err := rows.Scan(&g.ID, &g.Nom, &g.CreateurID, &g.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
