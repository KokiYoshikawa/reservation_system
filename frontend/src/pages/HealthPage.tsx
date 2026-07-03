import { useEffect, useState } from 'react'
import { fetchHealth, type HealthResponse } from '../lib/api'

type RequestState = 'idle' | 'loading' | 'success' | 'error'

export function HealthPage() {
  const [state, setState] = useState<RequestState>('idle')
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true

    async function loadHealth() {
      setState('loading')
      setError('')

      try {
        const response = await fetchHealth()
        if (!active) {
          return
        }

        setHealth(response)
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

    void loadHealth()

    return () => {
      active = false
    }
  }, [])

  return (
    <section className="page">
      <span className="eyebrow">Backend Health</span>
      <h1>ヘルスチェックページ</h1>
      <p className="lead">
        Axios クライアント経由で `/api/v1/health` を確認しています。
      </p>
      <div className="status-card">
        <p className="status-label">API status</p>
        {state === 'loading' || state === 'idle' ? (
          <p className="status-value">Loading...</p>
        ) : null}
        {state === 'success' && health ? (
          <p className="status-value success">{health.status}</p>
        ) : null}
        {state === 'error' ? <p className="status-value error">{error}</p> : null}
      </div>
    </section>
  )
}
