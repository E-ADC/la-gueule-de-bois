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
import { GroupesPage } from './pages/GroupesPage'
import { AmisPage } from './pages/AmisPage'
import { ProfilePage } from './pages/ProfilePage'
import { PublicProfilePage } from './pages/PublicProfilePage'
import { ModerationPage } from './pages/ModerationPage'
import { NotFoundPage } from './pages/NotFoundPage'

/**
 * Écrans couverts : UC01/02/03 (auth), UC04/05 (profil, édition et public),
 * UC06/07/08/10 (soirées, upload photo), UC09/11/12/13 (témoignages, votes,
 * invitation de témoin, signalement), UC15/16 (badges + score), UC17/20
 * (classement), UC18/19 (groupes), UC21 (amis), UC22 (modération).
 *
 * UC12 (swipe) volontairement absent (cf. consigne).
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
            <Route path="/groupes" element={<GroupesPage />} />
            <Route path="/amis" element={<AmisPage />} />
            <Route path="/profil" element={<ProfilePage />} />
            <Route path="/utilisateurs/:id" element={<PublicProfilePage />} />
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
