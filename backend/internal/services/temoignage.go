package services

import (
	"context"
	"strings"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// TemoignageService implémente UC11 (ajouter un témoignage) et UC09
// (inviter un témoin, pré-condition de UC11).
type TemoignageService struct {
	temoignages repository.TemoignageRepository
	invitations repository.TemoinInvitationRepository
	soirees     repository.SoireeRepository
	users       repository.UserRepository
	scoring     *ScoringService
	notifier    Notifier
}

func NewTemoignageService(
	temoignages repository.TemoignageRepository,
	invitations repository.TemoinInvitationRepository,
	soirees repository.SoireeRepository,
	users repository.UserRepository,
	scoring *ScoringService,
	notifier Notifier,
) *TemoignageService {
	return &TemoignageService{
		temoignages: temoignages,
		invitations: invitations,
		soirees:     soirees,
		users:       users,
		scoring:     scoring,
		notifier:    notifier,
	}
}

// InviteTemoin implémente UC09 : seul le propriétaire de la soirée peut
// inviter, l'invité doit exister (recherché par pseudo). Notifie l'invité
// par email (Resend).
func (s *TemoignageService) InviteTemoin(ctx context.Context, ownerID, soireeID int64, invitePseudo string) error {
	soiree, err := s.soirees.GetByID(ctx, soireeID)
	if err != nil {
		return err
	}
	if soiree.UserID != ownerID {
		return ErrForbidden
	}

	invite, err := s.users.GetByPseudo(ctx, invitePseudo)
	if err != nil {
		// "Utilisateur invité inexistant -> invitation refusée"
		return err
	}

	if err := s.invitations.Create(ctx, &models.TemoinInvitation{SoireeID: soireeID, InviteID: invite.ID}); err != nil {
		return err
	}

	if s.notifier != nil {
		_ = s.notifier.SendInvitation(ctx, invite.Email, invite.Pseudo, soiree.Titre)
	}
	return nil
}

// Add implémente UC11 : seul un témoin invité peut rédiger un témoignage
// sur la soirée. Inclut UC16 (recalcul du score du propriétaire de la
// soirée, cible réelle de l'appréciation portée par le témoignage).
func (s *TemoignageService) Add(ctx context.Context, auteurID, soireeID int64, contenu string) (*models.Temoignage, error) {
	if strings.TrimSpace(contenu) == "" {
		return nil, ErrValidation
	}

	soiree, err := s.soirees.GetByID(ctx, soireeID)
	if err != nil {
		return nil, err
	}

	invited, err := s.invitations.IsInvited(ctx, soireeID, auteurID)
	if err != nil {
		return nil, err
	}
	if !invited {
		return nil, ErrForbidden
	}

	t := &models.Temoignage{SoireeID: soireeID, AuteurID: auteurID, Contenu: contenu}
	if err := s.temoignages.Create(ctx, t); err != nil {
		return nil, err
	}

	if _, err := s.scoring.Recalculate(ctx, soiree.UserID); err != nil {
		return t, err
	}
	return t, nil
}

func (s *TemoignageService) ListBySoiree(ctx context.Context, soireeID int64) ([]models.Temoignage, error) {
	return s.temoignages.ListBySoiree(ctx, soireeID)
}

func (s *TemoignageService) Get(ctx context.Context, id int64) (*models.Temoignage, error) {
	return s.temoignages.GetByID(ctx, id)
}
