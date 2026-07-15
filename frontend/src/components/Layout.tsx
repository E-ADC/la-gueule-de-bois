import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'

function navClass({ isActive }: { isActive: boolean }) {
  return isActive ? 'nav-link nav-link-active' : 'nav-link'
}

export function Layout() {
  const { user, logout } = useAuth()

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
