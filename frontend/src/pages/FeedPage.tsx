import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { soireesApi } from '../api/soirees'
import { ApiError } from '../api/client'
import type { Soiree } from '../api/types'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; soirees: Soiree[] }

export function FeedPage() {
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)

  useEffect(() => {
    const controller = new AbortController()
    setState({ status: 'loading' })
    soireesApi
      .list(controller.signal)
      .then((soirees) => setState({ status: 'ready', soirees }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Le chargement des soirées a échoué.',
        })
      })
    return () => controller.abort()
  }, [reloadToken])

  return (
    <div className="page">
      <div className="page-header">
        {/* GET /soirees ne renvoie que les soirées de l'utilisateur connecté (UC10). */}
        <h1>Mes soirées</h1>
        <Link to="/soirees/nouvelle" className="btn btn-primary">
          Nouvelle soirée
        </Link>
      </div>

      {state.status === 'loading' && <Loading label="Chargement des soirées…" />}

      {state.status === 'error' && (
        <ErrorState message={state.message} onRetry={() => setReloadToken((token) => token + 1)} />
      )}

      {state.status === 'ready' && state.soirees.length === 0 && (
        <EmptyState
          title="Aucune soirée pour l’instant"
          message="Tu n’as encore enregistré aucune soirée."
          action={
            <Link to="/soirees/nouvelle" className="btn btn-primary">
              Créer une soirée
            </Link>
          }
        />
      )}

      {state.status === 'ready' && state.soirees.length > 0 && (
        <ul className="soiree-grid">
          {state.soirees.map((soiree) => (
            <li key={soiree.id}>
              <Link to={`/soirees/${soiree.id}`} className="card soiree-card">
                <p className="card-title">{soiree.titre}</p>
                <p className="card-meta">
                  {soiree.lieu} · {new Date(soiree.date).toLocaleDateString('fr-FR')}
                </p>
                {soiree.description && <p>{soiree.description}</p>}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
