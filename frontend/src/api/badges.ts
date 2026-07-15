import { api } from './client'
import type { MyBadgesResponse } from './types'

// UC15 : consulter ses badges. GET /me/badges renvoie une enveloppe
// { obtenus: UserBadge[], tous: Badge[] } — l'état "obtenu / à débloquer"
// se reconstruit côté front en croisant les deux listes.
export const badgesApi = {
  mine: (signal?: AbortSignal) => api.get<MyBadgesResponse>('/me/badges', signal),
}
