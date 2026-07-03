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
