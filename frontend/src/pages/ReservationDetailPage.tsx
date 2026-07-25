import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import { fetchReservationDetail, type ReservationItem } from '../lib/api'
import {
  formatReservationDateTime,
  reservationStatus,
} from '../lib/reservationPresentation'

type RequestState = 'loading' | 'success' | 'error'

export function ReservationDetailPage() {
  const { reservationId: reservationIdParam } = useParams()
  const reservationId = Number(reservationIdParam)
  const validReservationId = Number.isInteger(reservationId) && reservationId > 0
  const [state, setState] = useState<RequestState>(
    validReservationId ? 'loading' : 'error',
  )
  const [reservation, setReservation] = useState<ReservationItem | null>(null)
  const [error, setError] = useState(
    validReservationId ? '' : '予約IDが正しくありません。',
  )
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    if (!validReservationId) {
      return
    }

    let active = true

    async function loadReservation() {
      setState('loading')
      setError('')

      try {
        const response = await fetchReservationDetail(reservationId)
        if (!active) {
          return
        }

        setReservation(response.data)
        setState('success')
      } catch (requestError) {
        if (!active) {
          return
        }

        setReservation(null)
        setError(
          requestError instanceof Error ? requestError.message : 'Unknown error',
        )
        setState('error')
      }
    }

    void loadReservation()

    return () => {
      active = false
    }
  }, [reloadKey, reservationId, validReservationId])

  const status = reservation ? reservationStatus(reservation.status) : null
  const canManage = reservation?.status === 'reserved'

  return (
    <section className="page services-page reservations-page">
      <Link className="reservation-back-link" to="/reservations">
        ← マイ予約一覧へ戻る
      </Link>
      <span className="eyebrow">Reservation Detail</span>
      <h1>予約詳細</h1>

      {state === 'loading' ? (
        <LoadingIndicator
          label="予約詳細を取得しています..."
          className="status-feedback"
        />
      ) : null}

      {state === 'error' ? (
        <div className="reservation-request-error">
          <ErrorMessage
            message={error}
            title="予約詳細の取得に失敗しました"
            className="status-feedback"
          />
          {validReservationId ? (
            <button
              type="button"
              className="button secondary retry-button"
              onClick={() => setReloadKey((current) => current + 1)}
            >
              もう一度試す
            </button>
          ) : null}
        </div>
      ) : null}

      {state === 'success' && reservation && status ? (
        <article className="reservation-detail-card">
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

          <dl className="reservation-detail-grid">
            <div className="reservation-detail-wide">
              <dt>予約日時</dt>
              <dd>
                {formatReservationDateTime(reservation.startTime)}
                <span className="reservation-date-separator">〜</span>
                {formatReservationDateTime(reservation.endTime)}
              </dd>
            </div>
            <div>
              <dt>サービス</dt>
              <dd>{reservation.service.name}</dd>
            </div>
            <div>
              <dt>予約枠ID</dt>
              <dd>{reservation.slotId}</dd>
            </div>
            <div>
              <dt>予約受付日時</dt>
              <dd>{formatReservationDateTime(reservation.reservedAt)}</dd>
            </div>
            <div>
              <dt>最終更新日時</dt>
              <dd>{formatReservationDateTime(reservation.updatedAt)}</dd>
            </div>
            <div className="reservation-detail-wide">
              <dt>備考</dt>
              <dd className="reservation-note">{reservation.note || 'なし'}</dd>
            </div>
            {reservation.cancelledAt ? (
              <div className="reservation-detail-wide">
                <dt>キャンセル日時</dt>
                <dd>{formatReservationDateTime(reservation.cancelledAt)}</dd>
              </div>
            ) : null}
          </dl>

          <div className="reservation-management">
            <div>
              <h2>予約の変更・キャンセル</h2>
              <p>
                {canManage
                  ? '変更・キャンセルの受付機能は準備中です。'
                  : 'この予約は変更・キャンセルできません。'}
              </p>
            </div>
            <div className="reservation-management-actions">
              <button type="button" className="button secondary" disabled>
                予約を変更する
              </button>
              <button type="button" className="button danger" disabled>
                予約をキャンセル
              </button>
            </div>
          </div>
        </article>
      ) : null}
    </section>
  )
}
