import type { ReservationSlotItem, ServiceItem } from './api'

export type ReservationConfirmRouteState = {
  service: Pick<ServiceItem, 'id' | 'name' | 'duration_minutes' | 'price'>
  slot: ReservationSlotItem
  date: string
}

export type ReservationCompleteRouteState = {
  reservationId: number
  status: string
  service: {
    id: number
    name: string
  }
  startTime: string
  endTime: string
  reservedAt: string
  note: string
}
