import { api } from './client'
import type { SignalementView } from './types'

// UC22 : traiter un signalement (réservé au rôle "moderator").
export const signalementsApi = {
  listEnAttente: (signal?: AbortSignal) =>
    api.get<SignalementView[]>('/signalements', signal),
  traiter: (signalementId: number, action: 'rejeter' | 'supprimer') =>
    api.post<void>(`/signalements/${signalementId}/traiter`, { action }),
}
