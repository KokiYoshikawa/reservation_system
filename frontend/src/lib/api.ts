import axios from 'axios'

export type HealthResponse = {
  status: string
  error?: string
}

export type ServiceItem = {
  id: number
  name: string
  duration_minutes: number
  price: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export type ServicesResponse = {
  services: ServiceItem[]
}

export type ReservationSlotItem = {
  slotId: number
  startTime: string
  endTime: string
  capacity: number
  reservedCount: number
  available: boolean
}

export type ReservationSlotsResponse = {
  data: {
    date: string
    slots: ReservationSlotItem[]
  }
}

export type AuthUser = {
  email: string
  id?: number
  name?: string
  role?: string
}

export type LoginResponse = {
  token: string
  user: AuthUser
}

export type MeResponse = {
  user: AuthUser
}

export type CreateUserApiRequest = {
  name: string
  email: string
  password_hash: string
  role: string
}

const AUTH_STORAGE_KEY = 'reservation_system_auth'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use((config) => {
  const storedValue = window.localStorage.getItem(AUTH_STORAGE_KEY)

  if (!storedValue) {
    return config
  }

  try {
    const parsed = JSON.parse(storedValue) as { token?: string }
    if (parsed.token) {
      config.headers.Authorization = `Bearer ${parsed.token}`
    }
  } catch {
    window.localStorage.removeItem(AUTH_STORAGE_KEY)
  }

  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error)) {
      const message =
        error.response?.data?.error ??
        error.message ??
        'API request failed'

      return Promise.reject(new Error(message))
    }

    return Promise.reject(error)
  },
)

export async function fetchHealth() {
  const response = await apiClient.get<HealthResponse>('/health')
  return response.data
}

export async function fetchServices() {
  const response = await apiClient.get<ServicesResponse>('/services')
  return response.data
}

export async function fetchReservationSlots(date: string, serviceId: number) {
  const response = await apiClient.get<ReservationSlotsResponse>('/reservation-slots', {
    params: {
      date,
      serviceId,
    },
  })
  return response.data
}

export async function loginRequest(email: string, password: string) {
  const response = await apiClient.post<LoginResponse>('/auth/login', {
    email,
    password,
  })

  return response.data
}

export async function fetchCurrentUser() {
  const response = await apiClient.get<MeResponse>('/auth/me')
  return response.data
}

export async function logoutRequest() {
  await apiClient.post('/auth/logout')
}

export async function createUserRequest(payload: CreateUserApiRequest) {
  const response = await apiClient.post('/auth/register', payload)
  return response.data
}

export { apiClient, AUTH_STORAGE_KEY }
