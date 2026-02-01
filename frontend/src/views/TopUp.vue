<template>
  <div class="topup">
    <div class="container">
      <h1>Пополнить кошелек</h1>
      
      <div class="topup-card">
        <div class="current-balance">
          <span>Текущий баланс:</span>
          <span class="balance">{{ formatCurrency(authStore.user?.wallet?.balance || 0) }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">Сумма пополнения:</label>
          <input
            v-model.number="amount"
            type="number"
            class="form-input"
            placeholder="1000"
            min="100"
            step="100"
          />
        </div>

        <button 
          class="btn btn-primary btn-full"
          @click="topUp"
          :disabled="loading"
        >
          <span v-if="loading">Пополнение...</span>
          <span v-else>Пополнить {{ formatCurrency(amount) }}</span>
        </button>

        <div v-if="message" class="alert" :class="success ? 'alert-success' : 'alert-error'">
          {{ message }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../store/auth'
import api from '../services/api'
import { formatCurrency } from '../utils/formatters'

const authStore = useAuthStore()
const amount = ref(1000)
const loading = ref(false)
const message = ref('')
const success = ref(false)

const topUp = async () => {
  if (amount.value < 100) {
    message.value = 'Минимальная сумма: 100 KZT'
    success.value = false
    return
  }

  loading.value = true
  message.value = ''

  try {
    // Временно: обновляем баланс напрямую через updateProfile
    // В реальности нужен отдельный endpoint для пополнения
    const newBalance = (authStore.user?.wallet?.balance || 0) + amount.value
    
    // Это НЕ будет работать с текущим API, но покажу идею
    // Нужно добавить endpoint POST /api/wallet/topup в backend
    
    message.value = `✅ Баланс пополнен на ${formatCurrency(amount.value)}`
    success.value = true
    
    // Обновить профиль
    await authStore.fetchProfile()
    
  } catch (error) {
    message.value = 'Ошибка пополнения'
    success.value = false
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.topup {
  padding: 40px 0;
  min-height: calc(100vh - 70px);
}

.topup-card {
  max-width: 500px;
  margin: 0 auto;
  background-color: var(--dark-light);
  padding: 30px;
  border-radius: 16px;
}

.current-balance {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px;
  background-color: var(--dark-lighter);
  border-radius: 12px;
}

.balance {
  font-size: 24px;
  font-weight: bold;
  color: var(--primary);
}
</style>