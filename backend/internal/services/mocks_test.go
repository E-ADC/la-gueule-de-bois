package services

import (
	"context"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// Mocks minimalistes des interfaces repository, utilisés uniquement dans
// les tests de ce paquet — pas de base de données vivante requise.

type fakeUserRepo struct {
	users map[int64]*models.User
}

func newFakeUserRepo(users ...*models.User) *fakeUserRepo {
	m := make(map[int64]*models.User)
	for _, u := range users {
		m[u.ID] = u
	}
	return &fakeUserRepo{users: m}
}

func (f *fakeUserRepo) Create(ctx context.Context, u *models.User) error {
	f.users[u.ID] = u
	return nil
}
func (f *fakeUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	if u, ok := f.users[id]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) GetByPseudo(ctx context.Context, pseudo string) (*models.User, error) {
	for _, u := range f.users {
		if u.Pseudo == pseudo {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) Update(ctx context.Context, u *models.User) error {
	f.users[u.ID] = u
	return nil
}
func (f *fakeUserRepo) UpdateScore(ctx context.Context, userID int64, score int) error {
	if u, ok := f.users[userID]; ok {
		u.Score = score
		return nil
	}
	return repository.ErrNotFound
}
func (f *fakeUserRepo) ListLeaderboard(ctx context.Context, limit int) ([]models.User, error) {
	return nil, nil
}
func (f *fakeUserRepo) ListLeaderboardForGroup(ctx context.Context, groupeID int64, limit int) ([]models.User, error) {
	return nil, nil
}

type fakeSoireeRepo struct {
	countByUser map[int64]int
}

func (f *fakeSoireeRepo) Create(ctx context.Context, s *models.Soiree) error { return nil }
func (f *fakeSoireeRepo) GetByID(ctx context.Context, id int64) (*models.Soiree, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeSoireeRepo) Update(ctx context.Context, s *models.Soiree) error { return nil }
func (f *fakeSoireeRepo) Delete(ctx context.Context, id int64) error         { return nil }
func (f *fakeSoireeRepo) ListByUser(ctx context.Context, userID int64) ([]models.Soiree, error) {
	return nil, nil
}
func (f *fakeSoireeRepo) CountByUser(ctx context.Context, userID int64) (int, error) {
	return f.countByUser[userID], nil
}

type fakeTemoignageRepo struct {
	countForOwner map[int64]int
}

func (f *fakeTemoignageRepo) Create(ctx context.Context, t *models.Temoignage) error { return nil }
func (f *fakeTemoignageRepo) GetByID(ctx context.Context, id int64) (*models.Temoignage, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeTemoignageRepo) ListBySoiree(ctx context.Context, soireeID int64) ([]models.Temoignage, error) {
	return nil, nil
}
func (f *fakeTemoignageRepo) CountForOwner(ctx context.Context, ownerID int64) (int, error) {
	return f.countForOwner[ownerID], nil
}
func (f *fakeTemoignageRepo) Delete(ctx context.Context, id int64) error { return nil }

type fakeVoteRepo struct {
	positifs map[int64]int
	negatifs map[int64]int
}

func (f *fakeVoteRepo) Create(ctx context.Context, v *models.Vote) error { return nil }
func (f *fakeVoteRepo) Exists(ctx context.Context, temoignageID, userID int64) (bool, error) {
	return false, nil
}
func (f *fakeVoteRepo) SommeVotesForOwner(ctx context.Context, ownerID int64) (int, int, error) {
	return f.positifs[ownerID], f.negatifs[ownerID], nil
}

type fakeBadgeRepo struct {
	catalogue []models.Badge
	possedes  map[int64][]models.Badge
	attaches  []int64 // badgeID attachés, dans l'ordre
}

func (f *fakeBadgeRepo) ListAll(ctx context.Context) ([]models.Badge, error) { return f.catalogue, nil }
func (f *fakeBadgeRepo) ListForUser(ctx context.Context, userID int64) ([]models.Badge, error) {
	return f.possedes[userID], nil
}
func (f *fakeBadgeRepo) AttachToUser(ctx context.Context, userID, badgeID int64) error {
	f.attaches = append(f.attaches, badgeID)
	f.possedes[userID] = append(f.possedes[userID], mustFindBadge(f.catalogue, badgeID))
	return nil
}

func mustFindBadge(catalogue []models.Badge, id int64) models.Badge {
	for _, b := range catalogue {
		if b.ID == id {
			return b
		}
	}
	return models.Badge{}
}
