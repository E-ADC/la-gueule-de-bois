import { api } from './client'
import type { PublicProfile } from './types'

// UC05 : consulter le profil public d'un autre utilisateur.
// GET /users/{id} renvoie une enveloppe { user, badges }.
export const usersApi = {
  publicProfile: (id: number | string, signal?: AbortSignal) =>
    api.get<PublicProfile>(`/users/${id}`, signal),
}
