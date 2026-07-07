import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import { useAuth } from '../context/AuthContext'

export function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const { isAuthenticated, login, user } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const redirectTo = (location.state as { from?: string | { pathname: string; search?: string; hash?: string; state?: unknown } } | null)?.from ?? '/mypage'

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setIsSubmitting(true)

    try {
      await login(email, password)
      navigate(redirectTo, { replace: true })
    } catch (loginError) {
      setError(loginError instanceof Error ? loginError.message : 'ログインに失敗しました。')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <section className="auth-stage">
      <div className="auth-window">
        <div className="window-dots" aria-hidden="true">
          <span className="dot red"></span>
          <span className="dot yellow"></span>
          <span className="dot green"></span>
        </div>

        <div className="login-card">
          <p className="login-title">予約管理システム</p>
          <h1 className="login-heading">ログイン</h1>

          <form className="login-form" onSubmit={handleSubmit}>
            <label className="input-row">
              <span className="input-icon" aria-hidden="true">
                ✉
              </span>
              <input
                type="email"
                placeholder="メールアドレス"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
              />
            </label>

            <label className="input-row">
              <span className="input-icon" aria-hidden="true">
                🔒
              </span>
              <input
                type="password"
                placeholder="パスワード"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
              />
            </label>

            <button type="submit" className="login-button" disabled={isSubmitting}>
              {isSubmitting ? 'ログイン中' : 'ログイン'}
            </button>
          </form>

          {isSubmitting ? (
            <LoadingIndicator
              label="認証情報を確認しています..."
              className="login-feedback"
              compact
            />
          ) : null}

          {error ? (
            <ErrorMessage
              message={error}
              title="ログインに失敗しました"
              className="login-feedback"
            />
          ) : null}
          {isAuthenticated && user ? (
            <p className="login-message success">{user.email} でログイン中です。</p>
          ) : null}

          <div className="auth-links">
            <a className="forgot-link" href="#">
              パスワードを忘れた方はこちら
            </a>
            <Link className="forgot-link secondary-link" to="/register">
              新規ユーザー登録はこちら
            </Link>
          </div>
        </div>
      </div>
    </section>
  )
}
