type LoadingIndicatorProps = {
  label?: string
  className?: string
  compact?: boolean
}

export function LoadingIndicator({
  label = '読み込み中...',
  className = '',
  compact = false,
}: LoadingIndicatorProps) {
  return (
    <div
      className={`loading-indicator ${compact ? 'compact' : ''} ${className}`.trim()}
      role="status"
      aria-live="polite"
    >
      <span className="loading-spinner" aria-hidden="true"></span>
      <span className="loading-label">{label}</span>
    </div>
  )
}
