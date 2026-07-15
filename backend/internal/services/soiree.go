package services

import (
	"context"
	"strings"
	"time"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// SoireeService implémente la logique métier de UC06 (créer), UC07
// (modifier), UC08 (supprimer) et UC10 (historique). Chaque mutation
// inclut UC16 (recalcul du score du propriétaire).
type SoireeService struct {
	soirees repository.SoireeRepository
	photos  repository.PhotoRepository
	scoring *ScoringService
}

func NewSoireeService(soirees repository.SoireeRepository, photos repository.PhotoRepository, scoring *ScoringService) *SoireeService {
	return &SoireeService{soirees: soirees, photos: photos, scoring: scoring}
}

// CreateInput regroupe les champs saisissables à la création (UC06).
type CreateSoireeInput struct {
	Titre       string
	DateSoiree  time.Time
	Lieu        string
	Description string
}

func (s *SoireeService) validate(in CreateSoireeInput) error {
	if strings.TrimSpace(in.Titre) == "" {
		return ErrValidation
	}
	if in.DateSoiree.IsZero() {
		return ErrValidation
	}
	return nil
}

// Create implémente UC06 : "champs obligatoires manquants -> création
// refusée" (titre + date requis, lieu/description optionnels).
func (s *SoireeService) Create(ctx context.Context, userID int64, in CreateSoireeInput) (*models.Soiree, error) {
	if err := s.validate(in); err != nil {
		return nil, err
	}

	soiree := &models.Soiree{
		UserID:      userID,
		Titre:       in.Titre,
		DateSoiree:  in.DateSoiree,
		Lieu:        in.Lieu,
		Description: in.Description,
	}
	if err := s.soirees.Create(ctx, soiree); err != nil {
		return nil, err
	}

	if _, err := s.scoring.Recalculate(ctx, userID); err != nil {
		return soiree, err
	}
	return soiree, nil
}

// Update implémente UC07 : seul le propriétaire peut modifier (403 sinon).
func (s *SoireeService) Update(ctx context.Context, userID, soireeID int64, in CreateSoireeInput) (*models.Soiree, error) {
	if err := s.validate(in); err != nil {
		return nil, err
	}

	existing, err := s.soirees.GetByID(ctx, soireeID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
		return nil, ErrForbidden
	}

	existing.Titre = in.Titre
	existing.DateSoiree = in.DateSoiree
	existing.Lieu = in.Lieu
	existing.Description = in.Description

	if err := s.soirees.Update(ctx, existing); err != nil {
		return nil, err
	}
	if _, err := s.scoring.Recalculate(ctx, userID); err != nil {
		return existing, err
	}
	return existing, nil
}

// Delete implémente UC08 : seul le propriétaire peut supprimer.
func (s *SoireeService) Delete(ctx context.Context, userID, soireeID int64) error {
	existing, err := s.soirees.GetByID(ctx, soireeID)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return ErrForbidden
	}

	if err := s.soirees.Delete(ctx, soireeID); err != nil {
		return err
	}
	_, err = s.scoring.Recalculate(ctx, userID)
	return err
}

// Get récupère une soirée par id (404 si absente).
func (s *SoireeService) Get(ctx context.Context, soireeID int64) (*models.Soiree, error) {
	return s.soirees.GetByID(ctx, soireeID)
}

// ListByUser implémente UC10 : historique des soirées d'un utilisateur.
func (s *SoireeService) ListByUser(ctx context.Context, userID int64) ([]models.Soiree, error) {
	return s.soirees.ListByUser(ctx, userID)
}

// AddPhoto attache une photo déjà écrite sur disque à une soirée. Seul le
// propriétaire de la soirée peut ajouter une photo.
func (s *SoireeService) AddPhoto(ctx context.Context, userID, soireeID int64, chemin string) (*models.Photo, error) {
	soiree, err := s.soirees.GetByID(ctx, soireeID)
	if err != nil {
		return nil, err
	}
	if soiree.UserID != userID {
		return nil, ErrForbidden
	}

	photo := &models.Photo{SoireeID: soireeID, Chemin: chemin}
	if err := s.photos.Create(ctx, photo); err != nil {
		return nil, err
	}
	return photo, nil
}

func (s *SoireeService) ListPhotos(ctx context.Context, soireeID int64) ([]models.Photo, error) {
	return s.photos.ListBySoiree(ctx, soireeID)
}
