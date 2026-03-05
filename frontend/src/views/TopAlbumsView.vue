<script setup>
import { ref, watch, onMounted } from 'vue'
import api from '../services/api'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import Pagination from '../components/ui/Pagination.vue'
import { Disc } from 'lucide-vue-next'

import FiltersBar from '../components/ui/FiltersBar.vue'

const loading = ref(true)
const albums = ref([])
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
// Grid columns minmax is 200px, gap is 24px.
const calculateDynamicLimit = () => {
  const availableWidth = window.innerWidth - 300 
  if (availableWidth > 0) {
    // Column width is approx 200px + 24px gap = 224px
    const cols = Math.floor(availableWidth / 224)
    // We want 3 rows always, minimum 1 column
    const newLimit = Math.max(1, cols) * 3
    
    if (newLimit !== limit.value) {
      limit.value = newLimit
      page.value = 1
      fetchTopAlbums()
    }
  }
}

const fetchTopAlbums = async () => {
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

    const response = await api.getTopAlbums(params)
    albums.value = response.data.data || []
    totalPages.value = response.data.total_pages || 1
    total.value = response.data.total || 0
  } catch (error) {
    console.error("Error fetching top albums", error)
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  page.value = 1
  fetchTopAlbums()
}

let resizeTimeout
onMounted(() => {
  calculateDynamicLimit()
  fetchTopAlbums()
  
  window.addEventListener('resize', () => {
    clearTimeout(resizeTimeout)
    resizeTimeout = setTimeout(calculateDynamicLimit, 200)
  })
})

watch(page, () => {
  fetchTopAlbums()
})

const formatNumber = (num) => new Intl.NumberFormat('es-CL').format(num || 0)
</script>

<template>
  <div class="top-albums-view view-container">
    <header class="view-header">
      <h1>Top Álbumes</h1>
    </header>

    <FiltersBar 
      v-model="filters" 
      @filter-change="handleFilterChange" 
    />

    <LoadingSpinner v-if="loading" />

    <div v-else class="content-wrapper">
      <div v-if="albums.length === 0" class="empty-state">
        <p>No se encontraron álbumes que coincidan con los filtros aplicados.</p>
      </div>
      <div v-else class="albums-grid">
        <div 
          v-for="(album, index) in albums" 
          :key="index" 
          class="album-card card"
        >
          <div class="album-image-container">
            <img v-if="album.album_image" :src="album.album_image" alt="Album Cover" class="album-image" />
            <div v-else class="album-initials"><Disc :size="48"/></div>
            <div class="ranking-badge">{{ album.ranking }}</div>
          </div>
          
          <div class="album-info-container">
            <h3 class="album-title">{{ album.album_name }}</h3>
            <div class="album-artist text-gray">{{ album.artist_name }}</div>
            <div class="album-stats text-green">
              {{ formatNumber(album.times_played) }} reproducciones
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

.albums-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 40px 24px;
  margin-bottom: 24px;
}

.album-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
  background-color: var(--spotify-dark-gray);
  transition: background-color 0.3s ease;
  height: 100%;
}

.album-card:hover {
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

.album-image-container {
  position: relative;
  width: 100%;
}

.album-image {
  width: 100%;
  aspect-ratio: 1 / 1;
  border-radius: 4px;
  object-fit: cover;
  display: block;
  background-color: var(--spotify-dark-gray);
  box-shadow: 0 8px 24px rgba(0,0,0,0.5);
}

.album-initials {
  width: 100%;
  aspect-ratio: 1 / 1;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #333333;
  background-color: #1a1a1a;
  box-shadow: 0 8px 24px rgba(0,0,0,0.5);
}

.album-info-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-grow: 1;
  justify-content: flex-end;
  margin-top: 24px;
}

.album-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--spotify-white);
  line-height: 1.2;
}

.album-artist {
  font-size: 14px;
  font-weight: 500;
}

.album-stats {
  font-size: 13px;
  font-weight: 600;
  margin-top: auto;
}
</style>
