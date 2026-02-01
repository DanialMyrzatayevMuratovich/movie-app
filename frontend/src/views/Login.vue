<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <h1 class="login-title">
            <span class="logo-icon">🎬</span>
            CinemaHub
          </h1>
          <p class="login-subtitle">
            {{ isLogin ? 'Войдите в свой аккаунт' : 'Создайте новый аккаунт' }}
          </p>
        </div>

        <!-- Error Alert -->
        <div v-if="error" class="alert alert-error">
          {{ error }}
        </div>

        <!-- Login Form -->
        <form v-if="isLogin" @submit.prevent="handleLogin">
          <div class="form-group">
            <label class="form-label">Email</label>
            <input
              v-model="loginForm.email"
              type="email"
              class="form-input"
              placeholder="your@email.com"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">Пароль</label>
            <input
              v-model="loginForm.password"
              type="password"
              class="form-input"
              placeholder="••••••••"
              required
            />
          </div>

          <button
            type="submit"
            class="btn btn-primary btn-full"
            :disabled="loading"
          >
            <span v-if="loading">Вход...</span>
            <span v-else>Войти</span>
          </button>
        </form>

        <!-- Register Form -->
        <form v-else @submit.prevent="handleRegister">
          <div class="form-group">
            <label class="form-label">Полное имя</label>
            <input
              v-model="registerForm.fullName"
              type="text"
              class="form-input"
              placeholder="Иван Иванов"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">Email</label>
            <input
              v-model="registerForm.email"
              type="email"
              class="form-input"
              placeholder="your@email.com"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">Телефон</label>
            <input
              v-model="registerForm.phone"
              type="tel"
              class="form-input"
              placeholder="+77771234567"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Пароль</label>
            <input
              v-model="registerForm.password"
              type="password"
              class="form-input"
              placeholder="••••••••"
              required
              minlength="6"
            />
          </div>

          <button
            type="submit"
            class="btn btn-primary btn-full"
            :disabled="loading"
          >
            <span v-if="loading">Регистрация...</span>
            <span v-else>Зарегистрироваться</span>
          </button>
        </form>

        <!-- Toggle Form -->
        <div class="form-toggle">
          <p v-if="isLogin">
            Нет аккаунта?
            <button class="toggle-btn" @click="isLogin = false">
              Зарегистрироваться
            </button>
          </p>
          <p v-else>
            Уже есть аккаунт?
            <button class="toggle-btn" @click="isLogin = true">
              Войти
            </button>
          </p>
        </div>

        <!-- Demo Accounts -->
        <div class="demo-accounts">
          <p class="demo-title">Тестовые аккаунты:</p>
          <div class="demo-buttons">
            <button
              class="btn btn-secondary btn-sm"
              @click="fillDemoCredentials('user')"
            >
              👤 Пользователь
            </button>
            <button
              class="btn btn-secondary btn-sm"
              @click="fillDemoCredentials('admin')"
            >
              👑 Админ
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const router = useRouter()
const authStore = useAuthStore()

const isLogin = ref(true)
const loading = ref(false)
const error = ref('')

const loginForm = ref({
  email: '',
  password: ''
})

const registerForm = ref({
  email: '',
  password: '',
  fullName: '',
  phone: ''
})

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    const result = await authStore.login(loginForm.value)

    if (result.success) {
      router.push('/')
    } else {
      error.value = result.error
    }
  } catch (err) {
    error.value = 'Ошибка входа. Попробуйте еще раз.'
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  loading.value = true
  error.value = ''

  try {
    const result = await authStore.register(registerForm.value)

    if (result.success) {
      router.push('/')
    } else {
      error.value = result.error
    }
  } catch (err) {
    error.value = 'Ошибка регистрации. Попробуйте еще раз.'
  } finally {
    loading.value = false
  }
}

const fillDemoCredentials = (type) => {
  isLogin.value = true

  if (type === 'user') {
    loginForm.value.email = 'user1@example.com'
    loginForm.value.password = 'user123'
  } else if (type === 'admin') {
    loginForm.value.email = 'admin@cinema.kz'
    loginForm.value.password = 'admin123'
  }
}
</script>

<style scoped>
.login-page {
  min-height: calc(100vh - 70px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, rgba(229, 9, 20, 0.1), rgba(255, 193, 7, 0.1));
}

.login-container {
  width: 100%;
  max-width: 480px;
}

.login-card {
  background-color: var(--dark-light);
  padding: 40px;
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.logo-icon {
  font-size: 40px;
}

.login-subtitle {
  color: var(--text-gray);
  font-size: 16px;
}

.btn-full {
  width: 100%;
  margin-top: 8px;
}

.form-toggle {
  text-align: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--dark-lighter);
}

.form-toggle p {
  color: var(--text-gray);
}

.toggle-btn {
  color: var(--primary);
  background: none;
  border: none;
  font-weight: 600;
  margin-left: 4px;
  transition: color 0.3s ease;
}

.toggle-btn:hover {
  color: var(--secondary);
  text-decoration: underline;
}

.demo-accounts {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--dark-lighter);
  text-align: center;
}

.demo-title {
  color: var(--text-gray);
  font-size: 14px;
  margin-bottom: 12px;
}

.demo-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}

@media (max-width: 480px) {
  .login-card {
    padding: 30px 20px;
  }

  .login-title {
    font-size: 28px;
  }
}
</style>