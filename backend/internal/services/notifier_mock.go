package services

import "context"

// MockNotifier est une implémentation en mémoire de Notifier, utilisée en
// tests (aucun appel réseau) et disponible en développement local sans clé
// Resend.
type MockNotifier struct {
	Sent []MockEmail
}

// MockEmail trace un envoi simulé, inspectable dans les tests.
type MockEmail struct {
	Type    string
	ToEmail string
	Subject string
}

func NewMockNotifier() *MockNotifier {
	return &MockNotifier{}
}

func (m *MockNotifier) SendInvitation(ctx context.Context, toEmail, toPseudo, soireeTitre string) error {
	m.Sent = append(m.Sent, MockEmail{Type: "invitation", ToEmail: toEmail, Subject: soireeTitre})
	return nil
}

func (m *MockNotifier) SendBadgeUnlocked(ctx context.Context, toEmail, toPseudo, badgeNom string) error {
	m.Sent = append(m.Sent, MockEmail{Type: "badge_unlocked", ToEmail: toEmail, Subject: badgeNom})
	return nil
}

func (m *MockNotifier) SendFriendRequest(ctx context.Context, toEmail, toPseudo, fromPseudo string) error {
	m.Sent = append(m.Sent, MockEmail{Type: "friend_request", ToEmail: toEmail, Subject: fromPseudo})
	return nil
}

func (m *MockNotifier) SendReportResolved(ctx context.Context, toEmail, toPseudo string, temoignageSupprime bool) error {
	subject := "conserve"
	if temoignageSupprime {
		subject = "supprime"
	}
	m.Sent = append(m.Sent, MockEmail{Type: "report_resolved", ToEmail: toEmail, Subject: subject})
	return nil
}
