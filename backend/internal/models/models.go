// Package models définit les entités métier partagées par les couches
// repository, services et handlers. Alignées sur le diagramme de classes
// et les fiches de cas d'utilisation (UC01 à UC22).
package models

import "time"

// User représente un utilisateur inscrit (UC01/UC02/UC04/UC05).
type User struct {
	ID           int64  `json:"id"`
	Pseudo       string `json:"pseudo"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Avatar       string `json:"avatar,omitempty"`
	Bio          string `json:"bio,omitempty"`
	Score        int    `json:"score"`
	// Role distingue "user" (défaut) et "moderator" (UC22). La fiche ne
	// détaille pas la gestion des rôles : choix simple, un champ texte.
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// Session matérialise une session cookie HttpOnly (auth sans JWT).
type Session struct {
	Token     string    `json:"-"`
	UserID    int64     `json:"-"`
	ExpiresAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// Soiree représente une soirée enregistrée par un utilisateur (UC06/07/08/10).
type Soiree struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Titre       string    `json:"titre"`
	DateSoiree  time.Time `json:"date"`
	Lieu        string    `json:"lieu"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Photo est associée à une soirée (association 1..* — spec "photos minimales").
type Photo struct {
	ID        int64     `json:"id"`
	SoireeID  int64     `json:"soireeId"`
	Chemin    string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
}

// TemoinInvitation trace l'invitation d'un témoin sur une soirée (UC09).
// Nécessaire pour vérifier la pré-condition de UC11 ("l'utilisateur est
// témoin invité de la soirée") : la fiche ne nomme pas explicitement cette
// table, déduite du besoin métier.
type TemoinInvitation struct {
	ID        int64     `json:"id"`
	SoireeID  int64     `json:"soireeId"`
	InviteID  int64     `json:"inviteId"`
	CreatedAt time.Time `json:"createdAt"`
}

// Temoignage est rédigé par un témoin invité sur une soirée (UC11).
type Temoignage struct {
	ID        int64     `json:"id"`
	SoireeID  int64     `json:"soireeId"`
	AuteurID  int64     `json:"auteurId"`
	Contenu   string    `json:"contenu"`
	CreatedAt time.Time `json:"createdAt"`
}

// VoteValeur énumère les deux sens de vote possibles (UC12, swipe positif/négatif).
type VoteValeur int

const (
	VotePositif VoteValeur = 1
	VoteNegatif VoteValeur = -1
)

// Vote est un swipe sur un témoignage, unique par (temoignage, utilisateur).
type Vote struct {
	ID           int64      `json:"id"`
	TemoignageID int64      `json:"temoignageId"`
	UserID       int64      `json:"userId"`
	Valeur       VoteValeur `json:"valeur"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// Badge décrit un badge débloquable (UC14/UC15). Le seuil est exprimé en
// score : règle la plus simple qui reste alignée avec "évalue les critères
// de badges avec le score mis à jour" (fiche UC14, pas d'autre critère
// fourni).
type Badge struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Nom         string `json:"nom"`
	Description string `json:"description"`
	SeuilScore  int    `json:"seuilScore"`
}

// UserBadge trace le déblocage d'un badge par un utilisateur (table de liaison).
type UserBadge struct {
	UserID     int64     `json:"userId"`
	BadgeID    int64     `json:"badgeId"`
	DebloqueLe time.Time `json:"debloqueLe"`
}

// Groupe est un groupe d'amis (UC18/UC19/UC20).
type Groupe struct {
	ID         int64     `json:"id"`
	Nom        string    `json:"nom"`
	CreateurID int64     `json:"createurId"`
	CreatedAt  time.Time `json:"createdAt"`
}

// GroupeMembre est la table de liaison groupe <-> utilisateur.
type GroupeMembre struct {
	GroupeID int64     `json:"groupeId"`
	UserID   int64     `json:"userId"`
	JoinedAt time.Time `json:"joinedAt"`
}

// StatutSignalement énumère le cycle de vie d'un signalement (UC13/UC22).
type StatutSignalement string

const (
	SignalementEnAttente StatutSignalement = "en_attente"
	SignalementRejete    StatutSignalement = "rejete"
	SignalementSupprime  StatutSignalement = "temoignage_supprime"
)

// Signalement est créé par un utilisateur sur un témoignage jugé
// inapproprié (UC13), puis traité par un modérateur (UC22).
type Signalement struct {
	ID           int64             `json:"id"`
	TemoignageID int64             `json:"temoignageId"`
	AuteurID     int64             `json:"auteurId"`
	Motif        string            `json:"motif"`
	Statut       StatutSignalement `json:"statut"`
	TraiteParID  *int64            `json:"traiteParId,omitempty"`
	CreatedAt    time.Time         `json:"createdAt"`
	TraiteLe     *time.Time        `json:"traiteLe,omitempty"`
}

// StatutDemandeAmi énumère le cycle de vie d'une demande d'ami (UC21).
type StatutDemandeAmi string

const (
	DemandeAmiEnAttente StatutDemandeAmi = "en_attente"
	DemandeAmiAcceptee  StatutDemandeAmi = "acceptee"
	DemandeAmiRefusee   StatutDemandeAmi = "refusee"
)

// DemandeAmi représente une demande d'ami envoyée entre deux utilisateurs (UC21).
// Handler exposé en TODO dans ce squelette (hors périmètre "premières
// couches" demandées), le modèle et la migration sont prêts.
type DemandeAmi struct {
	ID             int64            `json:"id"`
	DemandeurID    int64            `json:"demandeurId"`
	DestinataireID int64            `json:"destinataireId"`
	Statut         StatutDemandeAmi `json:"statut"`
	CreatedAt      time.Time        `json:"createdAt"`
}
