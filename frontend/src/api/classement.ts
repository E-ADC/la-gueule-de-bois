import { api } from './client'
import type { User } from './types'

// UC17 : classement global ; UC20 : classement restreint à un groupe.
// Le backend renvoie des User[] triés par score décroissant (pas de champ
// rang : il se déduit de l'index côté front).
export const classementApi = {
  global: (signal?: AbortSignal) => api.get<User[]>('/classement', signal),
  groupe: (groupeId: number | string, signal?: AbortSignal) =>
    api.get<User[]>(`/groupes/${groupeId}/classement`, signal),
}
