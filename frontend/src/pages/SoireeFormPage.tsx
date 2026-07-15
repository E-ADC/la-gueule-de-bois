import { useEffect, useState } from 'react'
import type { ChangeEvent, FormEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { soireesApi } from '../api/soirees'
import { ApiError } from '../api/client'
import { Loading, ErrorState } from '../components/StateViews'

/** Formulaire UC06 (création) / UC07 (modification), avec upload de photos (multipart). */
export function SoireeFormPage() {
  const { id } = useParams<{ id: string }>()
  const isEditing = Boolean(id)
  const navigate = useNavigate()

  const [titre, setTitre] = useState('')
  const [lieu, setLieu] = useState('')
  const [date, setDate] = useState('')
  const [description, setDescription] = useState('')
  const [photos, setPhotos] = useState<File[]>([])
  const [existingPhotoUrls, setExistingPhotoUrls] = useState<string[]>([])

  const [loading, setLoading] = useState(isEditing)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!id) return
    const controller = new AbortController()
    setLoading(true)
    soireesApi
      .get(id, controller.signal) // enveloppe { soiree, photos }
      .then(({ soiree, photos: existingPhotos }) => {
        setTitre(soiree.titre)
        setLieu(soiree.lieu)
        setDate(soiree.date.slice(0, 10))
        setDescription(soiree.description ?? '')
        setExistingPhotoUrls(existingPhotos.map((photo) => photo.path))
      })
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setLoadError(
          err instanceof ApiError ? err.message : 'Impossible de charger cette soirée.',
        )
      })
      .finally(() => setLoading(false))
    return () => controller.abort()
  }, [id])

  function handlePhotosChange(event: ChangeEvent<HTMLInputElement>) {
    setPhotos(event.target.files ? Array.from(event.target.files) : [])
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setSubmitError(null)
    setSubmitting(true)
    try {
      const values = { titre, lieu, date, description, photos }
      const soiree = isEditing && id ? await soireesApi.update(id, values) : await soireesApi.create(values)
      navigate(`/soirees/${soiree.id}`)
    } catch (err) {
      setSubmitError(
        err instanceof ApiError ? err.message : "L'enregistrement de la soirée a échoué.",
      )
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return <Loading label="Chargement de la soirée…" />
  }

  if (loadError) {
    return <ErrorState message={loadError} />
  }

  return (
    <div className="page">
      <h1>{isEditing ? 'Modifier la soirée' : 'Nouvelle soirée'}</h1>
      <form className="card soiree-form" onSubmit={(event) => void handleSubmit(event)}>
        <label className="label" htmlFor="titre">
          Titre
        </label>
        <input
          id="titre"
          className="input"
          value={titre}
          onChange={(event) => setTitre(event.target.value)}
          required
        />

        <label className="label" htmlFor="lieu">
          Lieu
        </label>
        <input
          id="lieu"
          className="input"
          value={lieu}
          onChange={(event) => setLieu(event.target.value)}
          required
        />

        <label className="label" htmlFor="date">
          Date
        </label>
        <input
          id="date"
          type="date"
          className="input"
          value={date}
          onChange={(event) => setDate(event.target.value)}
          required
        />

        <label className="label" htmlFor="description">
          Description
        </label>
        <textarea
          id="description"
          className="textarea"
          value={description}
          onChange={(event) => setDescription(event.target.value)}
          rows={4}
        />

        <label className="label" htmlFor="photos">
          Photos (jpeg, png, webp — 10 Mo max chacune)
        </label>
        <input
          id="photos"
          type="file"
          className="input"
          accept="image/jpeg,image/png,image/webp"
          multiple
          onChange={handlePhotosChange}
        />
        {photos.length > 0 && (
          <p className="label">{photos.length} nouvelle(s) photo(s) sélectionnée(s)</p>
        )}
        {existingPhotoUrls.length > 0 && (
          <ul className="photo-preview-list">
            {existingPhotoUrls.map((url) => (
              <li key={url}>
                <img src={url} alt="Photo de la soirée" />
              </li>
            ))}
          </ul>
        )}

        {submitError && <p className="field-error">{submitError}</p>}

        <button type="submit" className="btn btn-primary" disabled={submitting}>
          {submitting ? 'Enregistrement…' : isEditing ? 'Enregistrer' : 'Créer la soirée'}
        </button>
      </form>
    </div>
  )
}
