import { api } from './client'
import type { User } from './types'

export interface RegisterPayload {
  pseudo: string
  email: string
  password: string
}

export interface LoginPayload {
  email: string
  password: string
}

// Endpoints supposés (session cookie HttpOnly, cf. spec) :
export const authApi = {
  register: (payload: RegisterPayload) => api.post<User>('/auth/register', payload),
  login: (payload: LoginPayload) => api.post<User>('/auth/login', payload),
  logout: () => api.post<void>('/auth/logout'),
  me: (signal?: AbortSignal) => api.get<User>('/auth/me', signal),
}
