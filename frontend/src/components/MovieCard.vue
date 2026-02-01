<template>
  <div class="movie-card" @click="goToMovie">
    <div class="movie-poster">
      <img 
        v-if="movie.posterUrl" 
        :src="movie.posterUrl" 
        :alt="movie.title"
        @error="handleImageError"
      />
      <div v-else class="poster-placeholder">
        <span class="poster-icon">🎬</span>
        <div class="poster-title">{{ movie.title }}</div>
      </div>

      <!-- Rating Badge -->
      <div class="rating-badge">
        <span class="star">⭐</span>
        <span>{{ movie.imdbRating }}</span>
      </div>

      <!-- Format Badge -->
      <div class="format-badges">
        <span v-for="format in formats" :key="format" class="badge badge-format">
          {{ format }}
        </span>
      </div>
    </div>

    <div class="movie-info">
      <h3 class="movie-title">{{ movie.title }}</h3>
      <p class="movie-subtitle text-gray">{{ movie.titleRu }}</p>
      
      <div class="movie-meta">
        <span class="meta-item">
          <span class="icon">🎭</span>
          {{ movie.genres?.join(', ') }}
        </span>
        <span class="meta-item">
          <span class="icon">⏱️</span>
          {{ formatDuration(movie.duration) }}
        </span>
        <span class="meta-item">
          <span class="icon">🔞</span>
          {{ movie.ageRestriction }}+
        </span>
      </div>

      <button class="btn btn-primary btn-full" @click.stop="goToMovie">
        Купить билет
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { formatDuration } from '../utils/formatters'

const props = defineProps({
  movie: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const formats = computed(() => {
  // Примерные форматы (можно получить из showtimes)
  return ['2D', '3D', 'IMAX'].slice(0, 2)
})

const goToMovie = () => {
  if (props.movie && props.movie._id) {
    router.push(`/movie/${props.movie._id}`)
  } else {
    console.error('Movie ID is missing:', props.movie)
  }
}

const handleImageError = (e) => {
  // Скрыть битое изображение и показать placeholder
  e.target.style.display = 'none'
  const placeholder = e.target.nextElementSibling
  if (placeholder) {
    placeholder.style.display = 'flex'
  }
}
</script>

<style scoped>
.movie-card {
  background-color: var(--dark-light);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.movie-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 30px rgba(229, 9, 20, 0.3);
}

.movie-poster {
  position: relative;
  width: 100%;
  aspect-ratio: 2/3;
  overflow: hidden;
  background-color: var(--dark-lighter);
}

.movie-poster img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.movie-card:hover .movie-poster img {
  transform: scale(1.05);
}

.poster-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--dark-lighter) 0%, var(--dark) 100%);
  padding: 20px;
  text-align: center;
}

.poster-icon {
  font-size: 64px;
  opacity: 0.5;
  margin-bottom: 12px;
}

.poster-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--text);
  opacity: 0.7;
  line-height: 1.3;
}

.rating-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background-color: rgba(0, 0, 0, 0.8);
  padding: 6px 12px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: bold;
  backdrop-filter: blur(10px);
  z-index: 2;
}

.star {
  font-size: 16px;
}

.format-badges {
  position: absolute;
  bottom: 12px;
  left: 12px;
  display: flex;
  gap: 8px;
  z-index: 2;
}

.badge-format {
  background-color: rgba(229, 9, 20, 0.9);
  color: white;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
}

.movie-info {
  padding: 16px;
}

.movie-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.movie-subtitle {
  font-size: 14px;
  margin-bottom: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.movie-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
  font-size: 13px;
  color: var(--text-gray);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.icon {
  font-size: 14px;
}

.btn-full {
  width: 100%;
}
</style>