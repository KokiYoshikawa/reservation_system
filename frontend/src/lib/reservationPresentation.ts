const dateTimeFormatter = new Intl.DateTimeFormat('ja-JP', {
  year: 'numeric',
  month: '2-digit',
  day: '2-digit',
  weekday: 'short',
  hour: '2-digit',
  minute: '2-digit',
})

const dateFormatter = new Intl.DateTimeFormat('ja-JP', {
  year: 'numeric',
  month: 'long',
  day: 'numeric',
  weekday: 'short',
})

const timeFormatter = new Intl.DateTimeFormat('ja-JP', {
  hour: '2-digit',
  minute: '2-digit',
})

export function formatReservationDateTime(value: string) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : dateTimeFormatter.format(date)
}

export function formatReservationDate(value: string) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : dateFormatter.format(date)
}

export function formatReservationTime(value: string) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : timeFormatter.format(date)
}

export function reservationStatus(status: string) {
  switch (status) {
    case 'reserved':
      return { label: '予約済み', className: 'reserved' }
    case 'cancelled':
      return { label: 'キャンセル済み', className: 'cancelled' }
    case 'completed':
      return { label: '完了', className: 'completed' }
    case 'no_show':
      return { label: '来店なし', className: 'no-show' }
    default:
      return { label: status, className: 'unknown' }
  }
}
