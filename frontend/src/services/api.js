import axios from 'axios'

// Базовый URL для API
const API_BASE_URL = 'http://localhost:8080/api'

// Создать axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Interceptor для добавления токена к каждому запросу
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Interceptor для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status

    if (status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    } else if (status === 403) {
      window.location.href = `/error?code=403&from=${encodeURIComponent(window.location.pathname)}`
    } else if (status === 500 || status === 502 || status === 503) {
      window.location.href = `/error?code=${status}&from=${encodeURIComponent(window.location.pathname)}`
    }

    return Promise.reject(error)
  }
)

// API методы
export default {
  // Auth
  register(data) {
    return api.post('/auth/register', data)
  },
  login(data) {
    return api.post('/auth/login', data)
  },
  getProfile() {
    return api.get('/profile')
  },
  updateProfile(data) {
    return api.put('/profile', data)
  },

  // Movies
  getMovies(params) {
    return api.get('/movies', { params })
  },
  getMovieDetails(id) {
    return api.get(`/movies/${id}`)
  },

  // Cinemas
  getCinemas(params) {
    return api.get('/cinemas', { params })
  },

  // Showtimes
  getShowtimes(params) {
    return api.get('/showtimes', { params })
  },

  // Bookings
  createBooking(data) {
    return api.post('/bookings', data)
  },
  getMyBookings(params) {
    return api.get('/bookings/my', { params })
  },
  confirmBooking(id) {
    return api.post(`/bookings/${id}/confirm`)
  },
  cancelBooking(id) {
    return api.delete(`/bookings/${id}`)
  },

  // Wallet
  topUpWallet(data) {
    return api.post('/wallet/topup', data)
  },

  // Analytics
  getPopularMovies(params) {
    return api.get('/analytics/popular-movies', { params })
  }
}