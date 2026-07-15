import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { usersApi } from '../api/users'
import { ApiError } from '../api/client'
import type { PublicProfile, Badge } from '../api/types'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'
import { badgeIconByCode } from '../components/BadgeIcons'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; profile: PublicProfile }

/** UC05 : consulter le profil public d'un autre utilisateur. */
export function PublicProfilePage() {
  const { id } = useParams<{ id: string }>()
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)

  useEffect(() => {
    if (!id) return
    const controller = new AbortController()
    setState({ status: 'loading' })
    usersApi
      .publicProfile(id, controller.signal)
      .then((publicProfile) => setState({ status: 'ready', profile: publicProfile }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger le profil.',
        })
      })
    return () => controller.abort()
  }, [id, reloadToken])

  if (state.status === 'loading') {
    return <Loading label="Chargement du profil…" />
  }

  if (state.status === 'error') {
    return (
      <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
    )
  }

  const { user, badges } = state.profile

  return (
    <div className="page">
      <div className="card profile-card">
        <div className="avatar profile-avatar">{user.pseudo.slice(0, 2).toUpperCase()}</div>
        <div>
          <p className="card-title">{user.pseudo}</p>
          {user.bio && <p>{user.bio}</p>}
          <p className="label">Score</p>
          <p className="score-value">{user.score} pts</p>
        </div>
      </div>

      <h2>Badges obtenus</h2>
      {badges.length === 0 ? (
        <EmptyState
          title="Aucun badge"
          message={`${user.pseudo} n'a pas encore obtenu de badge.`}
        />
      ) : (
        <ul className="badge-list">
          {badges.map((badge: Badge) => (
            <li key={badge.id}>
              <span className="badge" title={`${badge.description} (seuil : ${badge.seuilScore} pts)`}>
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
