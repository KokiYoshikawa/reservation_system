import { Link, Navigate, useLocation } from 'react-router-dom'
import type { ReservationCompleteRouteState } from '../lib/reservationFlow'

export function ReservationCompletePage() {
  const location = useLocation()
  const state = location.state as ReservationCompleteRouteState | null

  if (!state) {
    return <Navigate to="/reservations" replace />
  }

  return (
    <section className="page services-page">
      <span className="eyebrow">Reservation Complete</span>
      <h1>予約が完了しました</h1>
      <p className="lead">
        登録が完了しました。内容を確認して、必要に応じて予約一覧や空き枠検索へ進めます。
      </p>

      <div className="completion-card">
        <div className="completion-hero">
          <div>
            <p className="service-card-label">Reservation</p>
            <h2>{state.service.name}</h2>
          </div>
          <span className="service-card-badge">Success</span>
        </div>

        <dl className="reservation-summary-grid">
          <div>
            <dt>予約ID</dt>
            <dd>{state.reservationId}</dd>
          </div>
          <div>
            <dt>ステータス</dt>
            <dd>{state.status}</dd>
          </div>
          <div>
            <dt>開始時刻</dt>
            <dd>{state.startTime}</dd>
          </div>
          <div>
            <dt>終了時刻</dt>
            <dd>{state.endTime}</dd>
          </div>
          <div>
            <dt>予約日時</dt>
            <dd>{state.reservedAt}</dd>
          </div>
          <div>
            <dt>備考</dt>
            <dd>{state.note || 'なし'}</dd>
          </div>
        </dl>
      </div>

      <div className="actions">
        <Link className="button primary" to="/reservations">
          予約一覧へ
        </Link>
        <Link className="button secondary" to="/reservation-slots">
          もう一度空き枠を探す
        </Link>
      </div>
    </section>
  )
}
