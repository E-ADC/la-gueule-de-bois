import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { Loading } from './StateViews'

/** Garde de route : redirige vers /connexion si personne n'est authentifié. */
export function ProtectedRoute() {
  const { user, loading } = useAuth()
  const location = useLocation()

  if (loading) {
    return <Loading label="Vérification de la session…" />
  }

  if (!user) {
    return <Navigate to="/connexion" replace state={{ from: location }} />
  }

  return <Outlet />
}
