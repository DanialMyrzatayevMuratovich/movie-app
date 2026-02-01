<template>
  <div class="seat-map">
    <div class="screen">
      <div class="screen-line"></div>
      <div class="screen-text">ЭКРАН</div>
    </div>

    <div class="seats-container">
      <div 
        v-for="row in rows" 
        :key="row" 
        class="seat-row"
      >
        <div class="row-label">{{ row }}</div>
        
        <div class="seats">
          <div
            v-for="seat in getSeatsInRow(row)"
            :key="`${row}-${seat.number}`"
            class="seat"
            :class="getSeatClass(row, seat.number)"
            @click="toggleSeat(row, seat.number)"
          >
            <div class="seat-inner">{{ seat.number }}</div>
          </div>
        </div>

        <div class="row-label">{{ row }}</div>
      </div>
    </div>

    <!-- Legend -->
    <div class="legend">
      <div class="legend-item">
        <div class="seat seat-available"><div class="seat-inner"></div></div>
        <span>Доступно</span>
      </div>
      <div class="legend-item">
        <div class="seat seat-selected"><div class="seat-inner"></div></div>
        <span>Ваш выбор</span>
      </div>
      <div class="legend-item">
        <div class="seat seat-booked"><div class="seat-inner"></div></div>
        <span>Занято</span>
      </div>
      <div class="legend-item">
        <div class="seat seat-vip"><div class="seat-inner"></div></div>
        <span>VIP</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  seats: {
    type: Array,
    required: true
  },
  bookedSeats: {
    type: Array,
    default: () => []
  },
  selectedSeats: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['seat-selected'])

const rows = computed(() => {
  return [...new Set(props.seats.map(s => s.row))].sort()
})

const getSeatsInRow = (row) => {
  return props.seats
    .filter(s => s.row === row)
    .sort((a, b) => a.number - b.number)
}

const isSeatBooked = (row, number) => {
  return props.bookedSeats.some(s => s.row === row && s.number === number)
}

const isSeatSelected = (row, number) => {
  return props.selectedSeats.some(s => s.row === row && s.number === number)
}

const getSeatType = (row, number) => {
  const seat = props.seats.find(s => s.row === row && s.number === number)
  return seat?.type || 'regular'
}

const getSeatClass = (row, number) => {
  if (isSeatBooked(row, number)) return 'seat-booked'
  if (isSeatSelected(row, number)) return 'seat-selected'
  
  const type = getSeatType(row, number)
  if (type === 'vip') return 'seat-vip seat-available'
  if (type === 'couple') return 'seat-couple seat-available'
  
  return 'seat-available'
}

const toggleSeat = (row, number) => {
  if (isSeatBooked(row, number)) return
  
  const seat = props.seats.find(s => s.row === row && s.number === number)
  emit('seat-selected', { row, number, price: seat.price })
}
</script>

<style scoped>
.seat-map {
  max-width: 900px;
  margin: 0 auto;
}

.screen {
  text-align: center;
  margin-bottom: 40px;
}

.screen-line {
  width: 80%;
  height: 4px;
  background: linear-gradient(to right, transparent, var(--primary), transparent);
  margin: 0 auto 10px;
  border-radius: 50%;
  box-shadow: 0 0 20px rgba(229, 9, 20, 0.5);
}

.screen-text {
  color: var(--text-gray);
  font-size: 14px;
  font-weight: bold;
  letter-spacing: 2px;
}

.seats-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 30px;
}

.seat-row {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
}

.row-label {
  min-width: 30px;
  text-align: center;
  font-weight: bold;
  color: var(--text-gray);
  font-size: 14px;
}

.seats {
  display: flex;
  gap: 8px;
}

.seat {
  width: 36px;
  height: 36px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.seat-inner {
  width: 100%;
  height: 100%;
  border-radius: 8px 8px 0 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: bold;
  position: relative;
}

.seat-inner::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 50%;
  transform: translateX(-50%);
  width: 70%;
  height: 2px;
  border-radius: 2px;
}

.seat-available .seat-inner {
  background-color: var(--dark-lighter);
  color: var(--text-gray);
  border: 2px solid var(--dark-lighter);
}

.seat-available .seat-inner::after {
  background-color: var(--dark-lighter);
}

.seat-available:hover .seat-inner {
  background-color: var(--secondary);
  color: var(--dark);
  border-color: var(--secondary);
  transform: scale(1.1);
}

.seat-selected .seat-inner {
  background-color: var(--primary);
  color: white;
  border: 2px solid var(--primary);
  transform: scale(1.05);
}

.seat-selected .seat-inner::after {
  background-color: var(--primary);
}

.seat-booked .seat-inner {
  background-color: var(--dark);
  color: var(--text-gray);
  border: 2px solid var(--dark);
  cursor: not-allowed;
  opacity: 0.5;
}

.seat-vip.seat-available .seat-inner {
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  color: var(--dark);
  border: 2px solid #ffd700;
}

.seat-vip.seat-available .seat-inner::after {
  background: linear-gradient(135deg, #ffd700, #ffed4e);
}

.seat-couple {
  width: 76px;
}

.legend {
  display: flex;
  justify-content: center;
  gap: 24px;
  flex-wrap: wrap;
  padding: 20px;
  background-color: var(--dark-light);
  border-radius: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-gray);
}

.legend-item .seat {
  cursor: default;
}

@media (max-width: 768px) {
  .seat {
    width: 28px;
    height: 28px;
  }
  
  .seat-inner {
    font-size: 9px;
  }
  
  .seats {
    gap: 4px;
  }
}
</style>