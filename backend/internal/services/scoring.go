package services

import (
	"context"
	"fmt"

	"gueuledebois/backend/internal/repository"
)

// ScoreInput regroupe les données brutes nécessaires au calcul du score
// d'un utilisateur (UC16). Ces données sont fournies par les soirées qu'il
// a créées et les témoignages/votes reçus dessus.
type ScoreInput struct {
	NbSoirees     int
	NbTemoignages int
	VotesPositifs int
	VotesNegatifs int
}

// Barème du score (UC16). La fiche ne fournit aucun barème chiffré :
// choix documenté ici, à ajuster facilement si besoin.
//   - +10 points par soirée créée : c'est l'acte fondateur de l'app.
//   - +5 points par témoignage reçu sur ses soirées : encourage à inviter
//     des témoins et à jouer le jeu du récit.
//   - +1 / -1 point par vote (swipe) reçu sur ces témoignages : reflète
//     l'appréciation de la communauté.
//   - Le score ne descend jamais sous 0 (pas de score négatif affiché au
//     classement).
const (
	pointsParSoiree      = 10
	pointsParTemoignage  = 5
	pointsParVotePositif = 1
	pointsParVoteNegatif = 1
)

// ComputeScore est une fonction pure : facilement testable en table-driven
// sans dépendance à la base de données.
func ComputeScore(in ScoreInput) int {
	score := in.NbSoirees*pointsParSoiree +
		in.NbTemoignages*pointsParTemoignage +
		in.VotesPositifs*pointsParVotePositif -
		in.VotesNegatifs*pointsParVoteNegatif
	if score < 0 {
		score = 0
	}
	return score
}

// ScoringService orchestre le recalcul du score d'un utilisateur (UC16),
// puis l'évaluation des badges (UC14, inclus par UC16).
type ScoringService struct {
	users       repository.UserRepository
	soirees     repository.SoireeRepository
	temoignages repository.TemoignageRepository
	votes       repository.VoteRepository
	badges      *BadgeService
}

func NewScoringService(
	users repository.UserRepository,
	soirees repository.SoireeRepository,
	temoignages repository.TemoignageRepository,
	votes repository.VoteRepository,
	badges *BadgeService,
) *ScoringService {
	return &ScoringService{
		users:       users,
		soirees:     soirees,
		temoignages: temoignages,
		votes:       votes,
		badges:      badges,
	}
}

// Recalculate implémente UC16 : recalcule le score de l'utilisateur,
// l'enregistre, puis inclut UC14 (déblocage de badges).
func (s *ScoringService) Recalculate(ctx context.Context, userID int64) (newScore int, err error) {
	nbSoirees, err := s.soirees.CountByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("scoring: comptage soirées: %w", err)
	}
	nbTemoignages, err := s.temoignages.CountForOwner(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("scoring: comptage témoignages: %w", err)
	}
	votesPositifs, votesNegatifs, err := s.votes.SommeVotesForOwner(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("scoring: agrégation votes: %w", err)
	}

	score := ComputeScore(ScoreInput{
		NbSoirees:     nbSoirees,
		NbTemoignages: nbTemoignages,
		VotesPositifs: votesPositifs,
		VotesNegatifs: votesNegatifs,
	})

	if err := s.users.UpdateScore(ctx, userID, score); err != nil {
		return 0, fmt.Errorf("scoring: mise à jour score: %w", err)
	}

	if s.badges != nil {
		if _, err := s.badges.EvaluateAndUnlock(ctx, userID, score); err != nil {
			return score, fmt.Errorf("scoring: évaluation badges: %w", err)
		}
	}

	return score, nil
}
