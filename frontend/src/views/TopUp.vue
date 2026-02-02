<template>
  <div class="topup">
    <div class="container">
      <div class="topup-header">
        <button class="btn btn-secondary" @click="goBack">
          ← Назад
        </button>
        <h1 class="page-title">Пополнение кошелька</h1>
      </div>

      <div class="topup-card">
        <div class="current-balance">
          <div class="balance-icon">💳</div>
          <div class="balance-info">
            <div class="balance-label">Текущий баланс</div>
            <div class="balance-value">{{ formatCurrency(authStore.user?.wallet?.balance || 0) }}</div>
          </div>
        </div>

        <div class="topup-form">
          <h3>Введите сумму пополнения</h3>

          <!-- Quick amounts -->
          <div class="quick-amounts">
            <button
              v-for="qa in quickAmounts"
              :key="qa"
              class="quick-btn"
              :class="{ active: amount === qa }"
              @click="amount = qa"
            >
              {{ formatCurrency(qa) }}
            </button>
          </div>

          <!-- Custom amount -->
          <div class="input-group">
            <label for="amount">Сумма (KZT)</label>
            <input
              id="amount"
              v-model.number="amount"
              type="number"
              min="100"
              max="1000000"
              step="100"
              placeholder="Введите сумму"
              class="input-field"
            />
          </div>

          <div v-if="error" class="alert alert-error">
            {{ error }}
          </div>

          <div v-if="successMsg" class="alert alert-success">
            {{ successMsg }}
          </div>

          <button
            class="btn btn-primary btn-large"
            :disabled="loading || !amount || amount <= 0"
            @click="handleTopUp"
          >
            <span v-if="loading">Пополнение...</span>
            <span v-else>Пополнить {{ amount > 0 ? formatCurrency(amount) : '' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import api from '../services/api'
import { formatCurrency } from '../utils/formatters'

const router = useRouter()
const authStore = useAuthStore()

const amount = ref(5000)
const loading = ref(false)
const error = ref('')
const successMsg = ref('')

const quickAmounts = [1000, 2000, 5000, 10000, 20000, 50000]

const handleTopUp = async () => {
  if (!amount.value || amount.value <= 0) {
    error.value = 'Введите корректную сумму'
    return
  }

  if (amount.value < 100) {
    error.value = 'Минимальная сумма пополнения: 100 KZT'
    return
  }

  if (amount.value > 1000000) {
    error.value = 'Максимальная сумма пополнения: 1,000,000 KZT'
    return
  }

  loading.value = true
  error.value = ''
  successMsg.value = ''

  try {
    await api.topUpWallet({ amount: amount.value })

    // Обновить профиль (баланс изменился)
    await authStore.fetchProfile()

    successMsg.value = `Кошелек успешно пополнен на ${formatCurrency(amount.value)}`
  } catch (err) {
    error.value = err.response?.data?.error || 'Не удалось пополнить кошелек'
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.go(-1)
}
</script>

<style scoped>
.topup {
  padding: 40px 0;
  min-height: calc(100vh - 70px);
}

.topup-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.page-title {
  font-size: 32px;
  font-weight: bold;
}

.topup-card {
  max-width: 600px;
  margin: 0 auto;
  background: linear-gradient(135deg, var(--dark-light), var(--dark-lighter));
  padding: 30px;
  border-radius: 16px;
  border: 2px solid var(--dark-lighter);
}

.current-balance {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background-color: var(--dark);
  border-radius: 12px;
  border: 2px solid var(--primary);
  margin-bottom: 30px;
}

.balance-icon {
  font-size: 48px;
}

.balance-label {
  font-size: 14px;
  color: var(--text-gray);
  margin-bottom: 6px;
}

.balance-value {
  font-size: 32px;
  font-weight: bold;
  color: var(--primary);
}

.topup-form h3 {
  font-size: 20px;
  margin-bottom: 16px;
}

.quick-amounts {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.quick-btn {
  padding: 14px;
  background-color: var(--dark);
  color: var(--text);
  border: 2px solid var(--dark-lighter);
  border-radius: 12px;
  font-weight: 600;
  font-size: 15px;
  transition: all 0.3s ease;
}

.quick-btn:hover {
  border-color: var(--primary);
}

.quick-btn.active {
  background-color: var(--primary);
  border-color: var(--primary);
  color: white;
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  font-size: 14px;
  color: var(--text-gray);
  margin-bottom: 8px;
}

.input-field {
  width: 100%;
  padding: 14px 16px;
  background-color: var(--dark);
  color: var(--text);
  border: 2px solid var(--dark-lighter);
  border-radius: 12px;
  font-size: 18px;
  font-weight: 600;
  transition: border-color 0.3s ease;
}

.input-field:focus {
  outline: none;
  border-color: var(--primary);
}

.input-field::placeholder {
  color: var(--text-gray);
  font-weight: 400;
}

.btn-large {
  width: 100%;
  padding: 16px 32px;
  font-size: 18px;
  font-weight: bold;
  margin-top: 10px;
}

.alert {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 14px;
}

.alert-error {
  background-color: rgba(229, 62, 62, 0.15);
  color: #fc8181;
  border: 1px solid rgba(229, 62, 62, 0.3);
}

.alert-success {
  background-color: rgba(72, 187, 120, 0.15);
  color: #68d391;
  border: 1px solid rgba(72, 187, 120, 0.3);
}

@media (max-width: 768px) {
  .quick-amounts {
    grid-template-columns: repeat(2, 1fr);
  }

  .current-balance {
    flex-direction: column;
    text-align: center;
  }
}
</style>
