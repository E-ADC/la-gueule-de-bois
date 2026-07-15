package services

import (
	"context"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// VoteService implémente UC12 (swiper/voter sur un témoignage).
type VoteService struct {
	votes       repository.VoteRepository
	temoignages repository.TemoignageRepository
	soirees     repository.SoireeRepository
	scoring     *ScoringService
}

func NewVoteService(votes repository.VoteRepository, temoignages repository.TemoignageRepository, soirees repository.SoireeRepository, scoring *ScoringService) *VoteService {
	return &VoteService{votes: votes, temoignages: temoignages, soirees: soirees, scoring: scoring}
}

// Cast implémente UC12 : "l'utilisateur a déjà voté sur ce témoignage ->
// vote ignoré" (409 doublon, cf. mapping d'erreurs de la spec).
//
// Choix : la fiche UC12 ne liste pas explicitement UC16 parmi les cas
// inclus, contrairement à UC06/07/08/11. Comme le score (UC16) intègre
// les votes reçus par un utilisateur sur ses témoignages, on déclenche
// tout de même un recalcul ici pour éviter un score qui resterait figé
// entre deux créations de soirée/témoignage — c'est une extension
// raisonnable de la fiche, pas une contradiction.
func (s *VoteService) Cast(ctx context.Context, userID, temoignageID int64, valeur models.VoteValeur) (*models.Vote, error) {
	exists, err := s.votes.Exists(ctx, temoignageID, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, repository.ErrConflict
	}

	t, err := s.temoignages.GetByID(ctx, temoignageID)
	if err != nil {
		return nil, err
	}

	v := &models.Vote{TemoignageID: temoignageID, UserID: userID, Valeur: valeur}
	if err := s.votes.Create(ctx, v); err != nil {
		return nil, err
	}

	// Le score à recalculer est celui du propriétaire de la soirée liée au
	// témoignage voté (cf. logique de scoring.go), pas celui du votant.
	if soiree, err := s.soirees.GetByID(ctx, t.SoireeID); err == nil {
		_, _ = s.scoring.Recalculate(ctx, soiree.UserID)
	}

	return v, nil
}

// Counts renvoie les compteurs de votes positifs/négatifs d'un témoignage,
// pour l'affichage enrichi (SoireeDetailPage).
func (s *VoteService) Counts(ctx context.Context, temoignageID int64) (positifs int, negatifs int, err error) {
	return s.votes.CountForTemoignage(ctx, temoignageID)
}

// MonVote renvoie le vote de cet utilisateur sur ce témoignage (nil s'il
// n'a pas encore voté).
func (s *VoteService) MonVote(ctx context.Context, userID, temoignageID int64) (*models.Vote, error) {
	v, err := s.votes.GetByUserAndTemoignage(ctx, userID, temoignageID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}
