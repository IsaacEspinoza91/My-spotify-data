<script setup>
import { ref, watch, onMounted } from 'vue'
import api from '../services/api'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import Pagination from '../components/ui/Pagination.vue'

import FiltersBar from '../components/ui/FiltersBar.vue'

const loading = ref(true)
const artists = ref([])
const page = ref(1)
const totalPages = ref(1)
const total = ref(0)
const limit = ref(15) // Default to 15 (5 cols * 3 rows)

const filters = ref({
  search: '',
  start_date: '',
  end_date: '',
  start_hour: '',
  end_hour: ''
})

// Calculate limit based on available width and grid settings
// Grid columns minmax is 240px, gap is 24px.
const calculateDynamicLimit = () => {
  // Approximate sidebar width (240px) + padding
  const availableWidth = window.innerWidth - 300 
  if (availableWidth > 0) {
    // Column width is approx 240px + 24px gap = 264px
    const cols = Math.floor(availableWidth / 264)
    // We want 3 rows always, minimum 1 column
    const newLimit = Math.max(1, cols) * 3
    
    if (newLimit !== limit.value) {
      limit.value = newLimit
      page.value = 1 // Reset to first page when limit changes
      fetchTopArtists()
    }
  }
}

const fetchTopArtists = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      limit: limit.value,
      ...(filters.value.search && { search: filters.value.search }),
      ...(filters.value.start_date && { start_date: filters.value.start_date }),
      ...(filters.value.end_date && { end_date: filters.value.end_date }),
      ...(filters.value.start_hour !== '' && { start_hour: filters.value.start_hour }),
      ...(filters.value.end_hour !== '' && { end_hour: filters.value.end_hour })
    }
    
    const response = await api.getTopArtists(params)
    artists.value = response.data.data || []
    totalPages.value = response.data.total_pages || 1
    total.value = response.data.total || 0
  } catch (error) {
    console.error("Error fetching top artists", error)
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  page.value = 1
  fetchTopArtists()
}

let resizeTimeout
onMounted(() => {
  calculateDynamicLimit()
  fetchTopArtists()
  
  window.addEventListener('resize', () => {
    clearTimeout(resizeTimeout)
    resizeTimeout = setTimeout(calculateDynamicLimit, 200)
  })
})

watch(page, () => {
  fetchTopArtists()
})

const formatNumber = (num) => new Intl.NumberFormat('es-CL').format(num || 0)

const getInitials = (name) => {
  if (!name) return '?'
  const parts = name.split(' ').filter(p => p.length > 0)
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase()
  return (parts[0][0] + parts[1][0]).toUpperCase()
}
</script>

<template>
  <div class="top-artists-view view-container">
    <header class="view-header">
      <h1>Top Artistas</h1>
    </header>

    <FiltersBar 
      v-model="filters" 
      @filter-change="handleFilterChange" 
    />

    <LoadingSpinner v-if="loading" />

    <div v-else class="content-wrapper">
      <div v-if="artists.length === 0" class="empty-state">
        <p>No se encontraron artistas que coincidan con los filtros aplicados.</p>
      </div>
      <div v-else class="artists-grid">
        <div 
          v-for="artist in artists" 
          :key="artist.artist_name" 
          class="artist-card card"
        >
          <div class="artist-image-container">
            <img v-if="artist.artist_image" :src="artist.artist_image" alt="Artist" class="artist-image" />
            <div v-else class="artist-initials">{{ getInitials(artist.artist_name) }}</div>
            <div class="ranking-badge">{{ artist.ranking }}</div>
          </div>
          <div class="artist-info">
            <h3 class="artist-name">{{ artist.artist_name }}</h3>
            <div class="artist-stats text-gray">
              <span>{{ formatNumber(artist.times_played) }} reproducciones</span>
              <span>•</span>
              <span>{{ formatNumber(Math.round(artist.minutes_played / 60)) }} horas</span>
            </div>
          </div>
        </div>
      </div>
      
      <Pagination 
        v-model:page="page" 
        :totalPages="totalPages" 
        :total="total" 
      />
    </div>
  </div>
</template>

<style scoped>
.view-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.view-header h1 {
  font-size: 32px;
  font-weight: 800;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px;
  background-color: var(--spotify-dark-gray);
  border-radius: 8px;
  color: var(--spotify-text-gray);
  font-size: 16px;
  text-align: center;
}

.artists-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 24px;
  margin-bottom: 24px;
}

.artist-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px;
  transition: background-color 0.3s ease;
}

.artist-card:hover {
  background-color: var(--spotify-light-gray);
}

.ranking-badge {
  position: absolute;
  top: 0;
  right: 0;
  background-color: var(--spotify-green);
  color: var(--spotify-black);
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 14px;
  transform: translate(30%, -30%);
}

.artist-image-container {
  position: relative;
  width: 100%;
  margin-bottom: 8px;
}

.artist-image {
  width: 100%;
  aspect-ratio: 1 / 1;
  border-radius: 50%;
  object-fit: cover;
  display: block;
  background-color: var(--spotify-dark-gray);
  box-shadow: 0 8px 24px rgba(0,0,0,0.5);
}

.artist-initials {
  width: 100%;
  aspect-ratio: 1 / 1;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 700;
  color: var(--spotify-white);
  background-color: #333333;
  box-shadow: 0 8px 24px rgba(0,0,0,0.5);
}

.artist-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.artist-name {
  font-size: 18px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.artist-stats {
  font-size: 13px;
  display: flex;
  gap: 8px;
  align-items: center;
}
</style>
