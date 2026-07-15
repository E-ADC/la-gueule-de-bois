/**
 * Types des entités métier, alignés sur le contrat RÉEL du backend :
 * `backend/internal/models/models.go` + `backend/internal/handlers/`.
 *
 * Tous les ids sont des nombres (int64 JSON). Les champs `omitempty` côté Go
 * sont optionnels ici.
 */

export interface User {
  id: number
  pseudo: string
  email: string
  avatar?: string
  bio?: string
  score: number
  role: string
  createdAt: string
}

export interface Photo {
  id: number
  soireeId: number
  /** Chemin public de l'image, ex. "/uploads/xxx.png" (servi en statique). */
  path: string
  createdAt: string
}

export interface Soiree {
  id: number
  userId: number
  titre: string
  date: string
  lieu: string
  description?: string
  createdAt: string
  updatedAt: string
}

/** Enveloppe renvoyée par GET /api/soirees/{id}. */
export interface SoireeDetail {
  soiree: Soiree
  photos: Photo[]
}

export interface Temoignage {
  id: number
  soireeId: number
  auteurId: number
  contenu: string
  createdAt: string
  auteurPseudo: string
  votesPositifs: number
  votesNegatifs: number
  /** Vote de l'utilisateur courant sur ce témoignage, `null` s'il n'a pas voté. */
  monVote: VoteValeur | null
}

export type VoteValeur = 1 | -1

export interface Vote {
  id: number
  temoignageId: number
  userId: number
  valeur: VoteValeur
  createdAt: string
}

export interface Badge {
  id: number
  code: string
  nom: string
  description: string
  seuilScore: number
}

/** Enveloppe renvoyée par GET /api/me/badges (`obtenus` = Badge complets). */
export interface MyBadgesResponse {
  obtenus: Badge[]
  tous: Badge[]
}

export interface Groupe {
  id: number
  nom: string
  createurId: number
  createdAt: string
}

/** Enveloppe renvoyée par GET /api/users/{id} (UC05). */
export interface PublicProfile {
  user: User
  badges: Badge[]
}

export type StatutDemandeAmi = 'en_attente' | 'acceptee' | 'refusee'

/** Demande d'ami (UC21). */
export interface DemandeAmi {
  id: number
  demandeurId: number
  destinataireId: number
  statut: StatutDemandeAmi
  createdAt: string
}

export type SignalementStatut = 'en_attente' | 'rejete' | 'temoignage_supprime'

export interface Signalement {
  id: number
  // `null` une fois le témoignage supprimé par un modérateur (UC22) — le
  // signalement survit comme historique de modération.
  temoignageId: number | null
  auteurId: number
  motif: string
  statut: SignalementStatut
  traiteParId?: number
  createdAt: string
  traiteLe?: string
}

/** Signalement enrichi du contenu du témoignage signalé (GET /signalements, UC22). */
export interface SignalementView extends Signalement {
  temoignageContenu: string
}

/** Format d'erreur uniforme renvoyé par l'API (spec §API — gestion d'erreurs). */
export interface ApiErrorBody {
  error: string
  code: string
}
