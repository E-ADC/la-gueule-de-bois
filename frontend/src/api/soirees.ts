import { api } from './client'
import type { Photo, Soiree, SoireeDetail } from './types'

export interface SoireeFormValues {
  titre: string
  lieu: string
  date: string
  description: string
  photos: File[]
}

// Le backend attend un JSON (date RFC3339) puis un upload séparé par photo
// (POST /soirees/{id}/photos, champ `photo`) — cf. backend/internal/handlers/soiree.go
function toJson(values: SoireeFormValues) {
  return {
    titre: values.titre,
    lieu: values.lieu,
    date: `${values.date}T00:00:00Z`,
    description: values.description,
  }
}

async function uploadPhotos(soireeId: number | string, photos: File[]): Promise<void> {
  for (const photo of photos) {
    const form = new FormData()
    form.set('photo', photo)
    await api.postForm<Photo>(`/soirees/${soireeId}/photos`, form)
  }
}

// UC10 : historique -> GET /soirees ; UC06/07/08 : création/modif/suppression.
export const soireesApi = {
  list: (signal?: AbortSignal) => api.get<Soiree[]>('/soirees', signal),
  // GET /soirees/{id} renvoie une enveloppe { soiree, photos }.
  get: (id: string, signal?: AbortSignal) => api.get<SoireeDetail>(`/soirees/${id}`, signal),
  create: async (values: SoireeFormValues) => {
    const soiree = await api.post<Soiree>('/soirees', toJson(values))
    await uploadPhotos(soiree.id, values.photos)
    return soiree
  },
  update: async (id: string, values: SoireeFormValues) => {
    const soiree = await api.put<Soiree>(`/soirees/${id}`, toJson(values))
    await uploadPhotos(id, values.photos)
    return soiree
  },
  remove: (id: string) => api.del<void>(`/soirees/${id}`),
}
