import { api } from './client'
import type { Groupe } from './types'

// UC18 : créer un groupe ; UC19 : rejoindre un groupe.
// TODO(UC18/UC19) : routes encore en 501 côté backend (et GET /groupes
// n'existe pas encore : le sélecteur de groupe du classement encaisse
// silencieusement l'erreur). Pas de page dédiée pour l'instant.
export const groupesApi = {
  listMine: (signal?: AbortSignal) => api.get<Groupe[]>('/groupes', signal),
  create: (nom: string) => api.post<Groupe>('/groupes', { nom }),
  join: (groupeId: number | string) => api.post<void>(`/groupes/${groupeId}/membres`),
}
