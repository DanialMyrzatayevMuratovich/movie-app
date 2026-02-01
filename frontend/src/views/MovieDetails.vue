<template>
  <div class="movie-details">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <div v-else-if="movie" class="container">
      <!-- Hero Section -->
      <div class="movie-hero">
        <div class="movie-backdrop">
          <div class="backdrop-overlay"></div>
        </div>

        <div class="movie-hero-content">
          <div class="movie-poster-large">
            <img
              v-if="movie.posterFileId"
              :src="`http://localhost:8080/api/files/${movie.posterFileId}`"
              :alt="movie.title"
              @error="handleImageError"
            />
            <div v-else class="poster-placeholder">
              <span class="poster-icon">🎬</span>
            </div>
          </div>

          <div class="movie-info-main">
            <h1 class="movie-title">{{ movie.title }}</h1>
            <p class="movie-subtitle">{{ movie.titleRu }}</p>

            <div class="movie-meta-main">
              <div class="rating-large">
                <span class="star-large">⭐</span>
                <span class="rating-value">{{ movie.imdbRating }}</span>
                <span class="rating-label">IMDb</span>
              </div>

              <div class="meta-items">
                <span class="meta-tag">{{ movie.rating }}</span>
                <span class="meta-tag">{{ formatDuration(movie.duration) }}</span>
                <span class="meta-tag">{{ movie.ageRestriction }}+</span>
              </div>
            </div>

            <div class="genres">
              <span v-for="genre in movie.genres" :key="genre" class="genre-tag">
                {{ genre }}
              </span>
            </div>

            <p class="movie-description">{{ movie.description }}</p>

            <div class="movie-details-info">
              <div class="detail-item">
                <strong>Режиссер:</strong> {{ movie.director }}
              </div>
              <div class="detail-item">
                <strong>В ролях:</strong> {{ movie.cast?.join(', ') }}
              </div>
              <div class="detail-item">
                <strong>Дата выхода:</strong> {{ formatDate(movie.releaseDate) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Showtimes Section -->
      <div class="showtimes-section">
        <h2 class="section-title">Выберите сеанс</h2>

        <div v-if="loadingShowtimes" class="loading">
          <div class="spinner"></div>
        </div>

        <div v-else-if="showtimes.length > 0">
          <!-- Date Filter -->
          <div class="date-filters">
            <button
              v-for="date in availableDates"
              :key="date"
              class="date-btn"
              :class="{ active: selectedDate === date }"
              @click="selectDate(date)"
            >
              {{ formatShortDate(date) }}
            </button>
          </div>

          <!-- Showtimes by Cinema -->
          <div class="showtimes-list">
            <div
              v-for="cinema in groupedShowtimes"
              :key="cinema.cinemaId"
              class="cinema-showtimes"
            >
              <h3 class="cinema-name">🏢 {{ cinema.cinemaName }}</h3>

              <div class="time-slots">
                <button
                  v-for="showtime in cinema.showtimes"
                  :key="showtime._id"
                  class="time-slot"
                  @click="selectShowtime(showtime)"
                >
                  <div class="time">{{ formatTime(showtime.startTime) }}</div>
                  <div class="format-price">
                    <span class="format">{{ showtime.format }}</span>
                    <span class="price">{{ formatCurrency(showtime.basePrice) }}</span>
                  </div>
                  <div class="seats-available">
                    {{ showtime.availableSeats }} мест
                  </div>
                </button>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="no-showtimes">
          <p>Нет доступных сеансов на выбранную дату</p>
        </div>
      </div>

      <!-- Reviews Section -->
      <div v-if="movie.reviews && movie.reviews.length > 0" class="reviews-section">
        <h2 class="section-title">Отзывы ({{ movie.reviews.length }})</h2>

        <div class="reviews-list">
          <div v-for="review in movie.reviews" :key="review.userId" class="review-card">
            <div class="review-header">
              <div class="review-rating">
                ⭐ {{ review.rating }}/10
              </div>
              <div class="review-date">{{ formatDate(review.createdAt) }}</div>
            </div>
            <p class="review-comment">{{ review.comment }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../services/api'
import { formatDate, formatTime, formatDuration, formatCurrency } from '../utils/formatters'

const route = useRoute()
const router = useRouter()

const movie = ref(null)
const showtimes = ref([])
const loading = ref(true)
const loadingShowtimes = ref(true)
const selectedDate = ref(null)

const availableDates = computed(() => {
  const dates = new Set()
  showtimes.value.forEach(st => {
    const date = new Date(st.startTime).toISOString().split('T')[0]
    dates.add(date)
  })
  return Array.from(dates).sort().slice(0, 7)
})

const groupedShowtimes = computed(() => {
  if (!selectedDate.value) return []

  const filtered = showtimes.value.filter(st => {
    const date = new Date(st.startTime).toISOString().split('T')[0]
    return date === selectedDate.value
  })

  const grouped = {}
  filtered.forEach(st => {
    if (!grouped[st.cinemaId]) {
      grouped[st.cinemaId] = {
        cinemaId: st.cinemaId,
        cinemaName: st.cinemaDetails?.name || 'Кинотеатр',
        showtimes: []
      }
    }
    grouped[st.cinemaId].showtimes.push(st)
  })

  return Object.values(grouped)
})

const fetchMovie = async () => {
  try {
    const response = await api.getMovieDetails(route.params.id)
    movie.value = response.data.data.movie || response.data.data
  } catch (error) {
    console.error('Failed to fetch movie:', error)
  } finally {
    loading.value = false
  }
}

const fetchShowtimes = async () => {
  loadingShowtimes.value = true
  try {
    const response = await api.getShowtimes({
      movieId: route.params.id,
      onlyFuture: 'true',
      includeDetails: 'true'
    })
    showtimes.value = response.data.data

    if (availableDates.value.length > 0) {
      selectedDate.value = availableDates.value[0]
    }
  } catch (error) {
    console.error('Failed to fetch showtimes:', error)
  } finally {
    loadingShowtimes.value = false
  }
}

const selectDate = (date) => {
  selectedDate.value = date
}

const selectShowtime = (showtime) => {
  router.push(`/booking/${showtime._id}`)
}

const formatShortDate = (dateStr) => {
  const date = new Date(dateStr)
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)

  if (date.toDateString() === today.toDateString()) {
    return 'Сегодня'
  } else if (date.toDateString() === tomorrow.toDateString()) {
    return 'Завтра'
  }

  return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
}

const handleImageError = (e) => {
  e.target.style.display = 'none'
}

onMounted(() => {
  fetchMovie()
  fetchShowtimes()
})
</script>

<style scoped>
.movie-details {
  min-height: calc(100vh - 70px);
  padding-bottom: 40px;
}

.movie-hero {
  position: relative;
  margin-bottom: 40px;
}

.movie-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 500px;
  background: linear-gradient(135deg, var(--dark) 0%, var(--dark-light) 100%);
  z-index: 0;
}

.backdrop-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 200px;
  background: linear-gradient(to bottom, transparent, var(--dark));
}

