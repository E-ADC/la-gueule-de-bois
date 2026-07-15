import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { Loading } from './StateViews'

/** Garde de route : accès réservé au rôle "moderator" (UC22). */
export function ModeratorRoute() {
  const { user, loading } = useAuth()

  if (loading) {
    return <Loading label="Vérification de la session…" />
  }

  if (!user || user.role !== 'moderator') {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}
