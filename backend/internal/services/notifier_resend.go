package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// resendAPIURL est l'endpoint HTTPS de l'API Resend (pas de SMTP sortant,
// cf. spec — évite un éventuel blocage du port 25 sur le VPS).
const resendAPIURL = "https://api.resend.com/emails"

// ResendFromAddress est l'expéditeur utilisé pour tous les emails. En mode
// test Resend (sans domaine vérifié), seules les adresses de l'équipe
// peuvent recevoir des emails — limitation acceptée pour la démo (cf. spec).
const ResendFromAddress = "La Gueule de Bois <onboarding@resend.dev>"

// ResendNotifier implémente Notifier via l'API HTTPS de Resend.
type ResendNotifier struct {
	apiKey     string
	httpClient *http.Client
}

func NewResendNotifier(apiKey string) *ResendNotifier {
	return &ResendNotifier{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func (n *ResendNotifier) send(ctx context.Context, toEmail, subject, html string) error {
	body, err := json.Marshal(resendEmailRequest{
		From:    ResendFromAddress,
		To:      []string{toEmail},
		Subject: subject,
		HTML:    html,
	})
	if err != nil {
		return fmt.Errorf("resend: encodage requête: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("resend: construction requête: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.apiKey)

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend: appel API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("resend: réponse API inattendue (%d)", resp.StatusCode)
	}
	return nil
}

func (n *ResendNotifier) SendInvitation(ctx context.Context, toEmail, toPseudo, soireeTitre string) error {
	subject := "Tu as été invité comme témoin"
	html := fmt.Sprintf("<p>Salut %s,</p><p>Tu as été invité comme témoin sur la soirée « %s ».</p>", toPseudo, soireeTitre)
	return n.send(ctx, toEmail, subject, html)
}

func (n *ResendNotifier) SendBadgeUnlocked(ctx context.Context, toEmail, toPseudo, badgeNom string) error {
	subject := "Nouveau badge débloqué !"
	html := fmt.Sprintf("<p>Bravo %s, tu as débloqué le badge « %s » !</p>", toPseudo, badgeNom)
	return n.send(ctx, toEmail, subject, html)
}

func (n *ResendNotifier) SendFriendRequest(ctx context.Context, toEmail, toPseudo, fromPseudo string) error {
	subject := "Nouvelle demande d'ami"
	html := fmt.Sprintf("<p>Salut %s,</p><p>%s souhaite t'ajouter comme ami.</p>", toPseudo, fromPseudo)
	return n.send(ctx, toEmail, subject, html)
}

func (n *ResendNotifier) SendReportResolved(ctx context.Context, toEmail, toPseudo string, temoignageSupprime bool) error {
	subject := "Ton témoignage signalé a été traité"
	decision := "a été conservé"
	if temoignageSupprime {
		decision = "a été supprimé"
	}
	html := fmt.Sprintf("<p>Salut %s,</p><p>Ton témoignage signalé %s après examen par un modérateur.</p>", toPseudo, decision)
	return n.send(ctx, toEmail, subject, html)
}
