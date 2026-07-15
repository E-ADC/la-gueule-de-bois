import { api } from './client'
import type { DemandeAmi, User } from './types'

// UC21 : envoyer une demande d'ami (par pseudo), lister les demandes reçues,
// y répondre, lister ses amis (demandes acceptées).
export const amisApi = {
  envoyer: (pseudo: string) => api.post<void>('/amis/demandes', { pseudo }),
  listRecues: (signal?: AbortSignal) => api.get<DemandeAmi[]>('/amis/demandes', signal),
  repondre: (demandeId: number, action: 'accepter' | 'refuser') =>
    api.post<void>(`/amis/demandes/${demandeId}/repondre`, { action }),
  listAmis: (signal?: AbortSignal) => api.get<User[]>('/amis', signal),
}
