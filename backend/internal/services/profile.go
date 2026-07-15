package services

import (
	"context"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// ProfileService regroupe les lectures simples liées au profil : UC05
// (consulter un profil), UC15 (consulter ses badges), UC17/UC20
// (classement global / de groupe).
type ProfileService struct {
	users   repository.UserRepository
	badges  repository.BadgeRepository
	groupes repository.GroupeRepository
}

func NewProfileService(users repository.UserRepository, badges repository.BadgeRepository, groupes repository.GroupeRepository) *ProfileService {
	return &ProfileService{users: users, badges: badges, groupes: groupes}
}

// GetPublicProfile implémente UC05.
func (s *ProfileService) GetPublicProfile(ctx context.Context, userID int64) (*models.User, error) {
	return s.users.GetByID(ctx, userID)
}

// ListBadges implémente UC15.
func (s *ProfileService) ListBadges(ctx context.Context, userID int64) ([]models.Badge, error) {
	return s.badges.ListForUser(ctx, userID)
}

// AllBadges liste le catalogue complet, pour afficher aussi les badges à
// débloquer (UC15 : "badges obtenus et à débloquer").
func (s *ProfileService) AllBadges(ctx context.Context) ([]models.Badge, error) {
	return s.badges.ListAll(ctx)
}

const leaderboardLimit = 100

// Leaderboard implémente UC17.
func (s *ProfileService) Leaderboard(ctx context.Context) ([]models.User, error) {
	return s.users.ListLeaderboard(ctx, leaderboardLimit)
}

// LeaderboardForGroup implémente UC20 : restreint au groupe, l'utilisateur
// doit en être membre ("non membre -> accès refusé", 403).
func (s *ProfileService) LeaderboardForGroup(ctx context.Context, requesterID, groupeID int64) ([]models.User, error) {
	member, err := s.groupes.IsMember(ctx, groupeID, requesterID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, ErrForbidden
	}
	return s.users.ListLeaderboardForGroup(ctx, groupeID, leaderboardLimit)
}
