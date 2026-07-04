import { useState } from 'react'
import { Link } from 'react-router-dom'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import { createUserRequest } from '../lib/api'

export function RegisterPage() {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState('user')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setSuccess('')
    setIsSubmitting(true)

    try {
      await createUserRequest({
        name,
        email,
        password_hash: password,
        role,
      })
      setSuccess('ユーザーを登録しました。ログイン画面から認証できます。')
      setName('')
      setEmail('')
      setPassword('')
      setRole('user')
    } catch (requestError) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : 'ユーザー登録に失敗しました。',
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <section className="auth-stage">
      <div className="auth-window register-window">
        <div className="window-dots" aria-hidden="true">
          <span className="dot red"></span>
          <span className="dot yellow"></span>
          <span className="dot green"></span>
        </div>

        <div className="login-card">
          <p className="login-title">予約管理システム</p>
          <h1 className="login-heading">ユーザー登録</h1>

          <form className="login-form" onSubmit={handleSubmit}>
            <label className="input-row">
              <span className="input-icon" aria-hidden="true">
                👤
              </span>
              <input
                type="text"
                placeholder="氏名"
                value={name}
                onChange={(event) => setName(event.target.value)}
              />
            </label>

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

            <label className="select-row">
              <span className="input-icon" aria-hidden="true">
                ⚙
              </span>
              <select value={role} onChange={(event) => setRole(event.target.value)}>
                <option value="user">user</option>
                <option value="admin">admin</option>
              </select>
            </label>

            <button type="submit" className="login-button" disabled={isSubmitting}>
              {isSubmitting ? '登録中' : '登録する'}
            </button>
          </form>

          {isSubmitting ? (
            <LoadingIndicator
              label="ユーザー情報を登録しています..."
              className="login-feedback"
              compact
            />
          ) : null}

          {error ? (
            <ErrorMessage
              message={error}
              title="ユーザー登録に失敗しました"
              className="login-feedback"
            />
          ) : null}

          {success ? <p className="login-message success">{success}</p> : null}

          <div className="auth-links">
            <Link className="forgot-link" to="/">
              ログイン画面へ戻る
            </Link>
          </div>
        </div>
      </div>
    </section>
  )
}
