package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// DemandeAmiRepo implémente repository.DemandeAmiRepository (UC21).
type DemandeAmiRepo struct {
	pool *pgxpool.Pool
}

func NewDemandeAmiRepo(pool *pgxpool.Pool) *DemandeAmiRepo {
	return &DemandeAmiRepo{pool: pool}
}

func (r *DemandeAmiRepo) Create(ctx context.Context, d *models.DemandeAmi) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO demandes_amis (demandeur_id, destinataire_id, statut)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		d.DemandeurID, d.DestinataireID, models.DemandeAmiEnAttente,
	).Scan(&d.ID, &d.CreatedAt)
	if isUniqueViolation(err) {
		return repository.ErrConflict
	}
	return err
}

func (r *DemandeAmiRepo) GetByID(ctx context.Context, id int64) (*models.DemandeAmi, error) {
	var d models.DemandeAmi
	err := r.pool.QueryRow(ctx, `
		SELECT id, demandeur_id, destinataire_id, statut, created_at
		FROM demandes_amis WHERE id = $1`, id,
	).Scan(&d.ID, &d.DemandeurID, &d.DestinataireID, &d.Statut, &d.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// ExisteEntre vérifie l'existence d'une demande (dans les deux sens) avec
// l'un des statuts donnés — utilisé pour bloquer un doublon ou une demande
// entre deux utilisateurs déjà amis.
func (r *DemandeAmiRepo) ExisteEntre(ctx context.Context, userAID, userBID int64, statuts []models.StatutDemandeAmi) (bool, error) {
	strStatuts := make([]string, len(statuts))
	for i, s := range statuts {
		strStatuts[i] = string(s)
	}
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM demandes_amis
			WHERE ((demandeur_id = $1 AND destinataire_id = $2)
			    OR (demandeur_id = $2 AND destinataire_id = $1))
			  AND statut = ANY($3)
		)`, userAID, userBID, strStatuts,
	).Scan(&exists)
	return exists, err
}

// ListRecues liste les demandes en attente reçues par un utilisateur.
func (r *DemandeAmiRepo) ListRecues(ctx context.Context, destinataireID int64) ([]models.DemandeAmi, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, demandeur_id, destinataire_id, statut, created_at
		FROM demandes_amis
		WHERE destinataire_id = $1 AND statut = $2
		ORDER BY created_at`, destinataireID, models.DemandeAmiEnAttente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.DemandeAmi{}
	for rows.Next() {
		var d models.DemandeAmi
		if err := rows.Scan(&d.ID, &d.DemandeurID, &d.DestinataireID, &d.Statut, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DemandeAmiRepo) MarkStatut(ctx context.Context, id int64, statut models.StatutDemandeAmi) error {
	ct, err := r.pool.Exec(ctx, `UPDATE demandes_amis SET statut = $1 WHERE id = $2`, statut, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ListAmis liste les utilisateurs devenus amis (demande acceptée, peu
// importe qui l'a envoyée).
func (r *DemandeAmiRepo) ListAmis(ctx context.Context, userID int64) ([]models.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.pseudo, u.email, u.password_hash, u.avatar, u.bio, u.score, u.role, u.created_at
		FROM users u
		JOIN demandes_amis d ON (
			(d.demandeur_id = $1 AND d.destinataire_id = u.id) OR
			(d.destinataire_id = $1 AND d.demandeur_id = u.id)
		)
		WHERE d.statut = $2
		ORDER BY u.pseudo`, userID, models.DemandeAmiAcceptee)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Pseudo, &u.Email, &u.PasswordHash, &u.Avatar, &u.Bio, &u.Score, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
