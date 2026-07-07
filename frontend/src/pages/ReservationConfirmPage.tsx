import { useState } from 'react'
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import { createReservationRequest } from '../lib/api'
import type {
  ReservationCompleteRouteState,
  ReservationConfirmRouteState,
} from '../lib/reservationFlow'

export function ReservationConfirmPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const state = location.state as ReservationConfirmRouteState | null
  const [note, setNote] = useState('')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  if (!state) {
    return <Navigate to="/reservation-slots" replace />
  }

  const reservationState = state

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setIsSubmitting(true)

    try {
      const response = await createReservationRequest({
        serviceId: reservationState.service.id,
        slotId: reservationState.slot.slotId,
        note,
      })

      const completeState: ReservationCompleteRouteState = {
        reservationId: response.data.reservationId,
        status: response.data.status,
        service: response.data.service,
        startTime: response.data.startTime,
        endTime: response.data.endTime,
        reservedAt: response.data.reservedAt,
        note,
      }

      navigate('/reservation-complete', {
        replace: true,
        state: completeState,
      })
    } catch (requestError) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : '予約登録に失敗しました。',
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <section className="page services-page">
      <span className="eyebrow">Reservation Confirm</span>
      <h1>予約内容の確認</h1>
      <p className="lead">
        選択したサービスと予約枠を確認して、必要であれば備考を入力してください。
      </p>

      <div className="confirmation-card">
        <div className="service-card-header">
          <div>
            <p className="service-card-label">Selected Service</p>
            <h2>{reservationState.service.name}</h2>
          </div>
          <span className="service-card-badge">Ready</span>
        </div>

        <dl className="reservation-summary-grid">
          <div>
            <dt>サービスID</dt>
            <dd>{reservationState.service.id}</dd>
          </div>
          <div>
            <dt>所要時間</dt>
            <dd>{reservationState.service.duration_minutes}分</dd>
          </div>
          <div>
            <dt>料金</dt>
            <dd>{reservationState.service.price.toLocaleString('ja-JP')}円</dd>
          </div>
          <div>
            <dt>日付</dt>
            <dd>{reservationState.date}</dd>
          </div>
          <div>
            <dt>開始時刻</dt>
            <dd>{reservationState.slot.startTime}</dd>
          </div>
          <div>
            <dt>終了時刻</dt>
            <dd>{reservationState.slot.endTime}</dd>
          </div>
        </dl>
      </div>

      <form className="confirmation-form" onSubmit={handleSubmit}>
        <label className="note-field">
          <span>備考</span>
          <textarea
            value={note}
            onChange={(event) => setNote(event.target.value)}
            placeholder="ご希望や補足事項があれば入力してください"
            rows={5}
          />
        </label>

        <div className="actions">
          <Link className="button secondary" to="/reservation-slots">
            空き枠選択に戻る
          </Link>
          <button type="submit" className="button primary confirm-submit-button" disabled={isSubmitting}>
            {isSubmitting ? '予約登録中...' : 'この内容で予約する'}
          </button>
        </div>
      </form>

      {isSubmitting ? (
        <LoadingIndicator
          label="予約情報を登録しています..."
          className="status-feedback"
        />
      ) : null}

      {error ? (
        <ErrorMessage
          title="予約登録に失敗しました"
          message={error}
          className="status-feedback"
        />
      ) : null}
    </section>
  )
}
