package services

import (
	"context"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// DemandeAmiService implémente UC21 (envoyer une demande d'ami) et sa
// réponse (accepter/refuser — nécessaire pour que UC21 soit exploitable,
// la fiche ne détaille pas cette étape mais la présuppose : "en attente de
// réponse").
type DemandeAmiService struct {
	demandes repository.DemandeAmiRepository
	users    repository.UserRepository
	notifier Notifier
}

func NewDemandeAmiService(demandes repository.DemandeAmiRepository, users repository.UserRepository, notifier Notifier) *DemandeAmiService {
	return &DemandeAmiService{demandes: demandes, users: users, notifier: notifier}
}

// statutsBloquants : une demande en attente ou déjà acceptée dans un sens ou
// l'autre bloque l'envoi d'une nouvelle demande ("déjà envoyée ou déjà amis
// -> ignorée").
var statutsBloquants = []models.StatutDemandeAmi{models.DemandeAmiEnAttente, models.DemandeAmiAcceptee}

// Envoyer implémente UC21 : destinataire recherché par pseudo.
func (s *DemandeAmiService) Envoyer(ctx context.Context, demandeurID int64, destinatairePseudo string) error {
	destinataire, err := s.users.GetByPseudo(ctx, destinatairePseudo)
	if err != nil {
		return err
	}
	if destinataire.ID == demandeurID {
		return ErrValidation
	}

	bloque, err := s.demandes.ExisteEntre(ctx, demandeurID, destinataire.ID, statutsBloquants)
	if err != nil {
		return err
	}
	if bloque {
		return repository.ErrConflict
	}

	if err := s.demandes.Create(ctx, &models.DemandeAmi{
		DemandeurID:    demandeurID,
		DestinataireID: destinataire.ID,
	}); err != nil {
		return err
	}

	if s.notifier != nil {
		demandeur, err := s.users.GetByID(ctx, demandeurID)
		if err == nil {
			_ = s.notifier.SendFriendRequest(ctx, destinataire.Email, destinataire.Pseudo, demandeur.Pseudo)
		}
	}
	return nil
}

// ListRecues liste les demandes en attente reçues par l'utilisateur.
func (s *DemandeAmiService) ListRecues(ctx context.Context, userID int64) ([]models.DemandeAmi, error) {
	return s.demandes.ListRecues(ctx, userID)
}

// Repondre accepte ou refuse une demande reçue. Seul le destinataire peut
// répondre ; une demande déjà traitée est ignorée (409).
func (s *DemandeAmiService) Repondre(ctx context.Context, userID, demandeID int64, accepter bool) error {
	demande, err := s.demandes.GetByID(ctx, demandeID)
	if err != nil {
		return err
	}
	if demande.DestinataireID != userID {
		return ErrForbidden
	}
	if demande.Statut != models.DemandeAmiEnAttente {
		return repository.ErrConflict
	}

	statut := models.DemandeAmiRefusee
	if accepter {
		statut = models.DemandeAmiAcceptee
	}
	return s.demandes.MarkStatut(ctx, demandeID, statut)
}
