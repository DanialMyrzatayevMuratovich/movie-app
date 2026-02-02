<template>
  <div class="home">
    <div class="container">
      <!-- Hero Section -->
      <div class="hero">
        <h1 class="hero-title">
          <span class="gradient-text">Добро пожаловать</span>
          <br />
          в CinemaHub
        </h1>
        <p class="hero-subtitle">
          Бронируйте билеты на лучшие фильмы онлайн
        </p>
      </div>

      <!-- Filters -->
      <div class="filters">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="🔍 Поиск фильмов..."
          class="search-input"
          @input="searchMovies"
        />

        <div class="filter-buttons">
          <button
            v-for="genre in genres"
            :key="genre"
            class="filter-btn"
            :class="{ active: selectedGenre === genre }"
            @click="selectGenre(genre)"
          >
            {{ genre }}
          </button>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
      </div>

      <!-- Movies Grid -->
      <div v-else-if="movies.length > 0" class="movies-grid">
        <MovieCard v-for="movie in movies" :key="movie._id || movie.id" :movie="movie" />
      </div>

      <!-- No Results -->
      <div v-else class="no-results">
        <span class="no-results-icon">🎬</span>
        <h3>Фильмы не найдены</h3>
        <p>Попробуйте изменить критерии поиска</p>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button
          class="btn btn-secondary"
          :disabled="currentPage === 1"
          @click="goToPage(currentPage - 1)"
        >
          ← Назад
        </button>

        <div class="page-numbers">
          <button
            v-for="page in visiblePages"
            :key="page"
            class="page-btn"
            :class="{ active: page === currentPage }"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
        </div>

        <button
          class="btn btn-secondary"
          :disabled="currentPage === totalPages"
          @click="goToPage(currentPage + 1)"
        >
          Далее →
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import MovieCard from '../components/MovieCard.vue'
import api from '../services/api'

const movies = ref([])
const loading = ref(true)
const searchQuery = ref('')
const selectedGenre = ref('Все')
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 12

const genres = ref([
  'Все',
  'Sci-Fi',
  'Action',
  'Drama',
  'Horror',
  'Animation',
  'Comedy',
  'Thriller'
])

const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)

  if (end - start < maxVisible - 1) {
    start = Math.max(1, end - maxVisible + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

const fetchMovies = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: limit,
      isActive: 'true'
    }

    if (searchQuery.value) {
      params.search = searchQuery.value
    }

    if (selectedGenre.value !== 'Все') {
      params.genre = selectedGenre.value
    }

    const response = await api.getMovies(params)
    movies.value = response.data.data
    totalPages.value = response.data.pagination.totalPages
  } catch (error) {
    console.error('Failed to fetch movies:', error)
    movies.value = []
  } finally {
    loading.value = false
  }
}

const searchMovies = () => {
  currentPage.value = 1
  fetchMovies()
}

const selectGenre = (genre) => {
  selectedGenre.value = genre
  currentPage.value = 1
  fetchMovies()
}

const goToPage = (page) => {
  currentPage.value = page
  fetchMovies()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  fetchMovies()
})
</script>

<style scoped>
.home {
  padding: 40px 0;
  min-height: calc(100vh - 70px);
}

.hero {
  text-align: center;
  margin-bottom: 60px;
  padding: 60px 20px;
  background: linear-gradient(135deg, rgba(229, 9, 20, 0.1), rgba(255, 193, 7, 0.1));
  border-radius: 20px;
}

.hero-title {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 16px;
  line-height: 1.2;
}

.gradient-text {
  background: linear-gradient(to right, var(--primary), var(--secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-subtitle {
  font-size: 20px;
  color: var(--text-gray);
}

.filters {
  margin-bottom: 40px;
}

.search-input {
  width: 100%;
  max-width: 600px;
  margin: 0 auto 24px;
  display: block;
  padding: 16px 24px;
  background-color: var(--dark-light);
  border: 2px solid var(--dark-lighter);
  border-radius: 30px;
  color: var(--text);
  font-size: 16px;
  transition: all 0.3s ease;
}

.search-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(229, 9, 20, 0.1);
}

.filter-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 10px 20px;
  background-color: var(--dark-light);
  color: var(--text);
  border: 2px solid var(--dark-lighter);
  border-radius: 20px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.filter-btn:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
}

.filter-btn.active {
  background-color: var(--primary);
  border-color: var(--primary);
  color: white;
}

.movies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 24px;
  margin-bottom: 40px;
}

.no-results {
  text-align: center;
  padding: 80px 20px;
}

.no-results-icon {
  font-size: 80px;
  display: block;
  margin-bottom: 20px;
  opacity: 0.5;
}

.no-results h3 {
  font-size: 24px;
  margin-bottom: 12px;
}

.no-results p {
  color: var(--text-gray);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 40px;
}

.page-numbers {
  display: flex;
  gap: 8px;
}

.page-btn {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background-color: var(--dark-light);
  color: var(--text);
  border: 2px solid var(--dark-lighter);
  font-weight: 600;
  transition: all 0.3s ease;
}

.page-btn:hover {
  border-color: var(--primary);
}

.page-btn.active {
  background-color: var(--primary);
  border-color: var(--primary);
  color: white;
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 32px;
  }

  .hero-subtitle {
    font-size: 16px;
  }

  .movies-grid {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 16px;
  }
}
</style>