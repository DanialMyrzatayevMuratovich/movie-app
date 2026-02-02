<template>
  <div class="booking-card">
    <div class="booking-header">
      <div class="booking-number">{{ booking.bookingNumber }}</div>
      <div class="booking-status" :class="`status-${booking.status}`">
        <span class="status-dot"></span>
        {{ getStatusText(booking.status) }}
      </div>
    </div>

    <div class="booking-content">
      <div class="booking-info">
        <div class="info-item">
          <span class="info-label">🎬 Фильм:</span>
          <span class="info-value">{{ movieTitle }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">📅 Дата и время:</span>
          <span class="info-value">{{ formatDate(booking.createdAt) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">🎫 Места:</span>
          <span class="info-value">{{ getSeatsText() }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">💳 Оплата:</span>
          <span class="info-value">{{ getPaymentMethod() }}</span>
        </div>
      </div>

      <div class="booking-amount">
        <div class="amount-label">Итого:</div>
        <div class="amount-value">{{ formatCurrency(booking.totalAmount) }}</div>
      </div>
    </div>

    <div class="booking-actions">
      <button 
        v-if="booking.status === 'pending'" 
        class="btn btn-primary btn-sm"
        @click="$emit('confirm', booking._id || booking.id)"
      >
        Подтвердить оплату
      </button>
      
      <button 
        v-if="booking.status === 'confirmed'" 
        class="btn btn-outline btn-sm"
        @click="showQR = !showQR"
      >
        {{ showQR ? 'Скрыть QR' : 'Показать QR' }}
      </button>

      <button 
        v-if="booking.status !== 'cancelled'" 
        class="btn btn-secondary btn-sm"
        @click="$emit('cancel', booking._id || booking.id)"
      >
        Отменить
      </button>
    </div>

    <div v-if="showQR && booking.status === 'confirmed'" class="qr-code">
      <div class="qr-content">
        <img v-if="qrDataUrl" :src="qrDataUrl" alt="QR Code" class="qr-image" />
        <div v-else class="qr-loading">Генерация QR...</div>
        <div class="qr-text">{{ booking.qrCode || booking.bookingNumber }}</div>
        <p class="qr-hint">Покажите этот код на входе</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import QRCode from 'qrcode'
import { formatDate, formatCurrency } from '../utils/formatters'

const props = defineProps({
  booking: {
    type: Object,
    required: true
  },
  movieTitle: {
    type: String,
    default: 'Загрузка...'
  }
})

defineEmits(['confirm', 'cancel'])

const showQR = ref(false)
const qrDataUrl = ref('')

watch(showQR, async (visible) => {
  if (visible && !qrDataUrl.value) {
    const qrData = props.booking.qrCode || props.booking.bookingNumber || (props.booking._id || props.booking.id)
    try {
      qrDataUrl.value = await QRCode.toDataURL(String(qrData), {
        width: 200,
        margin: 2,
        color: { dark: '#000000', light: '#ffffff' }
      })
    } catch (err) {
      console.error('QR generation failed:', err)
    }
  }
})

const getStatusText = (status) => {
  const statusMap = {
    pending: 'Ожидает оплаты',
    confirmed: 'Подтверждено',
    cancelled: 'Отменено',
    expired: 'Истекло'
  }
  return statusMap[status] || status
}

const getSeatsText = () => {
  return props.booking.seats
    .map(s => `${s.row}-${s.number}`)
    .join(', ')
}

const getPaymentMethod = () => {
  const methods = {
    wallet: 'Кошелек',
    card: 'Карта',
    cash: 'Наличные'
  }
  return methods[props.booking.payment?.method] || props.booking.payment?.method
}
</script>

<style scoped>
.booking-card {
  background-color: var(--dark-light);
  border-radius: 12px;
  padding: 20px;
  border: 2px solid var(--dark-lighter);
  transition: all 0.3s ease;
}

.booking-card:hover {
  border-color: var(--primary);
}

.booking-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.booking-number {
  font-weight: bold;
  font-size: 18px;
  color: var(--text);
}

.booking-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-pending {
  background-color: rgba(255, 193, 7, 0.2);
  color: var(--warning);
}

.status-pending .status-dot {
  background-color: var(--warning);
}

.status-confirmed {
  background-color: rgba(70, 211, 105, 0.2);
  color: var(--success);
}

.status-confirmed .status-dot {
  background-color: var(--success);
}

.status-cancelled, .status-expired {
  background-color: rgba(229, 9, 20, 0.2);
  color: var(--danger);
}

.status-cancelled .status-dot,
.status-expired .status-dot {
  background-color: var(--danger);
}

.booking-content {
  margin-bottom: 16px;
}

.booking-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.info-label {
  color: var(--text-gray);
  font-size: 14px;
}

.info-value {
  color: var(--text);
  font-weight: 500;
  text-align: right;
}

.booking-amount {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: var(--dark-lighter);
  border-radius: 8px;
}

.amount-label {
  color: var(--text-gray);
  font-size: 14px;
}

.amount-value {
  font-size: 24px;
  font-weight: bold;
  color: var(--primary);
}

.booking-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.btn-sm {
  padding: 8px 16px;
  font-size: 14px;
}

.qr-code {
  margin-top: 16px;
  padding: 20px;
  background-color: var(--dark-lighter);
  border-radius: 8px;
  text-align: center;
}

.qr-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.qr-image {
  width: 200px;
  height: 200px;
  border-radius: 8px;
}

.qr-loading {
  width: 200px;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-gray);
  font-size: 14px;
}

.qr-text {
  font-family: monospace;
  font-size: 18px;
  font-weight: bold;
  color: var(--text);
  padding: 12px 20px;
  background-color: var(--dark);
  border-radius: 8px;
}

.qr-hint {
  color: var(--text-gray);
  font-size: 13px;
}

@media (max-width: 768px) {
  .info-item {
    flex-direction: column;
    gap: 4px;
  }
  
  .info-value {
    text-align: left;
  }
}
</style>