import { useEffect, useState } from 'react'
import { badgesApi } from '../api/badges'
import { ApiError } from '../api/client'
import type { MyBadgesResponse } from '../api/types'
import { useAuth } from '../auth/AuthContext'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; badges: MyBadgesResponse }

/** UC04 (profil courant), UC15 (badges) et affichage du score (UC16 côté lecture). */
export function ProfilePage() {
  const { user } = useAuth()
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)

  useEffect(() => {
    const controller = new AbortController()
    setState({ status: 'loading' })
    badgesApi
      .mine(controller.signal)
      .then((badges) => setState({ status: 'ready', badges }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger les badges.',
        })
      })
    return () => controller.abort()
  }, [reloadToken])

  if (!user) {
    return <Loading label="Chargement du profil…" />
  }

  // GET /me/badges renvoie { obtenus: Badge[], tous: Badge[] } :
  // l'état "obtenu / à débloquer" se reconstruit en croisant les deux listes.
  const obtenusIds =
    state.status === 'ready'
      ? new Set(state.badges.obtenus.map((badge) => badge.id))
      : new Set<number>()

  return (
    <div className="page">
      <div className="card profile-card">
        <div className="avatar profile-avatar">{user.pseudo.slice(0, 2).toUpperCase()}</div>
        <div>
          <p className="card-title">{user.pseudo}</p>
          <p className="card-meta">{user.email}</p>
          {user.bio && <p>{user.bio}</p>}
          <p className="label">Score</p>
          <p className="score-value">{user.score} pts</p>
        </div>
      </div>

      {/* TODO(UC04) : formulaire d'édition du profil (pseudo, avatar, bio) à ajouter ici */}

      <h2>Badges</h2>

      {state.status === 'loading' && <Loading label="Chargement des badges…" />}

      {state.status === 'error' && (
        <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
      )}

      {state.status === 'ready' && state.badges.tous.length === 0 && (
        <EmptyState
          title="Aucun badge"
          message="Aucun badge à afficher pour l’instant (débloqués automatiquement via le scoring, UC14)."
        />
      )}

      {state.status === 'ready' && state.badges.tous.length > 0 && (
        <ul className="badge-list">
          {state.badges.tous.map((badge) => (
            <li key={badge.id}>
              <span
                className={obtenusIds.has(badge.id) ? 'badge' : 'badge badge-locked'}
                title={`${badge.description} (seuil : ${badge.seuilScore} pts)`}
              >
                {badge.nom}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
