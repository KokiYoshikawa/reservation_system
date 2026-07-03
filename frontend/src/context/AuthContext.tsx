import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'

type AuthUser = {
  email: string
}

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
  user: AuthUser
}

const STORAGE_KEY = 'reservation_system_auth'

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<AuthUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const storedValue = window.localStorage.getItem(STORAGE_KEY)

    if (!storedValue) {
      setIsLoading(false)
      return
    }

    try {
      const parsed = JSON.parse(storedValue) as StoredAuth
      setToken(parsed.token)
      setUser(parsed.user)
    } catch {
      window.localStorage.removeItem(STORAGE_KEY)
    } finally {
      setIsLoading(false)
    }
  }, [])

  async function login(email: string, password: string) {
    if (!email || !password) {
      throw new Error('メールアドレスとパスワードを入力してください。')
    }

    const nextAuth: StoredAuth = {
      token: `demo-token-${Date.now()}`,
      user: { email },
    }

    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(nextAuth))
    setToken(nextAuth.token)
    setUser(nextAuth.user)
  }

  function logout() {
    window.localStorage.removeItem(STORAGE_KEY)
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
