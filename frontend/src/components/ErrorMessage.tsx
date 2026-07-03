type ErrorMessageProps = {
  message: string
  title?: string
  className?: string
}

export function ErrorMessage({
  message,
  title = 'エラーが発生しました',
  className = '',
}: ErrorMessageProps) {
  return (
    <div className={`error-message ${className}`.trim()} role="alert" aria-live="polite">
      <p className="error-message-title">{title}</p>
      <p className="error-message-text">{message}</p>
    </div>
  )
}
