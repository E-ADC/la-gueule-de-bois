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

func emailLayout(titre, corps string) string {
	return fmt.Sprintf(`<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4ede0; padding: 24px 0;">
  <tr><td align="center">
    <table role="presentation" width="480" cellpadding="0" cellspacing="0" style="background-color:#f9efdb; border-radius:8px; overflow:hidden; font-family: Helvetica, Arial, sans-serif;">
      <tr><td style="background-color:#b5651d; padding: 20px 24px; text-align:center;">
        <span style="color:#ffffff; font-size: 20px; font-weight: bold;">La Gueule de Bois</span>
      </td></tr>
      <tr><td style="padding: 24px 24px 0; color:#4a2c17; font-size: 16px; font-weight: bold;">
        %s
      </td></tr>
      <tr><td style="padding: 12px 24px 24px; color:#4a2c17; font-size: 15px; line-height: 1.6;">
        %s
      </td></tr>
      <tr><td style="padding: 16px 24px; text-align:center; border-top: 1px solid #e0d5c8;">
        <span style="color:#a08a73; font-size: 12px;">La Gueule de Bois — notification automatique</span>
      </td></tr>
    </table>
  </td></tr>
</table>`, titre, corps)
}

func (n *ResendNotifier) SendInvitation(ctx context.Context, toEmail, toPseudo, soireeTitre string) error {
	subject := "Tu as été invité comme témoin"
	titre := "Tu as été invité comme témoin"
	corps := fmt.Sprintf(`<p style="margin: 0 0 12px 0;">Salut %s,</p><p style="margin: 0;">Tu as été invité comme témoin sur la soirée « %s ».</p>`, toPseudo, soireeTitre)
	html := emailLayout(titre, corps)
	return n.send(ctx, toEmail, subject, html)
}

func (n *ResendNotifier) SendBadgeUnlocked(ctx context.Context, toEmail, toPseudo, badgeNom string) error {
	subject := "Nouveau badge débloqué !"
	titre := "Nouveau badge débloqué !"
	corps := fmt.Sprintf(`<p style="margin: 0 0 12px 0;">Bravo %s,</p><p style="margin: 0;">Tu as débloqué le badge « %s » !</p>`, toPseudo, badgeNom)
	html := emailLayout(titre, corps)
	return n.send(ctx, toEmail, subject, html)
}

func (n *ResendNotifier) SendFriendRequest(ctx context.Context, toEmail, toPseudo, fromPseudo string) error {
	subject := "Nouvelle demande d'ami"
	titre := "Nouvelle demande d'ami"
	corps := fmt.Sprintf(`<p style="margin: 0 0 12px 0;">Salut %s,</p><p style="margin: 0;">%s souhaite t'ajouter comme ami.</p>`, toPseudo, fromPseudo)
	html := emailLayout(titre, corps)
	return n.send(ctx, toEmail, subject, html)
}

func (n *ResendNotifier) SendReportResolved(ctx context.Context, toEmail, toPseudo string, temoignageSupprime bool) error {
	subject := "Ton témoignage signalé a été traité"
	titre := "Ton témoignage signalé a été traité"
	decision := "a été conservé"
	if temoignageSupprime {
		decision = "a été supprimé"
	}
	corps := fmt.Sprintf(`<p style="margin: 0 0 12px 0;">Salut %s,</p><p style="margin: 0;">Ton témoignage signalé %s après examen par un modérateur.</p>`, toPseudo, decision)
	html := emailLayout(titre, corps)
	return n.send(ctx, toEmail, subject, html)
}
