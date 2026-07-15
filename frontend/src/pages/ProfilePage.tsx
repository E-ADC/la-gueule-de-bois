import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { badgesApi } from '../api/badges'
import { api, ApiError } from '../api/client'
import type { MyBadgesResponse, User } from '../api/types'
import { useAuth } from '../auth/AuthContext'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'
import { badgeIconByCode } from '../components/BadgeIcons'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; badges: MyBadgesResponse }

interface EditProfileState {
  pseudo: string
  bio: string
  avatar: string
}

/** UC04 (profil courant), UC15 (badges) et affichage du score (UC16 côté lecture). */
export function ProfilePage() {
  const { user } = useAuth()
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)

  // État pour le formulaire d'édition de profil
  const [displayedUser, setDisplayedUser] = useState<User | null>(user ?? null)
  const [editForm, setEditForm] = useState<EditProfileState>(() => ({
    pseudo: user?.pseudo ?? '',
    bio: user?.bio ?? '',
    avatar: user?.avatar ?? '',
  }))
  const [editError, setEditError] = useState<string | null>(null)
  const [editSuccess, setEditSuccess] = useState<string | null>(null)
  const [editSubmitting, setEditSubmitting] = useState(false)

  useEffect(() => {
    if (user) {
      setDisplayedUser(user)
      setEditForm({
        pseudo: user.pseudo,
        bio: user.bio ?? '',
        avatar: user.avatar ?? '',
      })
    }
  }, [user])

  useEffect(() => {
    const controller = new AbortController()
    setState({ status: 'loading' })
    badgesApi
      .mine(controller.signal)
      .then((badges) =>
        setState({
          status: 'ready',
          badges: { obtenus: badges.obtenus ?? [], tous: badges.tous ?? [] },
        }),
      )
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger les badges.',
        })
      })
    return () => controller.abort()
  }, [reloadToken])

  async function handleEditSubmit(event: FormEvent) {
    event.preventDefault()
    setEditError(null)
    setEditSuccess(null)
    setEditSubmitting(true)

    try {
      const updatedUser = await api.put<User>('/me', {
        pseudo: editForm.pseudo,
        bio: editForm.bio,
        avatar: editForm.avatar,
      })
      setDisplayedUser(updatedUser)
      setEditSuccess('Profil mis à jour.')
      // Effacer le message de succès après 3s
      setTimeout(() => setEditSuccess(null), 3000)
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.status === 409) {
          setEditError('Ce pseudo est déjà pris.')
        } else {
          setEditError(err.message)
        }
      } else {
        setEditError('Impossible de mettre à jour le profil.')
      }
    } finally {
      setEditSubmitting(false)
    }
  }

  if (!user || !displayedUser) {
    return <Loading label="Chargement du profil…" />
  }

  return (
    <div className="page">
      <div className="card profile-card">
        <div className="avatar profile-avatar">{displayedUser.pseudo.slice(0, 2).toUpperCase()}</div>
        <div>
          <p className="card-title">{displayedUser.pseudo}</p>
          <p className="card-meta">{displayedUser.email}</p>
          {displayedUser.bio && <p>{displayedUser.bio}</p>}
          <p className="label">Score</p>
          <p className="score-value">{displayedUser.score} pts</p>
        </div>
      </div>

      {/* Formulaire d'édition du profil (UC04) */}
      <div className="card">
        <h2>Modifier mon profil</h2>
        <form onSubmit={(event) => void handleEditSubmit(event)}>
          <label className="label" htmlFor="edit-pseudo">
            Pseudo
          </label>
          <input
            id="edit-pseudo"
            type="text"
            className="input"
            value={editForm.pseudo}
            onChange={(event) => setEditForm((f) => ({ ...f, pseudo: event.target.value }))}
            required
            minLength={3}
          />

          <label className="label" htmlFor="edit-bio">
            Bio
          </label>
          <textarea
            id="edit-bio"
            className="textarea"
            value={editForm.bio}
            onChange={(event) => setEditForm((f) => ({ ...f, bio: event.target.value }))}
            rows={3}
          />

          <label className="label" htmlFor="edit-avatar">
            Avatar (URL)
          </label>
          <input
            id="edit-avatar"
            type="text"
            className="input"
            value={editForm.avatar}
            onChange={(event) => setEditForm((f) => ({ ...f, avatar: event.target.value }))}
            placeholder="https://..."
          />

          {editError && <p className="field-error">{editError}</p>}
          {editSuccess && <p className="label">{editSuccess}</p>}

          <button type="submit" className="btn btn-primary" disabled={editSubmitting}>
            {editSubmitting ? 'Enregistrement…' : 'Enregistrer'}
          </button>
        </form>
      </div>

      <h2>Badges</h2>

      {state.status === 'loading' && <Loading label="Chargement des badges…" />}

      {state.status === 'error' && (
        <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
      )}

      {state.status === 'ready' && state.badges.obtenus.length === 0 && (
        <EmptyState
          title="Aucun badge"
          message="Aucun badge débloqué pour l’instant (débloqués automatiquement via le scoring, UC14)."
        />
      )}

      {state.status === 'ready' && state.badges.obtenus.length > 0 && (
        <ul className="badge-list">
          {state.badges.obtenus.map((badge) => (
            <li key={badge.id}>
              <span className="badge" title={badge.description}>
                {badgeIconByCode(badge.code)}
                {badge.nom}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
