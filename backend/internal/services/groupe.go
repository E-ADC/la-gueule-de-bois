package services

import (
	"context"
	"strings"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// GroupeService implémente UC18 (créer un groupe) et UC19 (rejoindre un
// groupe).
type GroupeService struct {
	groupes repository.GroupeRepository
}

func NewGroupeService(groupes repository.GroupeRepository) *GroupeService {
	return &GroupeService{groupes: groupes}
}

// Create implémente UC18 : "nom déjà utilisé -> création refusée" remonte
// via repository.ErrConflict (contrainte unique groupes.nom). Le créateur
// devient automatiquement membre.
func (s *GroupeService) Create(ctx context.Context, createurID int64, nom string) (*models.Groupe, error) {
	if strings.TrimSpace(nom) == "" {
		return nil, ErrValidation
	}
	g := &models.Groupe{Nom: nom, CreateurID: createurID}
	if err := s.groupes.Create(ctx, g); err != nil {
		return nil, err
	}
	if err := s.groupes.AddMember(ctx, g.ID, createurID); err != nil {
		return nil, err
	}
	return g, nil
}

// Join implémente UC19 : "déjà membre -> action ignorée", AddMember étant
// idempotent (ON CONFLICT DO NOTHING) côté repository.
func (s *GroupeService) Join(ctx context.Context, userID, groupeID int64) error {
	if _, err := s.groupes.GetByID(ctx, groupeID); err != nil {
		return err
	}
	return s.groupes.AddMember(ctx, groupeID, userID)
}

// ListMine liste les groupes dont l'utilisateur est membre.
func (s *GroupeService) ListMine(ctx context.Context, userID int64) ([]models.Groupe, error) {
	return s.groupes.ListForUser(ctx, userID)
}
