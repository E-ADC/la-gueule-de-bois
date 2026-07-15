import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import { authApi, type LoginPayload, type RegisterPayload } from '../api/auth'
import { ApiError } from '../api/client'
import type { User } from '../api/types'

interface AuthContextValue {
  user: User | null
  /** Chargement initial de la session (`GET /auth/me`). */
  loading: boolean
  login: (payload: LoginPayload) => Promise<void>
  register: (payload: RegisterPayload) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const controller = new AbortController()
    authApi
      .me(controller.signal)
      .then(setUser)
      .catch((err: unknown) => {
        // 401 attendu si personne n'est connecté : pas une erreur à afficher.
        if (!(err instanceof ApiError) || err.status !== 401) {
          console.error("Impossible de vérifier la session en cours :", err)
        }
        setUser(null)
      })
      .finally(() => setLoading(false))
    return () => controller.abort()
  }, [])

  const login = useCallback(async (payload: LoginPayload) => {
    const loggedInUser = await authApi.login(payload)
    setUser(loggedInUser)
  }, [])

  const register = useCallback(async (payload: RegisterPayload) => {
    const createdUser = await authApi.register(payload)
    setUser(createdUser)
  }, [])

  const logout = useCallback(async () => {
    await authApi.logout()
    setUser(null)
  }, [])

  const value = useMemo(
    () => ({ user, loading, login, register, logout }),
    [user, loading, login, register, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth doit être utilisé à l’intérieur de <AuthProvider>.')
  }
  return ctx
}
