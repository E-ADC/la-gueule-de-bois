package services

import (
	"context"
	"strings"
	"time"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// SignalementService implémente UC13 (signaler un témoignage) et UC22
// (traiter un signalement, acteur Modérateur).
type SignalementService struct {
	signalements repository.SignalementRepository
	temoignages  repository.TemoignageRepository
	soirees      repository.SoireeRepository
	users        repository.UserRepository
	scoring      *ScoringService
	notifier     Notifier
}

func NewSignalementService(
	signalements repository.SignalementRepository,
	temoignages repository.TemoignageRepository,
	soirees repository.SoireeRepository,
	users repository.UserRepository,
	scoring *ScoringService,
	notifier Notifier,
) *SignalementService {
	return &SignalementService{
		signalements: signalements,
		temoignages:  temoignages,
		soirees:      soirees,
		users:        users,
		scoring:      scoring,
		notifier:     notifier,
	}
}

// Report implémente UC13 : signale un témoignage. La contrainte unique
// (temoignage_id, auteur_id) en base fait remonter repository.ErrConflict
// si cet utilisateur a déjà signalé ce témoignage ("signalement ignoré").
func (s *SignalementService) Report(ctx context.Context, auteurID, temoignageID int64, motif string) error {
	if strings.TrimSpace(motif) == "" {
		return ErrValidation
	}
	if _, err := s.temoignages.GetByID(ctx, temoignageID); err != nil {
		return err
	}
	return s.signalements.Create(ctx, &models.Signalement{
		TemoignageID: &temoignageID,
		AuteurID:     auteurID,
		Motif:        motif,
	})
}

// ListEnAttente implémente UC22 étape 1 : les signalements en attente de
// traitement par un modérateur.
func (s *SignalementService) ListEnAttente(ctx context.Context) ([]models.Signalement, error) {
	return s.signalements.ListEnAttente(ctx)
}

// Traiter implémente UC22 : le modérateur rejette le signalement ou
// supprime le témoignage signalé (ce qui déclenche UC16, le score de
// l'auteur de la soirée dépendant du nombre de témoignages reçus), puis
// notifie l'auteur du témoignage par email dans les deux cas.
func (s *SignalementService) Traiter(ctx context.Context, moderateurID, signalementID int64, supprimer bool) error {
	signalement, err := s.signalements.GetByID(ctx, signalementID)
	if err != nil {
		return err
	}
	// "Signalement déjà traité -> action ignorée"
	if signalement.Statut != models.SignalementEnAttente {
		return repository.ErrConflict
	}
	// Un signalement "en_attente" référence toujours un témoignage existant
	// (temoignage_id ne passe à NULL qu'après ce traitement, cf. migration
	// 000003) : le pointeur est garanti non nil ici.
	temoignage, err := s.temoignages.GetByID(ctx, *signalement.TemoignageID)
	if err != nil {
		return err
	}
	auteur, err := s.users.GetByID(ctx, temoignage.AuteurID)
	if err != nil {
		return err
	}

	statut := models.SignalementRejete
	if supprimer {
		statut = models.SignalementSupprime
		if err := s.temoignages.Delete(ctx, temoignage.ID); err != nil {
			return err
		}
		soiree, err := s.soirees.GetByID(ctx, temoignage.SoireeID)
		if err != nil {
			return err
		}
		if s.scoring != nil {
			if _, err := s.scoring.Recalculate(ctx, soiree.UserID); err != nil {
				return err
			}
		}
	}

	if err := s.signalements.MarkTraite(ctx, signalementID, statut, moderateurID, time.Now()); err != nil {
		return err
	}

	if s.notifier != nil {
		_ = s.notifier.SendReportResolved(ctx, auteur.Email, auteur.Pseudo, supprimer)
	}
	return nil
}
