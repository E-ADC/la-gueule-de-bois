import { api } from './client'
import type { User } from './types'

// UC05 : consulter le profil public d'un autre utilisateur.
// GET /users/{id} renvoie un User.
// TODO(UC05) : pas encore de page dédiée (route /utilisateurs/:id à ajouter).
export const usersApi = {
  publicProfile: (id: number | string, signal?: AbortSignal) =>
    api.get<User>(`/users/${id}`, signal),
}
