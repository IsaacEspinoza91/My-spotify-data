<script setup>
import { ref, watch, onMounted } from 'vue'
import api from '../services/api'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import Pagination from '../components/ui/Pagination.vue'
import { Play, Music } from 'lucide-vue-next'

import FiltersBar from '../components/ui/FiltersBar.vue'

const loading = ref(true)
const songs = ref([])
const page = ref(1)
const totalPages = ref(1)
const total = ref(0)
const limit = 20

const filters = ref({
  search: '',
  start_date: '',
  end_date: '',
  start_hour: '',
  end_hour: ''
})

const fetchTopSongs = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      limit,
      ...(filters.value.search && { search: filters.value.search }),
      ...(filters.value.start_date && { start_date: filters.value.start_date }),
      ...(filters.value.end_date && { end_date: filters.value.end_date }),
      ...(filters.value.start_hour !== '' && { start_hour: filters.value.start_hour }),
      ...(filters.value.end_hour !== '' && { end_hour: filters.value.end_hour })
    }

    const response = await api.getTopSongs(params)
    songs.value = response.data.data || []
    totalPages.value = response.data.total_pages || 1
    total.value = response.data.total || 0
  } catch (error) {
    console.error("Error fetching top songs", error)
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  page.value = 1
  fetchTopSongs()
}

onMounted(() => {
  fetchTopSongs()
})

watch(page, () => {
  fetchTopSongs()
})

const formatNumber = (num) => new Intl.NumberFormat('es-CL').format(num || 0)
</script>

<template>
  <div class="top-songs-view view-container">
    <div class="view-hero-bg"></div>
    
    <div class="view-content-wrapper relative-z">
      <header class="view-header">
        <h1>Top Canciones</h1>
        <p class="text-gray subtitle">Los temas indiscutibles que te acompañan cada día.</p>
      </header>

    <FiltersBar 
      v-model="filters" 
      @filter-change="handleFilterChange" 
    />

    <LoadingSpinner v-if="loading" />

    <div v-else class="content-wrapper">
      <div v-if="songs.length === 0" class="empty-state">
        <p>No se encontraron canciones que coincidan con los filtros aplicados.</p>
      </div>
      <div v-else class="songs-table-container">
        <table class="songs-table">
          <thead>
            <tr>
              <th class="col-rank">#</th>
              <th class="col-title">Título</th>
              <th class="col-plays">Reproducciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(song, index) in songs" :key="index" class="song-row">
              <td class="col-rank">
                <span class="rank-number">{{ song.ranking }}</span>
                <Play class="play-icon" :size="16" />
              </td>
              <td class="col-title">
                <div class="song-info-wrapper">
                  <img v-if="song.song_image" :src="song.song_image" alt="Cover" class="song-image" />
                  <div v-else class="song-initials"><Music :size="20"/></div>
                  <div class="song-info">
                    <div class="song-name">{{ song.track_name }}</div>
                    <div class="song-artist text-gray">{{ song.artist_name }}</div>
                  </div>
                </div>
              </td>

              <td class="col-plays text-gray">{{ formatNumber(song.times_played) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <Pagination 
        v-model:page="page" 
        :totalPages="totalPages" 
        :total="total" 
      />
    </div>
    </div>
  </div>
</template>

<style scoped>
.view-container {
  display: flex;
  flex-direction: column;
  position: relative;
  min-height: 100%;
}

.view-hero-bg {
  position: absolute;
  top: -24px;
  left: -24px;
  right: -24px;
  height: 200px;
  background: linear-gradient(180deg, rgba(40, 60, 100, 1) 0%, rgba(24, 24, 24, 1) 100%);
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
  z-index: 0;
  pointer-events: none;
}

.relative-z {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.view-header h1 {
  font-size: 32px;
  font-weight: 800;
  margin-top: 32px;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
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

.songs-table-container {
  width: 100%;
}

.songs-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.songs-table th {
  color: var(--spotify-text-gray);
  font-size: 13px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 8px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.song-row {
  transition: background-color 0.2s ease;
}

.song-row:hover {
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.song-row:hover .rank-number {
  display: none;
}

.song-row:hover .play-icon {
  display: block;
  color: var(--spotify-white);
}

.col-rank {
  width: 50px;
  text-align: right;
  padding: 12px 16px;
  font-size: 16px;
  color: var(--spotify-text-gray);
  position: relative;
}

.play-icon {
  display: none;
  position: absolute;
  top: 50%;
  right: 16px;
  transform: translateY(-50%);
}

.col-title {
  padding: 12px 16px;
}

.song-info-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.song-image, .song-initials {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  object-fit: cover;
  flex-shrink: 0;
}

.song-initials {
  background-color: #333;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #b3b3b3;
}

.song-info {
  display: flex;
  flex-direction: column;
}

.song-name {
  font-size: 16px;
  color: var(--spotify-white);
  font-weight: 500;
}

.song-artist {
  font-size: 14px;
}

.col-plays {
  width: 150px;
  padding: 12px 16px;
  font-variant-numeric: tabular-nums;
}
</style>
