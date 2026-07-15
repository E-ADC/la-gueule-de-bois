package services

import "context"

// Notifier abstrait l'envoi de notifications par email. Les méthodes
// correspondent aux notifications décrites par les fiches UC09, UC14,
// UC21 et UC22 (acteur secondaire "Service d'email", prestataire Resend).
// Deux implémentations : ResendNotifier (appel HTTPS réel) et MockNotifier
// (tests, sans réseau).
type Notifier interface {
	// SendInvitation notifie un utilisateur invité comme témoin d'une
	// soirée (UC09).
	SendInvitation(ctx context.Context, toEmail, toPseudo, soireeTitre string) error
	// SendBadgeUnlocked notifie l'utilisateur qu'il a débloqué un badge
	// (UC14).
	SendBadgeUnlocked(ctx context.Context, toEmail, toPseudo, badgeNom string) error
	// SendFriendRequest notifie le destinataire d'une demande d'ami (UC21).
	SendFriendRequest(ctx context.Context, toEmail, toPseudo, fromPseudo string) error
	// SendReportResolved notifie l'auteur d'un témoignage signalé de la
	// décision du modérateur (UC22).
	SendReportResolved(ctx context.Context, toEmail, toPseudo string, temoignageSupprime bool) error
}
