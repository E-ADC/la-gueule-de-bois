package services

import (
	"context"
	"reflect"
	"testing"

	"gueuledebois/backend/internal/models"
)

func TestEvaluateBadges(t *testing.T) {
	catalogue := []models.Badge{
		{ID: 1, Code: "premiere-cuite", SeuilScore: 10},
		{ID: 2, Code: "habitue-du-bar", SeuilScore: 50},
		{ID: 3, Code: "legende-de-la-soiree", SeuilScore: 150},
	}

	tests := []struct {
		name         string
		score        int
		dejaPossedes map[int64]bool
		wantCodes    []string
	}{
		{
			name:         "score sous le premier seuil -> aucun badge",
			score:        5,
			dejaPossedes: map[int64]bool{},
			wantCodes:    nil,
		},
		{
			name:         "score franchit le premier seuil uniquement",
			score:        10,
			dejaPossedes: map[int64]bool{},
			wantCodes:    []string{"premiere-cuite"},
		},
		{
			name:         "score franchit deux seuils d'un coup",
			score:        75,
			dejaPossedes: map[int64]bool{},
			wantCodes:    []string{"premiere-cuite", "habitue-du-bar"},
		},
		{
			name:         "badge déjà possédé n'est pas re-débloqué",
			score:        75,
			dejaPossedes: map[int64]bool{1: true},
			wantCodes:    []string{"habitue-du-bar"},
		},
		{
			name:         "tous les badges déjà possédés -> rien de nouveau",
			score:        999,
			dejaPossedes: map[int64]bool{1: true, 2: true, 3: true},
			wantCodes:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateBadges(tt.score, catalogue, tt.dejaPossedes)
			var gotCodes []string
			for _, b := range got {
				gotCodes = append(gotCodes, b.Code)
			}
			if !reflect.DeepEqual(gotCodes, tt.wantCodes) {
				t.Errorf("EvaluateBadges() codes = %v, want %v", gotCodes, tt.wantCodes)
			}
		})
	}
}

func TestBadgeService_EvaluateAndUnlock_AucunCritereAtteint(t *testing.T) {
	// Exception de la fiche UC14 : "Aucun critère atteint -> aucun badge attribué".
	user := &models.User{ID: 1, Email: "u1@test.local", Pseudo: "u1"}
	users := newFakeUserRepo(user)
	badgeRepo := &fakeBadgeRepo{
		catalogue: []models.Badge{{ID: 1, Code: "premiere-cuite", SeuilScore: 10}},
		possedes:  map[int64][]models.Badge{},
	}
	notifier := NewMockNotifier()
	svc := NewBadgeService(badgeRepo, users, notifier)

	got, err := svc.EvaluateAndUnlock(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("EvaluateAndUnlock() erreur inattendue: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("badges débloqués = %v, want aucun", got)
	}
	if len(notifier.Sent) != 0 {
		t.Errorf("notifications envoyées = %v, want aucune", notifier.Sent)
	}
	if len(badgeRepo.attaches) != 0 {
		t.Errorf("badges attachés en base = %v, want aucun", badgeRepo.attaches)
	}
}
