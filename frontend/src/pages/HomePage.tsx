export function HomePage() {
  return (
    <section className="login-stage">
      <div className="login-window">
        <div className="window-dots" aria-hidden="true">
          <span className="dot red"></span>
          <span className="dot yellow"></span>
          <span className="dot green"></span>
        </div>

        <div className="login-card">
          <p className="login-title">予約管理システム</p>
          <h1 className="login-heading">ログイン</h1>

          <form className="login-form">
            <label className="input-row">
              <span className="input-icon" aria-hidden="true">
                ✉
              </span>
              <input type="email" placeholder="メールアドレス" />
            </label>

            <label className="input-row">
              <span className="input-icon" aria-hidden="true">
                🔒
              </span>
              <input type="password" placeholder="パスワード" />
            </label>

            <button type="button" className="login-button">
              ログイン
            </button>
          </form>

          <a className="forgot-link" href="#">
            パスワードを忘れた方はこちら
          </a>
        </div>
      </div>
    </section>
  )
}
