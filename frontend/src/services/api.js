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
    if (error.response?.status === 401) {
      // Токен истек или невалиден
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
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

  // Analytics
  getPopularMovies(params) {
    return api.get('/analytics/popular-movies', { params })
  }
}