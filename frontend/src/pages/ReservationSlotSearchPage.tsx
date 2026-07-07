import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import {
  fetchReservationSlots,
  fetchServices,
  type ReservationSlotItem,
  type ServiceItem,
} from '../lib/api'
import type { ReservationConfirmRouteState } from '../lib/reservationFlow'

type RequestState = 'idle' | 'loading' | 'success' | 'error'

function todayDateString() {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function ReservationSlotSearchPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [services, setServices] = useState<ServiceItem[]>([])
  const [servicesState, setServicesState] = useState<RequestState>('idle')
  const [servicesError, setServicesError] = useState('')
  const [selectedServiceId, setSelectedServiceId] = useState('')
  const [selectedDate, setSelectedDate] = useState(todayDateString)

  const [slots, setSlots] = useState<ReservationSlotItem[]>([])
  const [slotsState, setSlotsState] = useState<RequestState>('idle')
  const [slotsError, setSlotsError] = useState('')
  const [searchedDate, setSearchedDate] = useState('')
  const [selectedSlotId, setSelectedSlotId] = useState<number | null>(null)
  const initialServiceId = searchParams.get('serviceId') ?? ''
  const selectedService = services.find((service) => String(service.id) === selectedServiceId) ?? null
  const selectedSlot = slots.find((slot) => slot.slotId === selectedSlotId) ?? null

  useEffect(() => {
    let active = true

    async function loadServices() {
      setServicesState('loading')
      setServicesError('')

      try {
        const response = await fetchServices()
        if (!active) {
          return
        }

        setServices(response.services)
        setSelectedServiceId((current) =>
          current ||
          (response.services.some((service) => String(service.id) === initialServiceId)
            ? initialServiceId
            : response.services[0]
              ? String(response.services[0].id)
              : ''),
        )
        setServicesState('success')
      } catch (requestError) {
        if (!active) {
          return
        }

        setServicesError(
          requestError instanceof Error ? requestError.message : 'Unknown error',
        )
        setServicesState('error')
      }
    }

    void loadServices()

    return () => {
      active = false
    }
  }, [initialServiceId])

  async function handleSearch(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!selectedServiceId || !selectedDate) {
      setSlotsError('サービスと日付を選択してください。')
      setSlotsState('error')
      return
    }

    setSlotsState('loading')
    setSlotsError('')
    setSelectedSlotId(null)

    try {
      const response = await fetchReservationSlots(selectedDate, Number(selectedServiceId))
      setSlots(response.data.slots)
      setSearchedDate(response.data.date)
      setSlotsState('success')
    } catch (requestError) {
      setSlots([])
      setSearchedDate(selectedDate)
      setSlotsError(
        requestError instanceof Error ? requestError.message : 'Unknown error',
      )
      setSlotsState('error')
    }
  }

  function handleProceedToConfirm() {
    if (!selectedService || !selectedSlot || !searchedDate) {
      return
    }

    const routeState: ReservationConfirmRouteState = {
      service: {
        id: selectedService.id,
        name: selectedService.name,
        duration_minutes: selectedService.duration_minutes,
        price: selectedService.price,
      },
      slot: selectedSlot,
      date: searchedDate,
    }

    navigate('/reservation-confirm', { state: routeState })
  }

  return (
    <section className="page services-page">
      <span className="eyebrow">Reservation Slots</span>
      <h1>空き枠検索</h1>
      <p className="lead">
        サービスと日付を選択して、予約可能な時間帯を確認できます。
      </p>

      <form className="slot-search-form" onSubmit={handleSearch}>
        <label className="slot-search-field">
          <span>サービス</span>
          <select
            value={selectedServiceId}
            onChange={(event) => setSelectedServiceId(event.target.value)}
            disabled={servicesState === 'loading' || services.length === 0}
          >
            {services.length === 0 ? <option value="">サービスを選択</option> : null}
            {services.map((service) => (
              <option key={service.id} value={service.id}>
                {service.name}
              </option>
            ))}
          </select>
        </label>

        <label className="slot-search-field">
          <span>日付</span>
          <input
            type="date"
            value={selectedDate}
            onChange={(event) => setSelectedDate(event.target.value)}
          />
        </label>

        <button
          type="submit"
          className="button primary slot-search-button"
          disabled={servicesState === 'loading'}
        >
          空き枠を検索
        </button>
      </form>

      {servicesState === 'loading' ? (
        <LoadingIndicator label="サービス一覧を読み込んでいます..." className="status-feedback" />
      ) : null}

      {servicesState === 'error' ? (
        <ErrorMessage
          message={servicesError}
          title="サービス一覧の取得に失敗しました"
          className="status-feedback"
        />
      ) : null}

      {slotsState === 'loading' ? (
        <LoadingIndicator label="空き枠を検索しています..." className="status-feedback" />
      ) : null}

      {slotsState === 'error' ? (
        <ErrorMessage
          message={slotsError}
          title="空き枠検索に失敗しました"
          className="status-feedback"
        />
      ) : null}

      {slotsState === 'success' ? (
        slots.length > 0 ? (
          <>
            <div className="service-list">
              {slots.map((slot) => (
                <article
                  key={slot.slotId}
                  className={`service-card ${selectedSlotId === slot.slotId ? 'selected' : ''}`}
                >
                  <div className="service-card-header">
                    <div>
                      <p className="service-card-label">Slot</p>
                      <h2>{searchedDate}</h2>
                    </div>
                    <span className={`service-card-badge ${slot.available ? '' : 'inactive'}`}>
                      {slot.available ? 'Available' : 'Full'}
                    </span>
                  </div>
                  <dl className="service-meta">
                    <div>
                      <dt>開始時刻</dt>
                      <dd>{slot.startTime}</dd>
                    </div>
                    <div>
                      <dt>終了時刻</dt>
                      <dd>{slot.endTime}</dd>
                    </div>
                    <div>
                      <dt>定員</dt>
                      <dd>{slot.capacity}名</dd>
                    </div>
                    <div>
                      <dt>予約済み</dt>
                      <dd>{slot.reservedCount}名</dd>
                    </div>
                  </dl>
                  <div className="slot-card-actions">
                    <button
                      type="button"
                      className={`button ${selectedSlotId === slot.slotId ? 'secondary' : 'primary'} slot-select-button`}
                      disabled={!slot.available}
                      onClick={() => setSelectedSlotId(slot.slotId)}
                    >
                      {slot.available
                        ? selectedSlotId === slot.slotId
                          ? '選択中'
                          : 'この枠を選択'
                        : '選択不可'}
                    </button>
                  </div>
                </article>
              ))}
            </div>

            {selectedService && selectedSlot ? (
              <section className="reservation-summary-card">
                <div>
                  <p className="service-card-label">Selection</p>
                  <h2>予約内容の確認へ進みます</h2>
                </div>
                <dl className="reservation-summary-grid">
                  <div>
                    <dt>サービス</dt>
                    <dd>{selectedService.name}</dd>
                  </div>
                  <div>
                    <dt>日付</dt>
                    <dd>{searchedDate}</dd>
                  </div>
                  <div>
                    <dt>開始時刻</dt>
                    <dd>{selectedSlot.startTime}</dd>
                  </div>
                  <div>
                    <dt>終了時刻</dt>
                    <dd>{selectedSlot.endTime}</dd>
                  </div>
                </dl>
                <div className="actions">
                  <button type="button" className="button primary" onClick={handleProceedToConfirm}>
                    予約内容を確認する
                  </button>
                </div>
              </section>
            ) : null}
          </>
        ) : (
          <div className="status-card">
            <p className="status-label">Reservation Slots</p>
            <p className="status-value">該当日の空き枠は見つかりませんでした。</p>
          </div>
        )
      ) : null}
    </section>
  )
}
