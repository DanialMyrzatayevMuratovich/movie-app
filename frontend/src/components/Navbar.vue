<template>
  <nav class="navbar">
    <div class="container">
      <div class="navbar-content">
        <!-- Logo -->
        <router-link to="/" class="navbar-logo">
          <span class="logo-icon">🎬</span>
          <span class="logo-text">CinemaHub</span>
        </router-link>

        <!-- Navigation Links -->
        <div class="navbar-links">
          <router-link to="/" class="nav-link">Фильмы</router-link>
          
          <template v-if="authStore.isAuthenticated">
            <router-link to="/profile" class="nav-link">Мои брони</router-link>
            
            <!-- User Menu -->
            <div class="user-menu">
              <button class="user-button" @click="toggleDropdown">
                <span class="user-name">{{ authStore.user?.fullName }}</span>
                <span class="user-icon">👤</span>
              </button>

              <!-- Dropdown -->
              <div v-if="showDropdown" class="dropdown">
                <div class="dropdown-item">
                  <span class="wallet-icon">💳</span>
                  <span>{{ formatCurrency(authStore.user?.wallet?.balance || 0) }}</span>
                </div>
                <hr class="dropdown-divider" />
                <router-link to="/profile" class="dropdown-item" @click="showDropdown = false">
                  <span>📊</span>
                  <span>Профиль</span>
                </router-link>
                <button class="dropdown-item" @click="logout">
                  <span>🚪</span>
                  <span>Выйти</span>
                </button>
              </div>
            </div>
          </template>

          <template v-else>
            <router-link to="/login" class="btn btn-primary">Войти</router-link>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { formatCurrency } from '../utils/formatters'

const router = useRouter()
const authStore = useAuthStore()
const showDropdown = ref(false)

const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
}

const logout = () => {
  authStore.logout()
  showDropdown.value = false
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  background-color: var(--dark-light);
  padding: 16px 0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.navbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 24px;
  font-weight: bold;
  color: var(--text);
}

.logo-icon {
  font-size: 32px;
}

.logo-text {
  background: linear-gradient(to right, var(--primary), var(--secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.navbar-links {
  display: flex;
  align-items: center;
  gap: 24px;
}

.nav-link {
  color: var(--text);
  font-weight: 500;
  transition: color 0.3s ease;
  position: relative;
}

.nav-link:hover {
  color: var(--primary);
}

.nav-link.router-link-active::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  right: 0;
  height: 3px;
  background-color: var(--primary);
}

.user-menu {
  position: relative;
}

.user-button {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: var(--dark-lighter);
  padding: 8px 16px;
  border-radius: 20px;
  transition: all 0.3s ease;
}

.user-button:hover {
  background-color: var(--dark);
}

.user-name {
  color: var(--text);
  font-weight: 500;
}

.user-icon {
  font-size: 20px;
}

.dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background-color: var(--dark-light);
  border-radius: 12px;
  min-width: 200px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
  padding: 8px;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: var(--text);
  border-radius: 8px;
  transition: background-color 0.3s ease;
  width: 100%;
  text-align: left;
  background: none;
}

.dropdown-item:hover {
  background-color: var(--dark-lighter);
}

.dropdown-divider {
  border: none;
  border-top: 1px solid var(--dark-lighter);
  margin: 8px 0;
}

.wallet-icon {
  font-size: 18px;
}

@media (max-width: 768px) {
  .user-name {
    display: none;
  }
}
</style>