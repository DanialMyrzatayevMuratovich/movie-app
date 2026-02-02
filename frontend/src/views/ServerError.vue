<template>
  <div class="error-page">
    <div class="container">
      <div class="error-content">
        <div class="error-code">{{ code }}</div>
        <h1 class="error-title">{{ title }}</h1>
        <p class="error-message">{{ message }}</p>
        <div class="error-actions">
          <router-link to="/" class="btn btn-primary">На главную</router-link>
          <button class="btn btn-secondary" @click="retry">Попробовать снова</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const code = computed(() => route.query.code || '500')

const title = computed(() => {
  const titles = {
    '500': 'Ошибка сервера',
    '502': 'Сервер недоступен',
    '503': 'Сервис временно недоступен',
    '403': 'Доступ запрещён'
  }
  return titles[code.value] || 'Произошла ошибка'
})

const message = computed(() => {
  const messages = {
    '500': 'Произошла внутренняя ошибка сервера. Попробуйте позже.',
    '502': 'Сервер временно не отвечает. Попробуйте обновить страницу.',
    '503': 'Сервис временно недоступен. Ведутся технические работы.',
    '403': 'У вас нет прав для доступа к этой странице.'
  }
  return messages[code.value] || 'Что-то пошло не так. Попробуйте позже.'
})

const retry = () => {
  const returnTo = route.query.from
  if (returnTo) {
    router.push(returnTo)
  } else {
    router.push('/')
  }
}
</script>

<style scoped>
.error-page {
  min-height: calc(100vh - 70px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.error-content {
  text-align: center;
  max-width: 500px;
  margin: 0 auto;
}

.error-code {
  font-size: 120px;
  font-weight: 900;
  background: linear-gradient(to right, var(--danger, #e50914), var(--secondary, #ffc107));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1;
  margin-bottom: 16px;
}

.error-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 12px;
}

.error-message {
  font-size: 16px;
  color: var(--text-gray);
  margin-bottom: 32px;
  line-height: 1.5;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

@media (max-width: 768px) {
  .error-code {
    font-size: 80px;
  }

  .error-title {
    font-size: 22px;
  }
}
</style>
