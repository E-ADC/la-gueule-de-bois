import type { ReactNode } from 'react'

/** Petits composants d'état partagés (chargement / erreur / vide) pour uniformiser les pages. */

export function Loading({ label = 'Chargement…' }: { label?: string }) {
  return (
    <div className="state-view">
      <p className="label">{label}</p>
    </div>
  )
}

export function ErrorState({
  message,
  onRetry,
}: {
  message: string
  onRetry?: () => void
}) {
  return (
    <div className="state-view state-view-error card">
      <p className="card-title">Oups</p>
      <p>{message}</p>
      {onRetry && (
        <button type="button" className="btn btn-ghost" onClick={onRetry}>
          Réessayer
        </button>
      )}
    </div>
  )
}

export function EmptyState({
  title,
  message,
  action,
}: {
  title: string
  message: string
  action?: ReactNode
}) {
  return (
    <div className="state-view state-view-empty">
      <p className="card-title">{title}</p>
      <p>{message}</p>
      {action}
    </div>
  )
}
