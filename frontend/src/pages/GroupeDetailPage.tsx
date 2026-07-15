import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { groupesApi } from '../api/groupes'
import { classementApi } from '../api/classement'
import { ApiError } from '../api/client'
import type { Groupe, User } from '../api/types'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; groupe: Groupe; membres: User[] }

/** Détail d'un groupe (UC18/19) : infos + classement de ses membres (UC20). */
export function GroupeDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)

  useEffect(() => {
    if (!id) return
    const controller = new AbortController()
    setState({ status: 'loading' })
    Promise.all([groupesApi.get(id, controller.signal), classementApi.groupe(id, controller.signal)])
      .then(([groupe, membres]) => setState({ status: 'ready', groupe, membres: membres ?? [] }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger ce groupe.',
        })
      })
    return () => controller.abort()
  }, [id, reloadToken])

  if (state.status === 'loading') {
    return <Loading label="Chargement du groupe…" />
  }

  if (state.status === 'error') {
    return <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
  }

  const { groupe, membres } = state

  return (
    <div className="page">
      <div className="page-header">
        <div>
          <h1>{groupe.nom}</h1>
          <p className="card-meta">
            Groupe n°{groupe.id} · créé le {new Date(groupe.createdAt).toLocaleDateString('fr-FR')}
          </p>
        </div>
      </div>

      <div className="card">
        <p className="label">Pour inviter quelqu'un à rejoindre ce groupe, partage cet ID :</p>
        <p className="card-title">{groupe.id}</p>
      </div>

      <hr className="sep" />
      <h2>Membres</h2>

      {membres.length === 0 ? (
        <EmptyState title="Aucun membre" message="Ce groupe n'a pas encore de membre." />
      ) : (
        <div className="card">
          {membres.map((membre, index) => (
            <div key={membre.id} className="rank-row">
              <span className="rank">{index + 1}</span>
              <span className="avatar">{membre.pseudo.slice(0, 2).toUpperCase()}</span>
              <span>{membre.pseudo}</span>
              <span className="pts">{membre.score} pts</span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
