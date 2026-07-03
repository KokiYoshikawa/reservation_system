import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function ReservationsPage() {
  const { isAuthenticated, user } = useAuth()

  return (
    <section className="page">
      <span className="eyebrow">Reservations</span>
      <h1>予約一覧ページ</h1>
      {isAuthenticated && user ? (
        <p className="lead">
          {user.email} としてログイン中です。この画面に今後、予約検索・一覧・詳細遷移を追加していけます。
        </p>
      ) : (
        <>
          <p className="lead">
            この画面はログイン後の利用を想定しています。認証状態管理の確認用として、トップからログインしてください。
          </p>
          <div className="actions">
            <Link className="button primary" to="/">
              ログイン画面へ戻る
            </Link>
          </div>
        </>
      )}
    </section>
  )
}
