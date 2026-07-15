import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { soireesApi } from '../api/soirees'
import { temoignagesApi } from '../api/temoignages'
import { ApiError } from '../api/client'
import type { Photo, Soiree, Temoignage } from '../api/types'
import { useAuth } from '../auth/AuthContext'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; soiree: Soiree; photos: Photo[]; temoignages: Temoignage[] }

export function SoireeDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const navigate = useNavigate()
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)
  const [nouveauTemoignage, setNouveauTemoignage] = useState('')
  const [publishError, setPublishError] = useState<string | null>(null)
  const [publishing, setPublishing] = useState(false)
  const [voteMessage, setVoteMessage] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    const controller = new AbortController()
    setState({ status: 'loading' })
    Promise.all([
      soireesApi.get(id, controller.signal), // enveloppe { soiree, photos }
      temoignagesApi.listBySoiree(id, controller.signal),
    ])
      .then(([detail, temoignages]) =>
        // Filet de sécurité : une liste vide côté API peut arriver en JSON `null`
        // (slice Go non initialisée) plutôt qu'en `[]` — normalisé ici pour ne
        // jamais planter le rendu (.length/.map sur null).
        setState({
          status: 'ready',
          soiree: detail.soiree,
          photos: detail.photos ?? [],
          temoignages: temoignages ?? [],
        }),
      )
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger cette soirée.',
        })
      })
    return () => controller.abort()
  }, [id, reloadToken])

  async function handleDelete() {
    if (!id) return
    if (!window.confirm('Supprimer définitivement cette soirée ?')) return
    try {
      await soireesApi.remove(id)
      navigate('/')
    } catch (err) {
      window.alert(err instanceof ApiError ? err.message : 'La suppression a échoué.')
    }
  }

  async function handlePublishTemoignage(event: FormEvent) {
    event.preventDefault()
    if (!id || !nouveauTemoignage.trim()) return
    setPublishError(null)
    setPublishing(true)
    try {
      await temoignagesApi.create(id, nouveauTemoignage.trim())
      setNouveauTemoignage('')
      setReloadToken((token) => token + 1)
    } catch (err) {
      setPublishError(
        err instanceof ApiError ? err.message : 'La publication du témoignage a échoué.',
      )
    } finally {
      setPublishing(false)
    }
  }

  async function handleVote(temoignageId: number, valeur: 1 | -1) {
    setVoteMessage(null)
    try {
      await temoignagesApi.vote(temoignageId, valeur)
      setVoteMessage('Vote enregistré.')
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setVoteMessage('Tu as déjà voté sur ce témoignage.')
      } else {
        setVoteMessage(err instanceof ApiError ? err.message : 'Le vote a échoué.')
      }
    }
  }

  if (state.status === 'loading') {
    return <Loading label="Chargement de la soirée…" />
  }

  if (state.status === 'error') {
    return <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
  }

  const { soiree, photos, temoignages } = state
  const estAuteur = user?.id === soiree.userId

  return (
    <div className="page">
      <div className="page-header">
        <div>
          <h1>{soiree.titre}</h1>
          <p className="card-meta">
            {soiree.lieu} · {new Date(soiree.date).toLocaleDateString('fr-FR')}
          </p>
        </div>
        {estAuteur && (
          <div className="header-actions">
            <Link to={`/soirees/${soiree.id}/modifier`} className="btn btn-ghost">
              Modifier
            </Link>
            <button type="button" className="btn btn-danger" onClick={() => void handleDelete()}>
              Supprimer
            </button>
          </div>
        )}
      </div>

      {soiree.description && <p>{soiree.description}</p>}

      {photos.length > 0 && (
        <ul className="photo-preview-list">
          {photos.map((photo) => (
            <li key={photo.id}>
              <img src={photo.path} alt={`Photo de ${soiree.titre}`} />
            </li>
          ))}
        </ul>
      )}

      <hr className="sep" />

      <h2>Témoignages</h2>

      <form className="card" onSubmit={(event) => void handlePublishTemoignage(event)}>
        <label className="label" htmlFor="temoignage">
          Ajouter un témoignage
        </label>
        <textarea
          id="temoignage"
          className="textarea"
          rows={3}
          value={nouveauTemoignage}
          onChange={(event) => setNouveauTemoignage(event.target.value)}
          placeholder="Raconte ce qu'il s'est passé…"
        />
        {publishError && <p className="field-error">{publishError}</p>}
        <button type="submit" className="btn btn-primary" disabled={publishing}>
          {publishing ? 'Publication…' : 'Publier'}
        </button>
      </form>

      {voteMessage && <p className="label vote-message">{voteMessage}</p>}

      {temoignages.length === 0 ? (
        <EmptyState
          title="Aucun témoignage"
          message="Personne n'a encore témoigné sur cette soirée."
        />
      ) : (
        <ul className="temoignage-list">
          {/* TODO(backend) : DTO enrichi côté backend (pseudo de l'auteur,
              compteurs de votes, vote de l'utilisateur courant) — en attendant,
              affichage sans compteurs ni état "déjà voté". */}
          {temoignages.map((temoignage) => (
            <li key={temoignage.id} className="card temoignage-card">
              <p className="card-meta">
                {new Date(temoignage.createdAt).toLocaleDateString('fr-FR')}
              </p>
              <p>{temoignage.contenu}</p>
              <div className="vote-row">
                <button
                  type="button"
                  className="btn btn-ghost"
                  onClick={() => void handleVote(temoignage.id, 1)}
                >
                  +1
                </button>
                <button
                  type="button"
                  className="btn btn-ghost"
                  onClick={() => void handleVote(temoignage.id, -1)}
                >
                  −1
                </button>
              </div>
              {/* TODO(UC13) : bouton "signaler" à brancher sur temoignagesApi.signaler
                  (route encore en 501 côté backend) */}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
