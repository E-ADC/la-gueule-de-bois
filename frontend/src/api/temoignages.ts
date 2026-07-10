import { api } from './client'
import type { Temoignage, Vote, VoteValeur } from './types'

// UC09 : inviter un témoin ; UC11 : ajouter un témoignage ; UC12 : voter ;
// UC13 : signaler.
export const temoignagesApi = {
  listBySoiree: (soireeId: number | string, signal?: AbortSignal) =>
    api.get<Temoignage[]>(`/soirees/${soireeId}/temoignages`, signal),
  create: (soireeId: number | string, contenu: string) =>
    api.post<Temoignage>(`/soirees/${soireeId}/temoignages`, { contenu }),
  // TODO(UC09) : pas encore d'UI, client prêt (POST /soirees/{id}/temoins).
  inviteTemoin: (soireeId: number | string, inviteId: number) =>
    api.post<void>(`/soirees/${soireeId}/temoins`, { inviteId }),
  // Renvoie le Vote créé ; 409 si l'utilisateur a déjà voté sur ce témoignage.
  vote: (temoignageId: number, valeur: VoteValeur) =>
    api.post<Vote>(`/temoignages/${temoignageId}/votes`, { valeur }),
  // TODO(UC13) : route encore en 501 côté backend.
  signaler: (temoignageId: number, motif: string) =>
    api.post<void>(`/temoignages/${temoignageId}/signalements`, { motif }),
}
