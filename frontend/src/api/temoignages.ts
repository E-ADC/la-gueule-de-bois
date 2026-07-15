import { api } from './client'
import type { Temoignage, Vote, VoteValeur } from './types'

// UC09 : inviter un témoin ; UC11 : ajouter un témoignage ; UC12 : voter ;
// UC13 : signaler.
export const temoignagesApi = {
  listBySoiree: (soireeId: number | string, signal?: AbortSignal) =>
    api.get<Temoignage[]>(`/soirees/${soireeId}/temoignages`, signal),
  create: (soireeId: number | string, contenu: string) =>
    api.post<Temoignage>(`/soirees/${soireeId}/temoignages`, { contenu }),
  // UC09 : invite par email (l'invité doit déjà avoir un compte).
  inviteTemoin: (soireeId: number | string, email: string) =>
    api.post<void>(`/soirees/${soireeId}/temoins`, { email }),
  // Renvoie le Vote créé ; 409 si l'utilisateur a déjà voté sur ce témoignage.
  vote: (temoignageId: number, valeur: VoteValeur) =>
    api.post<Vote>(`/temoignages/${temoignageId}/votes`, { valeur }),
  // UC13 : signaler un témoignage (409 si déjà signalé par cet utilisateur).
  signaler: (temoignageId: number, motif: string) =>
    api.post<void>(`/temoignages/${temoignageId}/signalements`, { motif }),
}
