import { useEffect, useState } from 'react'
import { signalementsApi } from '../api/signalements'
import { ApiError } from '../api/client'
import type { SignalementView } from '../api/types'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; signalements: SignalementView[] }

/** UC22 : le modérateur consulte les signalements en attente et décide d'une action. */
export function ModerationPage() {
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)
  const [actioningId, setActioningId] = useState<number | null>(null)
  const [actionMessage, setActionMessage] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()
    setState({ status: 'loading' })
    signalementsApi
      .listEnAttente(controller.signal)
      .then((signalements) => setState({ status: 'ready', signalements: signalements ?? [] }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger les signalements.',
        })
      })
    return () => controller.abort()
  }, [reloadToken])

  async function handleTraiter(id: number, action: 'rejeter' | 'supprimer') {
    setActionMessage(null)
    setActioningId(id)
    try {
      await signalementsApi.traiter(id, action)
      setReloadToken((token) => token + 1)
    } catch (err) {
      setActionMessage(err instanceof ApiError ? err.message : 'Le traitement a échoué.')
    } finally {
      setActioningId(null)
    }
  }

  if (state.status === 'loading') {
    return <Loading label="Chargement des signalements…" />
  }

  if (state.status === 'error') {
    return <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
  }

  const { signalements } = state

  return (
    <div className="page">
      <div className="page-header">
        <h1>Modération</h1>
      </div>

      {actionMessage && <p className="field-error">{actionMessage}</p>}

      {signalements.length === 0 ? (
        <EmptyState
          title="Aucun signalement en attente"
          message="Tous les témoignages signalés ont été traités."
        />
      ) : (
        <ul className="temoignage-list">
          {signalements.map((signalement) => (
            <li key={signalement.id} className="card temoignage-card">
              <p className="card-meta">
                Signalé le {new Date(signalement.createdAt).toLocaleDateString('fr-FR')}
              </p>
              <p className="card-title">Motif : {signalement.motif}</p>
              <p>{signalement.temoignageContenu}</p>
              <div className="vote-row">
                <button
                  type="button"
                  className="btn btn-ghost"
                  disabled={actioningId === signalement.id}
                  onClick={() => void handleTraiter(signalement.id, 'rejeter')}
                >
                  Rejeter le signalement
                </button>
                <button
                  type="button"
                  className="btn btn-danger"
                  disabled={actioningId === signalement.id}
                  onClick={() => void handleTraiter(signalement.id, 'supprimer')}
                >
                  Supprimer le témoignage
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
