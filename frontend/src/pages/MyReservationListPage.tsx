import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import { fetchReservations, type ReservationItem } from '../lib/api'
import {
  formatReservationDate,
  formatReservationDateTime,
  formatReservationTime,
  reservationStatus,
} from '../lib/reservationPresentation'

type RequestState = 'loading' | 'success' | 'error'

export function MyReservationListPage() {
  const [state, setState] = useState<RequestState>('loading')
  const [reservations, setReservations] = useState<ReservationItem[]>([])
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    let active = true

    async function loadReservations() {
      setState('loading')
      setError('')

      try {
        const response = await fetchReservations()
        if (!active) {
          return
        }

        setReservations(response.data.reservations)
        setState('success')
      } catch (requestError) {
        if (!active) {
          return
        }

        setError(
          requestError instanceof Error ? requestError.message : 'Unknown error',
        )
        setState('error')
      }
    }

    void loadReservations()

    return () => {
      active = false
    }
  }, [reloadKey])

  return (
    <section className="page services-page reservations-page">
      <span className="eyebrow">My Reservations</span>
      <div className="reservation-page-heading">
        <div>
          <h1>マイ予約一覧</h1>
          <p className="lead">予約日時や現在の状態を確認できます。</p>
        </div>
        <Link className="button primary" to="/reservation-slots">
          新しく予約する
        </Link>
      </div>

      {state === 'loading' ? (
        <LoadingIndicator
          label="予約一覧を取得しています..."
          className="status-feedback"
        />
      ) : null}

      {state === 'error' ? (
        <div className="reservation-request-error">
          <ErrorMessage
            message={error}
            title="予約一覧の取得に失敗しました"
            className="status-feedback"
          />
          <button
            type="button"
            className="button secondary retry-button"
            onClick={() => setReloadKey((current) => current + 1)}
          >
            もう一度試す
          </button>
        </div>
      ) : null}

      {state === 'success' ? (
        reservations.length > 0 ? (
          <div className="reservation-list">
            {reservations.map((reservation) => {
              const status = reservationStatus(reservation.status)

              return (
                <article key={reservation.reservationId} className="reservation-card">
                  <div className="reservation-card-header">
                    <div>
                      <p className="service-card-label">
                        Reservation #{reservation.reservationId}
                      </p>
                      <h2>{reservation.service.name}</h2>
                    </div>
                    <span className={`reservation-status ${status.className}`}>
                      {status.label}
                    </span>
                  </div>

                  <div className="reservation-schedule">
                    <p className="reservation-date">
                      {formatReservationDate(reservation.startTime)}
                    </p>
                    <p className="reservation-time">
                      {formatReservationTime(reservation.startTime)}
                      <span aria-hidden="true"> – </span>
                      {formatReservationTime(reservation.endTime)}
                    </p>
                  </div>

                  <dl className="reservation-card-meta">
                    <div>
                      <dt>予約受付日時</dt>
                      <dd>{formatReservationDateTime(reservation.reservedAt)}</dd>
                    </div>
                    <div>
                      <dt>備考</dt>
                      <dd>{reservation.note || 'なし'}</dd>
                    </div>
                  </dl>

                  <div className="reservation-card-actions">
                    <Link
                      className="button secondary"
                      to={`/reservations/${reservation.reservationId}`}
                    >
                      詳細を見る
                    </Link>
                  </div>
                </article>
              )
            })}
          </div>
        ) : (
          <div className="reservation-empty-state">
            <p className="reservation-empty-icon" aria-hidden="true">○</p>
            <h2>予約はまだありません</h2>
            <p>空き枠を探して、最初の予約を登録しましょう。</p>
            <Link className="button primary" to="/reservation-slots">
              空き枠を探す
            </Link>
          </div>
        )
      ) : null}
    </section>
  )
}
