import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <div className="page">
      <div className="card">
        <p className="card-title">Page introuvable</p>
        <p>Cette page n’existe pas (encore ?).</p>
        <Link to="/" className="btn btn-primary">
          Retour aux soirées
        </Link>
      </div>
    </div>
  )
}
