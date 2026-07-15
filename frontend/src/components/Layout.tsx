import { useEffect, useState } from 'react'
import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { amisApi } from '../api/amis'
import { soireesApi } from '../api/soirees'

function navClass({ isActive }: { isActive: boolean }) {
  return isActive ? 'nav-link nav-link-active' : 'nav-link'
}

/** Icône de notifications : demandes d'ami reçues + invitations de témoin
 * en attente. Pas de suivi "lu/non lu" — juste un compteur en direct. */
function BellIcon({ size = 20 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" />
      <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" />
    </svg>
  )
}

export function Layout() {
  const { user, logout } = useAuth()
  const [notifCount, setNotifCount] = useState(0)

  useEffect(() => {
    if (!user) return
    const controller = new AbortController()
    Promise.all([
      amisApi.listRecues(controller.signal).catch(() => []),
      soireesApi.listInvitations(controller.signal).catch(() => []),
    ]).then(([demandes, invitations]) => {
      if (controller.signal.aborted) return
      setNotifCount((demandes?.length ?? 0) + (invitations?.length ?? 0))
    })
    return () => controller.abort()
  }, [user])

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="app-header-inner">
          <NavLink to="/" className="brand">
            La Gueule de Bois
          </NavLink>
          {user && (
            <nav className="main-nav">
              <NavLink to="/" className={navClass} end>
                Soirées
              </NavLink>
              <NavLink to="/classement" className={navClass}>
                Classement
              </NavLink>
              <NavLink to="/groupes" className={navClass}>
                Groupes
              </NavLink>
              <NavLink to="/amis" className={navClass}>
                Amis
              </NavLink>
              <NavLink to="/profil" className={navClass}>
                Profil
              </NavLink>
              {user.role === 'moderator' && (
                <NavLink to="/moderation" className={navClass}>
                  Modération
                </NavLink>
              )}
            </nav>
          )}
          <div className="header-actions">
            {user ? (
              <>
                <NavLink
                  to="/amis"
                  className="bell-link"
                  aria-label={`Notifications (${notifCount})`}
                  title="Demandes d'ami et invitations en attente"
                >
                  <BellIcon />
                  {notifCount > 0 && <span className="bell-badge">{notifCount}</span>}
                </NavLink>
                <span className="label">{user.pseudo}</span>
                <button type="button" className="btn btn-ghost" onClick={() => void logout()}>
                  Déconnexion
                </button>
              </>
            ) : (
              <>
                <NavLink to="/connexion" className="btn btn-ghost">
                  Connexion
                </NavLink>
                <NavLink to="/inscription" className="btn btn-primary">
                  Inscription
                </NavLink>
              </>
            )}
          </div>
        </div>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  )
}
