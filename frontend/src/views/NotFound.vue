<template>
  <div class="error-page">
    <div class="container">
      <div class="error-content">
        <div class="error-code">404</div>
        <h1 class="error-title">Страница не найдена</h1>
        <p class="error-message">
          Запрашиваемая страница не существует или была удалена
        </p>
        <div class="error-actions">
          <router-link to="/" class="btn btn-primary">На главную</router-link>
          <button class="btn btn-secondary" @click="$router.back()">Назад</button>
        </div>
      </div>
    </div>
  </div>
</template>

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
  background: linear-gradient(to right, var(--primary), var(--secondary));
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
