<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../services/api'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import { Play, Music } from 'lucide-vue-next'

const loading = ref(false)
const wrappedData = ref([])
const errorMsg = ref('')

const year = ref(new Date().getFullYear().toString())
const filterType = ref('year') // 'year', 'month', 'season'
const month = ref('1')
const season = ref('spring')

const months = [
  { value: '1', label: 'Enero' }, { value: '2', label: 'Febrero' },
  { value: '3', label: 'Marzo' }, { value: '4', label: 'Abril' },
  { value: '5', label: 'Mayo' }, { value: '6', label: 'Junio' },
  { value: '7', label: 'Julio' }, { value: '8', label: 'Agosto' },
  { value: '9', label: 'Septiembre' }, { value: '10', label: 'Octubre' },
  { value: '11', label: 'Noviembre' }, { value: '12', label: 'Diciembre' }
]

const seasons = [
  { value: 'spring', label: 'Primavera' },
  { value: 'summer', label: 'Verano' },
  { value: 'autumn', label: 'Otoño' },
  { value: 'winter', label: 'Invierno' }
]

const fetchWrapped = async () => {
  loading.value = true
  errorMsg.value = ''
  wrappedData.value = []
  
  const params = { year: year.value }
  if (filterType.value === 'month') {
    params.month = month.value
  } else if (filterType.value === 'season') {
    params.season = season.value
  }

  try {
    const response = await api.getWrapped(params)
    const data = response.data
    
    if (!data || data.length === 0) {
      errorMsg.value = "No hay datos para el periodo seleccionado."
      return
    }
    wrappedData.value = data
  } catch (error) {
    console.error("Error fetching wrapped data", error)
    errorMsg.value = "Hubo un error al buscar el Wrapped."
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchWrapped()
})

const formatNumber = (num) => new Intl.NumberFormat('es-CL').format(num || 0)
</script>

<template>
  <div class="wrapped-view view-container">
    <header class="view-header">
      <div class="header-content">
        <h1>Tu Wrapped</h1>
        <p class="text-gray subtitle">Descubre tus canciones más escuchadas por periodo</p>
      </div>
      
      <div class="filters-container">
        <div class="filter-group">
          <label>Periodo</label>
          <select v-model="filterType" class="spotify-select">
            <option value="year">Anual</option>
            <option value="month">Mensual</option>
            <option value="season">Por Estación</option>
          </select>
        </div>

        <div class="filter-group">
          <label>Año</label>
          <input type="number" v-model="year" class="spotify-input" min="2010" max="2100" />
        </div>

        <div v-if="filterType === 'month'" class="filter-group">
          <label>Mes</label>
          <select v-model="month" class="spotify-select">
            <option v-for="m in months" :key="m.value" :value="m.value">{{ m.label }}</option>
          </select>
        </div>

        <div v-if="filterType === 'season'" class="filter-group">
          <label>Estación</label>
          <select v-model="season" class="spotify-select">
            <option v-for="s in seasons" :key="s.value" :value="s.value">{{ s.label }}</option>
          </select>
        </div>

        <button class="generate-btn" @click="fetchWrapped" :disabled="loading">
          Generar
        </button>
      </div>
    </header>

    <LoadingSpinner v-if="loading" />

    <div v-else-if="errorMsg" class="error-message card">
      {{ errorMsg }}
    </div>

    <div v-else-if="wrappedData.length > 0" class="content-wrapper">
      <div class="wrapped-banner gradient-bg">
        <h2>Top Tracks</h2>
        <div class="hero-content">
          <img v-if="wrappedData[0].song_image" :src="wrappedData[0].song_image" class="hero-image" />
          <div class="hero-image-fallback" v-else>
            <Music :size="48" />
          </div>
          <div class="hero-text">
            <div class="big-number">#1 {{ wrappedData[0].track_name }}</div>
            <p>de {{ wrappedData[0].artist_name }} ({{ wrappedData[0].times_played }} reproducciones)</p>
          </div>
        </div>
      </div>

      <div class="songs-table-container mt-4">
        <table class="songs-table">
          <thead>
            <tr>
              <th class="col-rank">#</th>
              <th class="col-title">Título</th>
              <th class="col-plays">Reproducciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(song, index) in wrappedData" :key="index" class="song-row">
              <td class="col-rank">
                <span class="rank-number">{{ song.ranking }}</span>
                <Play class="play-icon" :size="16" />
              </td>
              <td class="col-title">
                <div class="song-info-wrapper">
                  <img v-if="song.song_image" :src="song.song_image" class="table-song-image" />
                  <div class="table-song-fallback" v-else>
                    <Music :size="16" />
                  </div>
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
    </div>
  </div>
</template>

<style scoped>
.view-container {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.view-header {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.view-header h1 {
  font-size: 40px;
  font-weight: 900;
  letter-spacing: -1px;
}

.filters-container {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: 16px;
  background-color: var(--spotify-dark-gray);
  padding: 16px;
  border-radius: 8px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--spotify-text-gray);
  letter-spacing: 1px;
}

.spotify-select, .spotify-input {
  background-color: var(--spotify-light-gray);
  color: var(--spotify-white);
  border: 1px solid transparent;
  padding: 0 16px;
  border-radius: 4px;
  font-size: 14px;
  min-width: 140px;
  height: 40px;
  box-sizing: border-box;
}

.spotify-select:focus, .spotify-input:focus {
  outline: none;
  border-color: #727272;
}

.generate-btn {
  height: 40px;
  margin-left: auto;
}

.wrapped-banner {
  padding: 40px;
  border-radius: 8px;
  margin-bottom: 24px;
  color: white;
}

.gradient-bg {
  background: linear-gradient(135deg, #1ED760 0%, #006450 100%);
}

.wrapped-banner h2 {
  font-size: 18px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 2px;
  margin-bottom: 24px;
}

.hero-content {
  display: flex;
  align-items: center;
  gap: 32px;
}

.hero-image {
  width: 140px;
  height: 140px;
  border-radius: 8px;
  object-fit: cover;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
}

.hero-image-fallback {
  width: 140px;
  height: 140px;
  border-radius: 8px;
  background-color: var(--spotify-dark-gray);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--spotify-text-gray);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
}

.hero-text {
  display: flex;
  flex-direction: column;
}

.big-number {
  font-size: 48px;
  font-weight: 900;
  line-height: 1.1;
  margin-bottom: 8px;
  letter-spacing: -1px;
}

.error-message {
  padding: 32px;
  text-align: center;
  color: var(--spotify-text-gray);
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
  gap: 16px;
}

.table-song-image {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  object-fit: cover;
}

.table-song-fallback {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  background-color: var(--spotify-light-gray);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--spotify-text-gray);
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

.mt-4 {
  margin-top: 16px;
}
</style>
