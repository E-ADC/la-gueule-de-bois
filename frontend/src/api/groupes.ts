import { api } from './client'
import type { Groupe } from './types'

// UC18 : créer un groupe ; UC19 : rejoindre un groupe.
export const groupesApi = {
  listMine: (signal?: AbortSignal) => api.get<Groupe[]>('/groupes', signal),
  get: (groupeId: number | string, signal?: AbortSignal) =>
    api.get<Groupe>(`/groupes/${groupeId}`, signal),
  create: (nom: string) => api.post<Groupe>('/groupes', { nom }),
  join: (groupeId: number | string) => api.post<void>(`/groupes/${groupeId}/membres`),
}
