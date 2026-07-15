package services

import (
	"context"
	"testing"

	"gueuledebois/backend/internal/models"
)

func TestComputeScore(t *testing.T) {
	tests := []struct {
		name string
		in   ScoreInput
		want int
	}{
		{
			name: "aucune activité -> score nul",
			in:   ScoreInput{},
			want: 0,
		},
		{
			name: "une soirée créée",
			in:   ScoreInput{NbSoirees: 1},
			want: 10,
		},
		{
			name: "soirées + témoignages + votes positifs",
			in:   ScoreInput{NbSoirees: 2, NbTemoignages: 3, VotesPositifs: 5},
			want: 2*10 + 3*5 + 5*1, // = 40
		},
		{
			name: "votes négatifs réduisent le score",
			in:   ScoreInput{NbSoirees: 1, VotesNegatifs: 3},
			want: 10 - 3, // = 7
		},
		{
			name: "le score ne descend jamais sous zéro",
			in:   ScoreInput{NbSoirees: 0, NbTemoignages: 0, VotesNegatifs: 50},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeScore(tt.in)
			if got != tt.want {
				t.Errorf("ComputeScore(%+v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestScoringService_Recalculate(t *testing.T) {
	user := &models.User{ID: 1, Email: "u1@test.local", Pseudo: "u1"}

	users := newFakeUserRepo(user)
	soirees := &fakeSoireeRepo{countByUser: map[int64]int{1: 2}}           // 2*10 = 20
	temoignages := &fakeTemoignageRepo{countForOwner: map[int64]int{1: 1}} // 1*5 = 5
	votes := &fakeVoteRepo{
		positifs: map[int64]int{1: 3}, // 3*1 = 3
		negatifs: map[int64]int{1: 1}, // -1
	}
	// score attendu : 20 + 5 + 3 - 1 = 27

	badgeRepo := &fakeBadgeRepo{
		catalogue: []models.Badge{
			{ID: 1, Code: "premiere-cuite", Nom: "Première Cuite", SeuilScore: 10},
			{ID: 2, Code: "habitue-du-bar", Nom: "Habitué du Bar", SeuilScore: 50},
		},
		possedes: map[int64][]models.Badge{},
	}
	notifier := NewMockNotifier()
	badgeService := NewBadgeService(badgeRepo, users, notifier)

	scoring := NewScoringService(users, soirees, temoignages, votes, badgeService)

	got, err := scoring.Recalculate(context.Background(), 1)
	if err != nil {
		t.Fatalf("Recalculate() erreur inattendue: %v", err)
	}
	if got != 27 {
		t.Errorf("Recalculate() score = %d, want 27", got)
	}
	if user.Score != 27 {
		t.Errorf("le score stocké sur l'utilisateur = %d, want 27", user.Score)
	}

	// Le score de 27 franchit le seuil du badge "premiere-cuite" (10) mais
	// pas "habitue-du-bar" (50) : un seul badge doit être débloqué et notifié.
	if len(badgeRepo.attaches) != 1 || badgeRepo.attaches[0] != 1 {
		t.Errorf("badges attachés = %v, want [1]", badgeRepo.attaches)
	}
	if len(notifier.Sent) != 1 || notifier.Sent[0].Type != "badge_unlocked" {
		t.Errorf("notifications envoyées = %+v, want 1 notification badge_unlocked", notifier.Sent)
	}
}