.movie-hero-content {
  position: relative;
  z-index: 1;
  display: flex;
  gap: 40px;
  padding: 40px 0;
}

.movie-poster-large {
  flex-shrink: 0;
  width: 300px;
  aspect-ratio: 2/3;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.movie-poster-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.poster-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--dark-lighter), var(--dark));
}

.poster-icon {
  font-size: 80px;
  opacity: 0.3;
}

.movie-info-main {
  flex: 1;
}

.movie-title {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 8px;
}

.movie-subtitle {
  font-size: 24px;
  color: var(--text-gray);
  margin-bottom: 24px;
}

.movie-meta-main {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 20px;
}

.rating-large {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: var(--dark-light);
  padding: 12px 20px;
  border-radius: 12px;
}

.star-large {
  font-size: 32px;
}

.rating-value {
  font-size: 32px;
  font-weight: bold;
}

.rating-label {
  font-size: 14px;
  color: var(--text-gray);
}

.meta-items {
  display: flex;
  gap: 12px;
}

.meta-tag {
  padding: 8px 16px;
  background-color: var(--dark-lighter);
  border-radius: 8px;
  font-weight: 600;
}

.genres {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.genre-tag {
  padding: 6px 16px;
  background-color: var(--primary);
  color: white;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
}

.movie-description {
  font-size: 16px;
  line-height: 1.8;
  color: var(--text-gray);
  margin-bottom: 24px;
}

.movie-details-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-item {
  font-size: 15px;
  color: var(--text);
}

.detail-item strong {
  color: var(--text-gray);
  margin-right: 8px;
}

.showtimes-section,
.reviews-section {
  margin-top: 60px;
}

.section-title {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 24px;
}

.date-filters {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.date-btn {
  padding: 12px 24px;
  background-color: var(--dark-light);
  color: var(--text);
  border: 2px solid var(--dark-lighter);
  border-radius: 12px;
  font-weight: 600;
  white-space: nowrap;
  transition: all 0.3s ease;
}

.date-btn:hover {
  border-color: var(--primary);
}

.date-btn.active {
  background-color: var(--primary);
  border-color: var(--primary);
  color: white;
}

.showtimes-list {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.cinema-showtimes {
  background-color: var(--dark-light);
  padding: 24px;
  border-radius: 12px;
}

.cinema-name {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 16px;
}

.time-slots {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 12px;
}

.time-slot {
  padding: 16px;
  background-color: var(--dark-lighter);
  border: 2px solid transparent;
  border-radius: 12px;
  transition: all 0.3s ease;
  text-align: center;
}

.time-slot:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(229, 9, 20, 0.3);
}

.time {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 8px;
}

.format-price {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
  font-size: 13px;
}

.format {
  color: var(--secondary);
  font-weight: 600;
}

.price {
  color: var(--text-gray);
}

.seats-available {
  font-size: 12px;
  color: var(--text-gray);
}

.no-showtimes {
  text-align: center;
  padding: 40px;
  color: var(--text-gray);
}

.reviews-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.review-card {
  background-color: var(--dark-light);
  padding: 20px;
  border-radius: 12px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.review-rating {
  font-weight: bold;
  color: var(--secondary);
}

.review-date {
  font-size: 13px;
  color: var(--text-gray);
}

.review-comment {
  color: var(--text);
  line-height: 1.6;
}

@media (max-width: 768px) {
  .movie-hero-content {
    flex-direction: column;
    align-items: center;
  }

  .movie-poster-large {
    width: 100%;
    max-width: 300px;
  }

  .movie-title {
    font-size: 32px;
  }

  .movie-subtitle {
    font-size: 18px;
  }

  .time-slots {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }
}
</style>