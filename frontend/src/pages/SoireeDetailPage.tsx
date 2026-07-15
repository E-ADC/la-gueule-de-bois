import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { soireesApi } from '../api/soirees'
import { temoignagesApi } from '../api/temoignages'
import { ApiError } from '../api/client'
import type { Photo, Soiree, Temoignage } from '../api/types'
import { useAuth } from '../auth/AuthContext'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'
import { ThumbUpIcon, ThumbDownIcon } from '../components/VoteIcons'

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
  const [reportMessage, setReportMessage] = useState<string | null>(null)
  const [invitePseudo, setInvitePseudo] = useState('')
  const [inviteMessage, setInviteMessage] = useState<string | null>(null)
  const [inviting, setInviting] = useState(false)
  const [lightboxPhoto, setLightboxPhoto] = useState<Photo | null>(null)

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

  // Ferme la visionneuse plein écran avec Échap.
  useEffect(() => {
    if (!lightboxPhoto) return
    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') setLightboxPhoto(null)
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [lightboxPhoto])

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
      if (err instanceof ApiError && err.status === 403) {
        setPublishError(
          "Tu dois être invité comme témoin par l'auteur de cette soirée pour pouvoir témoigner (UC09/UC11).",
        )
      } else {
        setPublishError(
          err instanceof ApiError ? err.message : 'La publication du témoignage a échoué.',
        )
      }
    } finally {
      setPublishing(false)
    }
  }

  async function handleVote(temoignageId: number, valeur: 1 | -1) {
    setVoteMessage(null)
    try {
      await temoignagesApi.vote(temoignageId, valeur)
      setReloadToken((token) => token + 1)
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setVoteMessage('Tu as déjà voté sur ce témoignage.')
      } else {
        setVoteMessage(err instanceof ApiError ? err.message : 'Le vote a échoué.')
      }
    }
  }

  async function handleSignaler(temoignageId: number) {
    const motif = window.prompt('Pourquoi signaler ce témoignage ?')
    if (!motif || !motif.trim()) return
    setReportMessage(null)
    try {
      await temoignagesApi.signaler(temoignageId, motif.trim())
      setReportMessage('Témoignage signalé, un modérateur va l’examiner (UC13).')
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setReportMessage('Tu as déjà signalé ce témoignage.')
      } else {
        setReportMessage(err instanceof ApiError ? err.message : 'Le signalement a échoué.')
      }
    }
  }

  async function handleInvite(event: FormEvent) {
    event.preventDefault()
    if (!id || !invitePseudo.trim()) return
    setInviteMessage(null)
    setInviting(true)
    try {
      await temoignagesApi.inviteTemoin(id, invitePseudo.trim())
      setInviteMessage(`${invitePseudo.trim()} a été invité comme témoin.`)
      setInvitePseudo('')
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setInviteMessage("Aucun compte n'existe avec ce pseudo.")
      } else if (err instanceof ApiError && err.status === 409) {
        setInviteMessage('Cette personne est déjà invitée comme témoin.')
      } else {
        setInviteMessage(err instanceof ApiError ? err.message : "L'invitation a échoué.")
      }
    } finally {
      setInviting(false)
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
              <button
                type="button"
                className="photo-preview-button"
                onClick={() => setLightboxPhoto(photo)}
                aria-label="Agrandir la photo"
              >
                <img src={photo.path} alt={`Photo de ${soiree.titre}`} />
              </button>
            </li>
          ))}
        </ul>
      )}

      {lightboxPhoto && (
        <div
          className="lightbox-overlay"
          onClick={() => setLightboxPhoto(null)}
          role="dialog"
          aria-modal="true"
          aria-label="Photo en plein écran"
        >
          <img
            src={lightboxPhoto.path}
            alt={`Photo de ${soiree.titre} en plein écran`}
            className="lightbox-image"
          />
          <button
            type="button"
            className="lightbox-close"
            onClick={() => setLightboxPhoto(null)}
            aria-label="Fermer"
          >
            ✕
          </button>
        </div>
      )}

      {estAuteur && (
        <>
          <hr className="sep" />
          <h2>Inviter un témoin</h2>
          <form className="card" onSubmit={(event) => void handleInvite(event)}>
            <label className="label" htmlFor="invite-pseudo">
              Pseudo de la personne à inviter (UC09)
            </label>
            <input
              id="invite-pseudo"
              type="text"
              className="input"
              value={invitePseudo}
              onChange={(event) => setInvitePseudo(event.target.value)}
              placeholder="pseudo"
              required
            />
            {inviteMessage && <p className="label vote-message">{inviteMessage}</p>}
            <button type="submit" className="btn btn-primary" disabled={inviting}>
              {inviting ? 'Invitation…' : 'Inviter'}
            </button>
          </form>
        </>
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
      {reportMessage && <p className="label vote-message">{reportMessage}</p>}

      {temoignages.length === 0 ? (
        <EmptyState
          title="Aucun témoignage"
          message="Personne n'a encore témoigné sur cette soirée."
        />
      ) : (
        <ul className="temoignage-list">
          {temoignages.map((temoignage) => (
            <li key={temoignage.id} className="card temoignage-card">
              <div className="temoignage-header">
                <span className="avatar">
                  {temoignage.auteurAvatar ? (
                    <img src={temoignage.auteurAvatar} alt={temoignage.auteurPseudo} />
                  ) : (
                    temoignage.auteurPseudo.slice(0, 2).toUpperCase()
                  )}
                </span>
                <div>
                  <p className="card-meta">
                    <strong>{temoignage.auteurPseudo}</strong> ·{' '}
                    {new Date(temoignage.createdAt).toLocaleDateString('fr-FR')}
                  </p>
                </div>
              </div>
              <p>{temoignage.contenu}</p>
              <div className="vote-row">
                <button
                  type="button"
                  className={`vote-btn${temoignage.monVote === 1 ? ' vote-btn-active' : ''}`}
                  disabled={temoignage.monVote !== null}
                  onClick={() => void handleVote(temoignage.id, 1)}
                  aria-label="Voter positivement"
                >
                  <ThumbUpIcon /> {temoignage.votesPositifs}
                </button>
                <button
                  type="button"
                  className={`vote-btn${temoignage.monVote === -1 ? ' vote-btn-active' : ''}`}
                  disabled={temoignage.monVote !== null}
                  onClick={() => void handleVote(temoignage.id, -1)}
                  aria-label="Voter négativement"
                >
                  <ThumbDownIcon /> {temoignage.votesNegatifs}
                </button>
                <button
                  type="button"
                  className="btn btn-ghost"
                  onClick={() => void handleSignaler(temoignage.id)}
                >
                  Signaler
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
