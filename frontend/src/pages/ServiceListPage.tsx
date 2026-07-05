import { useEffect, useState } from 'react'
import { ErrorMessage } from '../components/ErrorMessage'
import { LoadingIndicator } from '../components/LoadingIndicator'
import { fetchServices, type ServiceItem } from '../lib/api'

type RequestState = 'idle' | 'loading' | 'success' | 'error'

export function ServiceListPage() {
  const [state, setState] = useState<RequestState>('idle')
  const [services, setServices] = useState<ServiceItem[]>([])
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true

    async function loadServices() {
      setState('loading')
      setError('')

      try {
        const response = await fetchServices()
        if (!active) {
          return
        }

        setServices(response.services)
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

    void loadServices()

    return () => {
      active = false
    }
  }, [])

  return (
    <section className="page services-page">
      <span className="eyebrow">Services</span>
      <h1>サービス一覧</h1>
      <p className="lead">現在利用可能なサービスを一覧表示しています。</p>

      {state === 'loading' || state === 'idle' ? (
        <LoadingIndicator label="サービス一覧を取得しています..." className="status-feedback" />
      ) : null}

      {state === 'error' ? (
        <ErrorMessage
          message={error}
          title="サービス一覧の取得に失敗しました"
          className="status-feedback"
        />
      ) : null}

      {state === 'success' ? (
        services.length > 0 ? (
          <div className="service-list">
            {services.map((service) => (
              <article key={service.id} className="service-card">
                <div className="service-card-header">
                  <div>
                    <p className="service-card-label">Service</p>
                    <h2>{service.name}</h2>
                  </div>
                  <span className="service-card-badge">Active</span>
                </div>
                <dl className="service-meta">
                  <div>
                    <dt>所要時間</dt>
                    <dd>{service.duration_minutes}分</dd>
                  </div>
                  <div>
                    <dt>料金</dt>
                    <dd>{service.price.toLocaleString('ja-JP')}円</dd>
                  </div>
                </dl>
              </article>
            ))}
          </div>
        ) : (
          <div className="status-card">
            <p className="status-label">Services</p>
            <p className="status-value">公開中のサービスはまだありません。</p>
          </div>
        )
      ) : null}
    </section>
  )
}
