<template>
  <div class="seat-selection">
    <div class="container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
      </div>

      <div v-else-if="showtime && hall">
        <!-- Header -->
        <div class="booking-header">
          <button class="btn btn-secondary" @click="goBack">
            ← Назад
          </button>
          <h1 class="page-title">Выбор мест</h1>
        </div>

        <!-- Movie Info -->
        <div class="movie-info-card">
          <div class="movie-info-content">
            <h2 class="movie-title">{{ movieTitle }}</h2>
            <div class="session-info">
              <span class="info-item">📅 {{ formatDate(showtime.startTime) }}</span>
              <span class="info-item">🏢 {{ cinemaName }}</span>
              <span class="info-item">🎬 {{ showtime.format }}</span>
              <span class="info-item">🗣️ {{ showtime.language }}</span>
            </div>
          </div>
        </div>

        <!-- Seat Map -->
        <div class="seat-selection-area">
          <SeatMap
            :seats="hall.seats"
            :booked-seats="showtime.bookedSeats"
            :selected-seats="selectedSeats"
            @seat-selected="toggleSeat"
          />
        </div>

        <!-- Booking Summary -->
        <div v-if="selectedSeats.length > 0" class="booking-summary">
          <div class="summary-content">
            <div class="selected-seats-info">
              <h3>Выбранные места:</h3>
              <div class="seats-list">
                <span
                  v-for="seat in selectedSeats"
                  :key="`${seat.row}-${seat.number}`"
                  class="seat-badge"
                >
                  {{ seat.row }}-{{ seat.number }}
                  <button class="remove-seat" @click="removeSeat(seat)">×</button>
                </span>
              </div>
            </div>

            <div class="total-price">
              <span class="price-label">Итого:</span>
              <span class="price-value">{{ formatCurrency(totalPrice) }}</span>
            </div>
          </div>

          <div class="payment-method">
            <h3>Способ оплаты:</h3>
            <div class="payment-options">
              <label class="payment-option">
                <input
                  v-model="paymentMethod"
                  type="radio"
                  value="wallet"
                  name="payment"
                />
                <span class="option-content">
                  <span class="option-icon">💳</span>
                  <span class="option-text">
                    <strong>Кошелек</strong>
                    <small>{{ formatCurrency(userBalance) }}</small>
                  </span>
                </span>
              </label>

              <label class="payment-option">
                <input
                  v-model="paymentMethod"
                  type="radio"
                  value="card"
                  name="payment"
                />
                <span class="option-content">
                  <span class="option-icon">💳</span>
                  <span class="option-text">
                    <strong>Банковская карта</strong>
                    <small>Visa / MasterCard</small>
                  </span>
                </span>
              </label>

              <label class="payment-option">
                <input
                  v-model="paymentMethod"
                  type="radio"
                  value="cash"
                  name="payment"
                />
                <span class="option-content">
                  <span class="option-icon">💵</span>
                  <span class="option-text">
                    <strong>Наличные</strong>
                    <small>Оплата в кассе</small>
                  </span>
                </span>
              </label>
            </div>
          </div>

          <div v-if="error" class="alert alert-error">
            {{ error }}
          </div>

          <button
            class="btn btn-primary btn-large"
            :disabled="bookingInProgress"
            @click="createBooking"
          >
            <span v-if="bookingInProgress">Бронирование...</span>
            <span v-else>Забронировать {{ formatCurrency(totalPrice) }}</span>
          </button>
        </div>

        <div v-else class="no-seats-selected">
          <p>Выберите места для бронирования</p>
        </div>
      </div>

      <div v-else class="error-state">
        <h2>Ошибка загрузки</h2>
        <p>Не удалось загрузить информацию о сеансе</p>
        <button class="btn btn-primary" @click="goBack">
          Вернуться назад
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import SeatMap from '../components/SeatMap.vue'
import api from '../services/api'
import { formatDate, formatCurrency } from '../utils/formatters'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const bookingInProgress = ref(false)
const showtime = ref(null)
const hall = ref(null)
const movieTitle = ref('')
const cinemaName = ref('')
const selectedSeats = ref([])
const paymentMethod = ref('wallet')
const error = ref('')

const userBalance = computed(() => {
  return authStore.user?.wallet?.balance || 0
})

const totalPrice = computed(() => {
  return selectedSeats.value.reduce((sum, seat) => sum + seat.price, 0)
})

const fetchShowtimeDetails = async () => {
  loading.value = true
  try {
    // Получить информацию о сеансе
    const showtimeResponse = await api.getShowtimes({
      includeDetails: 'true'
    })
    
    const showtimeData = showtimeResponse.data.data.find(
      st => st._id === route.params.showtimeId
    )

    if (!showtimeData) {
      throw new Error('Showtime not found')
    }

    showtime.value = showtimeData

    // Извлечь информацию из includeDetails
    if (showtimeData.movieDetails) {
      movieTitle.value = showtimeData.movieDetails.title
    }
    if (showtimeData.cinemaDetails) {
      cinemaName.value = showtimeData.cinemaDetails.name
    }

    // Получить информацию о зале (нужно создать endpoint или использовать существующие данные)
    // Для простоты создадим mock данные на основе bookedSeats
    hall.value = generateHallSeats(showtime.value)

  } catch (err) {
    console.error('Failed to fetch showtime:', err)
    error.value = 'Не удалось загрузить информацию о сеансе'
  } finally {
    loading.value = false
  }
}

