import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <section className="page">
      <span className="eyebrow">404</span>
      <h1>ページが見つかりません。</h1>
      <p className="lead">
        URL を確認するか、トップページから目的の画面へ移動してください。
      </p>
      <div className="actions">
        <Link className="button primary" to="/">
          ホームへ戻る
        </Link>
      </div>
    </section>
  )
}
