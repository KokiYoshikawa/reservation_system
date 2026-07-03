import { NavLink, Route, Routes } from 'react-router-dom'
import './App.css'
import { useAuth } from './context/AuthContext'
import { HealthPage } from './pages/HealthPage'
import { HomePage } from './pages/HomePage'
import { NotFoundPage } from './pages/NotFoundPage'
import { ReservationsPage } from './pages/ReservationsPage'

function App() {
  const { isAuthenticated, isLoading, logout, user } = useAuth()

  return (
    <div className="app-shell">
      <header className="site-header">
        <NavLink className="brand" to="/">
          Reservation System
        </NavLink>
        <nav className="site-nav" aria-label="global">
          <NavLink
            className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}
            to="/"
            end
          >
            Home
          </NavLink>
          <NavLink
            className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}
            to="/reservations"
          >
            Reservations
          </NavLink>
          <NavLink
            className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}
            to="/health"
          >
            Health
          </NavLink>
        </nav>
        <div className="auth-panel">
          {isLoading ? <span className="auth-meta">認証状態を確認中...</span> : null}
          {!isLoading && isAuthenticated && user ? (
            <>
              <span className="auth-meta">{user.email}</span>
              <button type="button" className="auth-button" onClick={logout}>
                ログアウト
              </button>
            </>
          ) : null}
          {!isLoading && !isAuthenticated ? (
            <span className="auth-meta">未ログイン</span>
          ) : null}
        </div>
      </header>

      <main className="content">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/reservations" element={<ReservationsPage />} />
          <Route path="/health" element={<HealthPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>
    </div>
  )
}

export default App