const generateHallSeats = (showtime) => {
  // Генерация мест зала (10 рядов по 15 мест)
  const rows = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J']
  const seatsPerRow = 15
  const seats = []

  rows.forEach((row, rowIndex) => {
    for (let number = 1; number <= seatsPerRow; number++) {
      let type = 'regular'
      let price = showtime.basePrice

      // VIP места (середина зала, средние ряды)
      if (rowIndex >= 3 && rowIndex <= 6 && number >= 5 && number <= 11) {
        type = 'vip'
        price = showtime.basePrice * 1.5
      }

      // Couple seats (последние 2 ряда, парные места)
      if (rowIndex >= 8 && number % 2 === 1 && number < seatsPerRow) {
        type = 'couple'
        price = showtime.basePrice * 1.3
      }

      seats.push({ row, number, type, price })
    }
  })

  return {
    _id: showtime.hallId,
    seats
  }
}

const toggleSeat = (seat) => {
  const index = selectedSeats.value.findIndex(
    s => s.row === seat.row && s.number === seat.number
  )

  if (index > -1) {
    selectedSeats.value.splice(index, 1)
  } else {
    if (selectedSeats.value.length >= 10) {
      error.value = 'Максимум 10 мест за одно бронирование'
      return
    }
    selectedSeats.value.push(seat)
    error.value = ''
  }
}

const removeSeat = (seat) => {
  const index = selectedSeats.value.findIndex(
    s => s.row === seat.row && s.number === seat.number
  )
  if (index > -1) {
    selectedSeats.value.splice(index, 1)
  }
}

const createBooking = async () => {
  if (selectedSeats.value.length === 0) {
    error.value = 'Выберите хотя бы одно место'
    return
  }

  if (paymentMethod.value === 'wallet' && totalPrice.value > userBalance.value) {
    error.value = `Недостаточно средств на кошельке. Необходимо: ${formatCurrency(totalPrice.value)}, доступно: ${formatCurrency(userBalance.value)}`
    return
  }

  bookingInProgress.value = true
  error.value = ''

  try {
    const bookingData = {
      showtimeId: route.params.showtimeId,
      seats: selectedSeats.value.map(s => ({
        row: s.row,
        number: s.number
      })),
      paymentMethod: paymentMethod.value
    }

    const response = await api.createBooking(bookingData)

    // Успешно создано
    const booking = response.data.data

    // Обновить профиль (баланс изменился)
    await authStore.fetchProfile()

    // Перейти в профиль
    router.push({
      name: 'Profile',
      query: { bookingCreated: booking._id }
    })
  } catch (err) {
    error.value = err.response?.data?.error || 'Не удалось создать бронь'
    console.error('Booking failed:', err)
  } finally {
    bookingInProgress.value = false
  }
}

const goBack = () => {
  router.go(-1)
}

onMounted(() => {
  fetchShowtimeDetails()
})
</script>

<style scoped>
.seat-selection {
  padding: 40px 0;
  min-height: calc(100vh - 70px);
}

.booking-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.page-title {
  font-size: 32px;
  font-weight: bold;
}

.movie-info-card {
  background: linear-gradient(135deg, var(--dark-light), var(--dark-lighter));
  padding: 24px;
  border-radius: 16px;
  margin-bottom: 40px;
  border: 2px solid var(--dark-lighter);
}

.movie-title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 12px;
}

.session-info {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  font-size: 14px;
  color: var(--text-gray);
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.seat-selection-area {
  margin-bottom: 40px;
}

.booking-summary {
  background-color: var(--dark-light);
  padding: 30px;
  border-radius: 16px;
  border: 2px solid var(--primary);
}

.summary-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 20px;
}

.selected-seats-info h3 {
  font-size: 18px;
  margin-bottom: 12px;
}

.seats-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.seat-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: var(--primary);
  color: white;
  border-radius: 20px;
  font-weight: 600;
  font-size: 14px;
}

.remove-seat {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  font-size: 16px;
  line-height: 1;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.remove-seat:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.total-price {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.price-label {
  font-size: 14px;
  color: var(--text-gray);
  margin-bottom: 4px;
}

.price-value {
  font-size: 36px;
  font-weight: bold;
  color: var(--primary);
}

.payment-method {
  margin-bottom: 24px;
}

.payment-method h3 {
  font-size: 18px;
  margin-bottom: 16px;
}

.payment-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.payment-option {
  position: relative;
  cursor: pointer;
}

.payment-option input[type="radio"] {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.option-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background-color: var(--dark-lighter);
  border: 2px solid var(--dark-lighter);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.payment-option input[type="radio"]:checked + .option-content {
  border-color: var(--primary);
  background-color: rgba(229, 9, 20, 0.1);
}

.option-icon {
  font-size: 32px;
}

.option-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.option-text strong {
  font-size: 16px;
}

.option-text small {
  font-size: 13px;
  color: var(--text-gray);
}

.btn-large {
  width: 100%;
  padding: 16px 32px;
  font-size: 18px;
  font-weight: bold;
}

.no-seats-selected {
  text-align: center;
  padding: 60px 20px;
  background-color: var(--dark-light);
  border-radius: 16px;
  color: var(--text-gray);
}

.error-state {
  text-align: center;
  padding: 80px 20px;
}

.error-state h2 {
  font-size: 32px;
  margin-bottom: 16px;
}

.error-state p {
  color: var(--text-gray);
  margin-bottom: 24px;
}

@media (max-width: 768px) {
  .summary-content {
    flex-direction: column;
    align-items: stretch;
  }

  .total-price {
    align-items: flex-start;
  }

  .payment-options {
    grid-template-columns: 1fr;
  }
}
</style>