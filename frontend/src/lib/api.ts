import axios from 'axios'

export type HealthResponse = {
  status: string
  error?: string
}

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use((config) => {
  const storedValue = window.localStorage.getItem('reservation_system_auth')

  if (!storedValue) {
    return config
  }

  try {
    const parsed = JSON.parse(storedValue) as { token?: string }
    if (parsed.token) {
      config.headers.Authorization = `Bearer ${parsed.token}`
    }
  } catch {
    window.localStorage.removeItem('reservation_system_auth')
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

export { apiClient }
