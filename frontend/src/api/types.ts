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

// TODO(backend) : DTO enrichi (pseudo de l'auteur, compteurs de votes,
// vote de l'utilisateur courant) — en attendant, le front affiche les
// témoignages sans compteurs.
export interface Temoignage {
  id: number
  soireeId: number
  auteurId: number
  contenu: string
  createdAt: string
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

/** Déblocage d'un badge par un utilisateur (table de liaison). */
export interface UserBadge {
  userId: number
  badgeId: number
  debloqueLe: string
}

/** Enveloppe renvoyée par GET /api/me/badges. */
export interface MyBadgesResponse {
  obtenus: UserBadge[]
  tous: Badge[]
}

export interface Groupe {
  id: number
  nom: string
  createurId: number
  createdAt: string
}

export type SignalementStatut = 'en_attente' | 'rejete' | 'temoignage_supprime'

export interface Signalement {
  id: number
  temoignageId: number
  auteurId: number
  motif: string
  statut: SignalementStatut
  traiteParId?: number
  createdAt: string
  traiteLe?: string
}

/** Format d'erreur uniforme renvoyé par l'API (spec §API — gestion d'erreurs). */
export interface ApiErrorBody {
  error: string
  code: string
}
