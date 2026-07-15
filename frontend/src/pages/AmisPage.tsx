import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { amisApi } from '../api/amis'
import { usersApi } from '../api/users'
import { ApiError } from '../api/client'
import type { DemandeAmi } from '../api/types'
import { Loading, ErrorState, EmptyState } from '../components/StateViews'

type LoadState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; demandes: DemandeAmi[] }

interface DemandeAffichee extends DemandeAmi {
  demandeurPseudo?: string
}

/**
 * UC21 : page Amis - envoyer une demande d\'ami et lister les demandes reçues.
 */
export function AmisPage() {
  const [state, setState] = useState<LoadState>({ status: 'loading' })
  const [reloadToken, setReloadToken] = useState(0)
  const [demandeur, setDemandeur] = useState('')
  const [sendError, setSendError] = useState<string | null>(null)
  const [sending, setSending] = useState(false)
  const [demandesAffichees, setDemandesAffichees] = useState<DemandeAffichee[]>([])
  const [resolvingPseudos, setResolvingPseudos] = useState(false)

  // Chargement des demandes reçues
  useEffect(() => {
    const controller = new AbortController()
    setState({ status: 'loading' })
    amisApi
      .listRecues(controller.signal)
      .then((demandes) => setState({ status: 'ready', demandes: demandes ?? [] }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({
          status: 'error',
          message:
            err instanceof ApiError ? err.message : 'Impossible de charger les demandes d\'ami.',
        })
      })
    return () => controller.abort()
  }, [reloadToken])

  // Résolution des pseudos des demandeurs
  useEffect(() => {
    if (state.status !== 'ready') {
      setDemandesAffichees([])
      return
    }

    setResolvingPseudos(true)
    const demandes = state.demandes

    Promise.all(
      demandes.map(async (demande) => {
        try {
          const profile = await usersApi.publicProfile(demande.demandeurId)
          return { ...demande, demandeurPseudo: profile.user.pseudo }
        } catch {
          // En cas d\'erreur de résolution du pseudo, on garde le demande sans pseudo
          return demande
        }
      }),
    )
      .then((resolved) => {
        setDemandesAffichees(resolved)
        setResolvingPseudos(false)
      })
      .catch(() => {
        // Fallback : affiche les demandes même si la résolution échoue
        setDemandesAffichees(demandes)
        setResolvingPseudos(false)
      })
  }, [state])

  async function handleSendRequest(event: FormEvent) {
    event.preventDefault()
    if (!demandeur.trim()) return
    setSendError(null)
    setSending(true)
    try {
      await amisApi.envoyer(demandeur.trim())
      setDemandeur('')
      setReloadToken((token) => token + 1)
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setSendError('Aucun compte avec ce pseudo.')
      } else if (err instanceof ApiError && err.status === 409) {
        setSendError('Demande déjà envoyée ou vous êtes déjà amis.')
      } else {
        setSendError(err instanceof ApiError ? err.message : 'Envoi de la demande échoué.')
      }
    } finally {
      setSending(false)
    }
  }

  async function handleRespond(demandeId: number, action: 'accepter' | 'refuser') {
    try {
      await amisApi.repondre(demandeId, action)
      setReloadToken((token) => token + 1)
    } catch (err) {
      window.alert(err instanceof ApiError ? err.message : 'Action échouée.')
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <h1>Amis</h1>
      </div>

      <h2>Envoyer une demande d\'ami</h2>
      <form className="card" onSubmit={(event) => void handleSendRequest(event)}>
        <label className="label" htmlFor="demandeur">
          Pseudo de la personne (UC21)
        </label>
        <input
          id="demandeur"
          type="text"
          className="input"
          value={demandeur}
          onChange={(event) => setDemandeur(event.target.value)}
          placeholder="pseudo"
          required
        />
        {sendError && <p className="field-error">{sendError}</p>}
        <button type="submit" className="btn btn-primary" disabled={sending}>
          {sending ? 'Envoi…' : 'Envoyer'}
        </button>
      </form>

      <hr className="sep" />

      <h2>Demandes reçues</h2>

      {state.status === 'loading' && <Loading label="Chargement des demandes d\'ami…" />}

      {state.status === 'error' && (
        <ErrorState message={state.message} onRetry={() => setReloadToken((t) => t + 1)} />
      )}

      {state.status === 'ready' && demandesAffichees.length === 0 && (
        <EmptyState
          title="Aucune demande d\'ami"
          message="Aucune demande d\'ami en attente."
        />
      )}

      {state.status === 'ready' && demandesAffichees.length > 0 && (
        <ul className="friend-request-list">
          {demandesAffichees.map((demande) => (
            <li key={demande.id} className="card friend-request-card">
              <div className="friend-request-header">
                <p className="card-title">
                  {demande.demandeurPseudo ?? `Utilisateur #${demande.demandeurId}`}
                </p>
                {!resolvingPseudos && (
                  <p className="card-meta">
                    {new Date(demande.createdAt).toLocaleDateString('fr-FR')}
                  </p>
                )}
              </div>
              <div className="vote-row">
                <button
                  type="button"
                  className="btn btn-primary"
                  onClick={() => void handleRespond(demande.id, 'accepter')}
                >
                  Accepter
                </button>
                <button
                  type="button"
                  className="btn btn-danger"
                  onClick={() => void handleRespond(demande.id, 'refuser')}
                >
                  Refuser
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
