import { useEffect, useState } from 'react'
import { groupesApi } from '../api/groupes'
import { ApiError } from '../api/client'
import type { Groupe } from '../api/types'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; groupes: Groupe[] }

export function GroupesPage() {
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)

  // Formulaire créer groupe
  const [createForm, setCreateForm] = useState({ nom: '', error: '' })

  // Formulaire rejoindre groupe
  const [joinForm, setJoinForm] = useState({ id: '', error: '' })

  // Chargement de la liste des groupes
  useEffect(() => {
    const controller = new AbortController()
    setState({ status: 'loading' })
    groupesApi
      .listMine(controller.signal)
      .then((groupes) => setState({ status: 'ready', groupes: groupes ?? [] }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message: err instanceof ApiError ? err.message : 'Impossible de charger les groupes.',
        })
      })
    return () => controller.abort()
  }, [reloadToken])

  const handleCreateGroup = async (e: React.FormEvent) => {
    e.preventDefault()
    setCreateForm({ ...createForm, error: '' })

    if (!createForm.nom.trim()) {
      setCreateForm({ ...createForm, error: 'Le nom du groupe ne peut pas être vide.' })
      return
    }

    try {
      await groupesApi.create(createForm.nom)
      setCreateForm({ nom: '', error: '' })
      setReloadToken((t) => t + 1)
    } catch (err: unknown) {
      if (err instanceof ApiError) {
        if (err.status === 409) {
          setCreateForm({ ...createForm, error: 'Ce nom de groupe est déjà utilisé.' })
        } else {
          setCreateForm({ ...createForm, error: err.message })
        }
      } else {
        setCreateForm({ ...createForm, error: 'Une erreur est survenue lors de la création.' })
      }
    }
  }

  const handleJoinGroup = async (e: React.FormEvent) => {
    e.preventDefault()
    setJoinForm({ ...joinForm, error: '' })

    if (!joinForm.id.trim()) {
      setJoinForm({ ...joinForm, error: "L'ID du groupe ne peut pas être vide." })
      return
    }

    const groupeId = parseInt(joinForm.id, 10)
    if (isNaN(groupeId)) {
      setJoinForm({ ...joinForm, error: "L'ID du groupe doit être un nombre." })
      return
    }

    try {
      await groupesApi.join(groupeId)
      setJoinForm({ id: '', error: '' })
      setReloadToken((t) => t + 1)
    } catch (err: unknown) {
      if (err instanceof ApiError) {
        setJoinForm({ ...joinForm, error: err.message })
      } else {
        setJoinForm({ ...joinForm, error: 'Une erreur est survenue lors de la tentative de rejoindre.' })
      }
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <h1>Groupes</h1>
      </div>

      {state.status === 'loading' && <Loading label="Chargement des groupes…" />}

      {state.status === 'error' && (
        <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
      )}

      {state.status === 'ready' && (
        <>
          {/* Formulaire créer groupe */}
          <div className="card">
            <h2 className="card-title">Créer un groupe</h2>
            <form onSubmit={handleCreateGroup}>
              <div>
                <label htmlFor="create-nom" className="label">
                  Nom du groupe
                </label>
                <input
                  id="create-nom"
                  type="text"
                  className="input"
                  placeholder="Ex: Les copains de la fac"
                  value={createForm.nom}
                  onChange={(e) => setCreateForm({ ...createForm, nom: e.target.value })}
                />
                {createForm.error && <p className="field-error">{createForm.error}</p>}
              </div>
              <button type="submit" className="btn btn-primary">
                Créer
              </button>
            </form>
          </div>

          {/* Formulaire rejoindre groupe */}
          <div className="card">
            <h2 className="card-title">Rejoindre un groupe</h2>
            <form onSubmit={handleJoinGroup}>
              <div>
                <label htmlFor="join-id" className="label">
                  ID du groupe
                </label>
                <input
                  id="join-id"
                  type="number"
                  className="input"
                  placeholder="Ex: 42"
                  value={joinForm.id}
                  onChange={(e) => setJoinForm({ ...joinForm, id: e.target.value })}
                />
                {joinForm.error && <p className="field-error">{joinForm.error}</p>}
              </div>
              <button type="submit" className="btn btn-primary">
                Rejoindre
              </button>
            </form>
          </div>

          {/* Liste des groupes */}
          {state.groupes.length === 0 && (
            <EmptyState
              title="Aucun groupe"
              message="Tu n'es membre d'aucun groupe pour l'instant."
            />
          )}

          {state.groupes.length > 0 && (
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {state.groupes.map((groupe) => (
                <li key={groupe.id}>
                  <div className="card">
                    <p className="card-title">{groupe.nom}</p>
                    <p className="card-meta">
                      Créé le {new Date(groupe.createdAt).toLocaleDateString('fr-FR')}
                    </p>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </>
      )}
    </div>
  )
}
