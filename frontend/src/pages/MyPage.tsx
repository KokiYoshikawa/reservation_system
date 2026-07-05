import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import { useAuth } from '../context/AuthContext'
import { fetchCurrentUser, type AuthUser } from '../lib/api'

export function MyPage() {
  const { isAuthenticated, isLoading, user: sessionUser } = useAuth()
  const [profile, setProfile] = useState<AuthUser | null>(sessionUser)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    setProfile(sessionUser)
  }, [sessionUser])

  useEffect(() => {
    if (isLoading || !isAuthenticated) {
      return
    }

    let active = true

    async function loadProfile() {
      setIsRefreshing(true)
      setError('')

      try {
        const response = await fetchCurrentUser()
        if (!active) {
          return
        }

        setProfile(response.user)
      } catch (requestError) {
        if (!active) {
          return
        }

        setError(
          requestError instanceof Error
            ? requestError.message
            : 'ユーザー情報の取得に失敗しました。',
        )
      } finally {
        if (active) {
          setIsRefreshing(false)
        }
      }
    }

    void loadProfile()

    return () => {
      active = false
    }
  }, [isAuthenticated, isLoading])

  if (isLoading) {
    return (
      <section className="page">
        <span className="eyebrow">My Page</span>
        <h1>マイページ</h1>
        <LoadingIndicator label="認証状態を確認しています..." className="status-feedback" />
      </section>
    )
  }

  if (!isAuthenticated) {
    return (
      <section className="page">
        <span className="eyebrow">My Page</span>
        <h1>マイページ</h1>
        <p className="lead">
          ログイン画面に戻り、ログインしてください。
        </p>
        <div className="actions">
          <Link className="button primary" to="/login">
            ログイン画面へ
          </Link>
        </div>
      </section>
    )
  }

  return (
    <section className="page">
      <span className="eyebrow">My Page</span>
      <h1>ログインユーザー情報</h1>
      <p className="lead">現在ログインしているユーザー情報を表示しています。</p>

      <div className="profile-card">
        <div className="profile-header">
          <span className="profile-badge">Authenticated</span>
        </div>

        <dl className="profile-grid">
          <div className="profile-row">
            <dt>ユーザーID</dt>
            <dd>{profile?.id ?? '取得できませんでした'}</dd>
          </div>
          <div className="profile-row">
            <dt>氏名</dt>
            <dd>{profile?.name ?? '取得できませんでした'}</dd>
          </div>
          <div className="profile-row">
            <dt>メールアドレス</dt>
            <dd>{profile?.email ?? '取得できませんでした'}</dd>
          </div>
          <div className="profile-row">
            <dt>ロール</dt>
            <dd>{profile?.role ?? '取得できませんでした'}</dd>
          </div>
        </dl>
      </div>

      {isRefreshing ? (
        <LoadingIndicator
          label="ユーザー情報を取得しています..."
          className="status-feedback"
          compact
        />
      ) : null}

      {error ? (
        <ErrorMessage
          title="マイページの読み込みに失敗しました"
          message={error}
          className="status-feedback"
        />
      ) : null}

      <div className="actions">
        <Link className="button secondary" to="/reservations">
          予約一覧へ
        </Link>
      </div>
    </section>
  )
}
