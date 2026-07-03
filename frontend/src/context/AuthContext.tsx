import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import {
  AUTH_STORAGE_KEY,
  fetchCurrentUser,
  loginRequest,
  logoutRequest,
  type AuthUser,
} from '../lib/api'

type AuthState = {
  isAuthenticated: boolean
  isLoading: boolean
  token: string | null
  user: AuthUser | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
}

type StoredAuth = {
  token: string
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<AuthUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const storedValue = window.localStorage.getItem(AUTH_STORAGE_KEY)

    if (!storedValue) {
      setIsLoading(false)
      return
    }

    try {
      const parsed = JSON.parse(storedValue) as StoredAuth
      setToken(parsed.token)
    } catch {
      window.localStorage.removeItem(AUTH_STORAGE_KEY)
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!token) {
      setUser(null)
      return
    }

    let active = true

    async function syncSession() {
      try {
        const response = await fetchCurrentUser()
        if (!active) {
          return
        }

        setUser(response.user)
      } catch {
        if (!active) {
          return
        }

        window.localStorage.removeItem(AUTH_STORAGE_KEY)
        setToken(null)
        setUser(null)
      } finally {
        if (active) {
          setIsLoading(false)
        }
      }
    }

    void syncSession()

    return () => {
      active = false
    }
  }, [token])

  async function login(email: string, password: string) {
    if (!email || !password) {
      throw new Error('メールアドレスとパスワードを入力してください。')
    }

    const response = await loginRequest(email, password)
    const nextAuth: StoredAuth = {
      token: response.token,
    }

    window.localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(nextAuth))
    setToken(nextAuth.token)
    setUser(response.user)
  }

  function logout() {
    void logoutRequest()
    window.localStorage.removeItem(AUTH_STORAGE_KEY)
    setToken(null)
    setUser(null)
  }

  const value = useMemo<AuthState>(
    () => ({
      isAuthenticated: Boolean(token),
      isLoading,
      token,
      user,
      login,
      logout,
    }),
    [isLoading, token, user],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }

  return context
}
