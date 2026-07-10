// Package repository définit les interfaces d'accès aux données. Les
// implémentations Postgres vivent dans le sous-paquet postgres ; les
// services utilisent ces interfaces pour rester testables sans base
// vivante (mocks écrits à la main dans les tests table-driven).
package repository

import (
	"context"
	"errors"
	"time"

	"gueuledebois/backend/internal/models"
)

// Erreurs sentinelles génériques, à mapper par les handlers vers les codes
// HTTP demandés par la spec (400/403/404/409).
var (
	ErrNotFound = errors.New("repository: ressource introuvable")
	ErrConflict = errors.New("repository: conflit (doublon)")
)

// UserRepository gère la persistance des utilisateurs (UC01/02/04/05/17).
type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPseudo(ctx context.Context, pseudo string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	UpdateScore(ctx context.Context, userID int64, score int) error
	ListLeaderboard(ctx context.Context, limit int) ([]models.User, error)
	// ListLeaderboardForGroup restreint le classement aux membres d'un
	// groupe (UC20, généralisation de UC17).
	ListLeaderboardForGroup(ctx context.Context, groupeID int64, limit int) ([]models.User, error)
}

// SessionRepository gère les sessions cookie HttpOnly (auth sans JWT).
type SessionRepository interface {
	Create(ctx context.Context, s *models.Session) error
	GetByToken(ctx context.Context, token string) (*models.Session, error)
	DeleteByToken(ctx context.Context, token string) error
}

// SoireeRepository gère les soirées (UC06/07/08/10).
type SoireeRepository interface {
	Create(ctx context.Context, s *models.Soiree) error
	GetByID(ctx context.Context, id int64) (*models.Soiree, error)
	Update(ctx context.Context, s *models.Soiree) error
	Delete(ctx context.Context, id int64) error
	ListByUser(ctx context.Context, userID int64) ([]models.Soiree, error)
	CountByUser(ctx context.Context, userID int64) (int, error)
}

// PhotoRepository gère les photos attachées à une soirée (spec "photos minimales").
type PhotoRepository interface {
	Create(ctx context.Context, p *models.Photo) error
	ListBySoiree(ctx context.Context, soireeID int64) ([]models.Photo, error)
	Delete(ctx context.Context, id int64) error
}

// TemoinInvitationRepository gère les invitations de témoins (UC09), et sert
// à vérifier la pré-condition de UC11.
type TemoinInvitationRepository interface {
	Create(ctx context.Context, inv *models.TemoinInvitation) error
	IsInvited(ctx context.Context, soireeID, userID int64) (bool, error)
}

// TemoignageRepository gère les témoignages (UC11).
type TemoignageRepository interface {
	Create(ctx context.Context, t *models.Temoignage) error
	GetByID(ctx context.Context, id int64) (*models.Temoignage, error)
	ListBySoiree(ctx context.Context, soireeID int64) ([]models.Temoignage, error)
	// CountForOwner et SommeVotesForOwner alimentent le calcul du score
	// (UC16) : nombre de témoignages reçus par les soirées d'un
	// utilisateur, et somme des votes reçus dessus.
	CountForOwner(ctx context.Context, ownerID int64) (int, error)
	Delete(ctx context.Context, id int64) error
}

// VoteRepository gère les votes/swipes sur témoignages (UC12).
type VoteRepository interface {
	Create(ctx context.Context, v *models.Vote) error
	Exists(ctx context.Context, temoignageID, userID int64) (bool, error)
	SommeVotesForOwner(ctx context.Context, ownerID int64) (positifs int, negatifs int, err error)
}

// BadgeRepository gère le catalogue de badges et leur attribution (UC14/UC15).
type BadgeRepository interface {
	ListAll(ctx context.Context) ([]models.Badge, error)
	ListForUser(ctx context.Context, userID int64) ([]models.Badge, error)
	AttachToUser(ctx context.Context, userID, badgeID int64) error
}

// GroupeRepository gère les groupes d'amis (UC18/19/20).
type GroupeRepository interface {
	Create(ctx context.Context, g *models.Groupe) error
	GetByID(ctx context.Context, id int64) (*models.Groupe, error)
	GetByNom(ctx context.Context, nom string) (*models.Groupe, error)
	AddMember(ctx context.Context, groupeID, userID int64) error
	IsMember(ctx context.Context, groupeID, userID int64) (bool, error)
}

// SignalementRepository gère les signalements de témoignages (UC13/UC22).
type SignalementRepository interface {
	Create(ctx context.Context, s *models.Signalement) error
	Exists(ctx context.Context, temoignageID, auteurID int64) (bool, error)
	GetByID(ctx context.Context, id int64) (*models.Signalement, error)
	ListEnAttente(ctx context.Context) ([]models.Signalement, error)
	MarkTraite(ctx context.Context, id int64, statut models.StatutSignalement, moderateurID int64, traiteLe time.Time) error
}
