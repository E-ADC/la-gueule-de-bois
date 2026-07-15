package services

import (
	"context"
	"fmt"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// EvaluateBadges est une fonction pure (UC14) : à partir d'un score, du
// catalogue complet de badges et de l'ensemble des badges déjà possédés,
// retourne la liste des badges nouvellement débloqués. Testable en
// table-driven sans base de données.
func EvaluateBadges(score int, catalogue []models.Badge, dejaPossedes map[int64]bool) []models.Badge {
	var debloques []models.Badge
	for _, b := range catalogue {
		if dejaPossedes[b.ID] {
			continue
		}
		if score >= b.SeuilScore {
			debloques = append(debloques, b)
		}
	}
	return debloques
}

// BadgeService orchestre l'évaluation et l'attribution des badges (UC14),
// puis la notification de l'utilisateur via Notifier.
type BadgeService struct {
	badges   repository.BadgeRepository
	users    repository.UserRepository
	notifier Notifier
}

func NewBadgeService(badges repository.BadgeRepository, users repository.UserRepository, notifier Notifier) *BadgeService {
	return &BadgeService{badges: badges, users: users, notifier: notifier}
}

// EvaluateAndUnlock implémente UC14 : évalue les critères de badges avec le
// score fourni (déjà recalculé par UC16), attribue les nouveaux badges et
// notifie l'utilisateur par email. "Aucun critère atteint -> aucun badge
// attribué" (exception de la fiche) : retourne simplement une liste vide.
func (s *BadgeService) EvaluateAndUnlock(ctx context.Context, userID int64, score int) ([]models.Badge, error) {
	catalogue, err := s.badges.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("badges: liste catalogue: %w", err)
	}
	possedes, err := s.badges.ListForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("badges: liste utilisateur: %w", err)
	}

	dejaPossedes := make(map[int64]bool, len(possedes))
	for _, b := range possedes {
		dejaPossedes[b.ID] = true
	}

	nouveaux := EvaluateBadges(score, catalogue, dejaPossedes)
	if len(nouveaux) == 0 {
		return nil, nil
	}

	for _, b := range nouveaux {
		if err := s.badges.AttachToUser(ctx, userID, b.ID); err != nil {
			return nil, fmt.Errorf("badges: attribution %s: %w", b.Code, err)
		}
	}

	if s.notifier != nil {
		if user, err := s.users.GetByID(ctx, userID); err == nil {
			for _, b := range nouveaux {
				_ = s.notifier.SendBadgeUnlocked(ctx, user.Email, user.Pseudo, b.Nom)
			}
		}
	}

	return nouveaux, nil
}
