import { useEffect, useState } from 'react'
import { classementApi } from '../api/classement'
import { groupesApi } from '../api/groupes'
import { ApiError } from '../api/client'
import type { Groupe, User } from '../api/types'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type Scope = 'global' | string // 'global' ou id de groupe

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; entries: User[] } // User[] triés par score décroissant

/** UC17 (classement global) et UC20 (classement de groupe, même vue). */
export function ClassementPage() {
  const [scope, setScope] = useState<Scope>('global')
  const [groupes, setGroupes] = useState<Groupe[]>([])
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)

  useEffect(() => {
    const controller = new AbortController()
    // GET /groupes n'existe pas encore côté backend (404/501) : on encaisse
    // silencieusement et le sélecteur ne propose que "Global".
    groupesApi
      .listMine(controller.signal)
      .then((groupes) => setGroupes(groupes ?? []))
      .catch(() => setGroupes([]))
    return () => controller.abort()
  }, [])

  useEffect(() => {
    const controller = new AbortController()
    setState({ status: 'loading' })
    const fetchEntries =
      scope === 'global'
        ? classementApi.global(controller.signal)
        : classementApi.groupe(scope, controller.signal)

    fetchEntries
      .then((entries) => setState({ status: 'ready', entries: entries ?? [] }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger le classement.',
        })
      })
    return () => controller.abort()
  }, [scope, reloadToken])

  return (
    <div className="page">
      <div className="page-header">
        <h1>Classement</h1>
        <select
          className="select"
          value={scope}
          onChange={(event) => setScope(event.target.value)}
          aria-label="Périmètre du classement"
        >
          <option value="global">Global</option>
          {groupes.map((groupe) => (
            <option key={groupe.id} value={groupe.id}>
              Groupe : {groupe.nom}
            </option>
          ))}
        </select>
      </div>

      {state.status === 'loading' && <Loading label="Chargement du classement…" />}

      {state.status === 'error' && (
        <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
      )}

      {state.status === 'ready' && state.entries.length === 0 && (
        <EmptyState
          title="Classement vide"
          message="Aucune donnée de classement disponible pour l’instant."
        />
      )}

      {state.status === 'ready' && state.entries.length > 0 && (
        <div className="card">
          {/* Pas de champ rang dans la réponse : déduit de l'index (liste déjà triée). */}
          {state.entries.map((entry, index) => (
            <div key={entry.id} className="rank-row">
              <span className="rank">{index + 1}</span>
              <span className="avatar">{entry.pseudo.slice(0, 2).toUpperCase()}</span>
              <span>{entry.pseudo}</span>
              <span className="pts">{entry.score} pts</span>
            </div>
          ))}
        </div>
      )}

      {/* TODO(UC18/UC19) : créer/rejoindre un groupe depuis cette page (pas de page dédiée encore) */}
    </div>
  )
}
