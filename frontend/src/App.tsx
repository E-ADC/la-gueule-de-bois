import { Navigate, Route, Routes } from 'react-router-dom'
import { AuthProvider } from './auth/AuthContext'
import { Layout } from './components/Layout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { ModeratorRoute } from './components/ModeratorRoute'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { FeedPage } from './pages/FeedPage'
import { SoireeFormPage } from './pages/SoireeFormPage'
import { SoireeDetailPage } from './pages/SoireeDetailPage'
import { ClassementPage } from './pages/ClassementPage'
import { ProfilePage } from './pages/ProfilePage'
import { ModerationPage } from './pages/ModerationPage'
import { NotFoundPage } from './pages/NotFoundPage'

/**
 * Écrans couverts par ce squelette : UC01/02/03 (auth), UC06/07/08/10 (soirées,
 * y compris upload photo), UC09/11/12/13 (témoignages, votes, invitation de
 * témoin, signalement), UC17/20 (classement), UC04 (lecture profil)/UC15/16
 * (badges + score), UC22 (modération, rôle Modérateur).
 *
 * TODO — cas d'utilisation pas encore couverts par une page (client API prêt
 * dans src/api/ quand ce sera le cas) :
 *   - UC04 (édition du profil)      → formulaire dans ProfilePage
 *   - UC05 (profil public d'un tiers) → route /utilisateurs/:id
 *   - UC18/19 (créer/rejoindre un groupe) → page Groupes (groupesApi prêt)
 *   - UC21 (demande d'ami)          → page/section Amis
 *   - UC12 (swipe) volontairement absent de ce lot (cf. consigne)
 */
export default function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/connexion" element={<LoginPage />} />
        <Route path="/inscription" element={<RegisterPage />} />

        <Route element={<Layout />}>
          <Route element={<ProtectedRoute />}>
            <Route path="/" element={<FeedPage />} />
            <Route path="/soirees/nouvelle" element={<SoireeFormPage />} />
            <Route path="/soirees/:id" element={<SoireeDetailPage />} />
            <Route path="/soirees/:id/modifier" element={<SoireeFormPage />} />
            <Route path="/classement" element={<ClassementPage />} />
            <Route path="/profil" element={<ProfilePage />} />
            <Route element={<ModeratorRoute />}>
              <Route path="/moderation" element={<ModerationPage />} />
            </Route>
          </Route>
          <Route path="/404" element={<NotFoundPage />} />
          <Route path="*" element={<Navigate to="/404" replace />} />
        </Route>
      </Routes>
    </AuthProvider>
  )
}
