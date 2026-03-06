<script setup>
import { ref } from 'vue'
import api from '../services/api'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import LineChart from '../components/ui/LineChart.vue'
import BarChart from '../components/ui/BarChart.vue'
import DoughnutChart from '../components/ui/DoughnutChart.vue'
import { Search } from 'lucide-vue-next'

const searchQuery = ref('')
const searchType = ref('track')
const loading = ref(false)
const evolutionData = ref(null)
const yearlyDataMinutes = ref(null)
const yearlyDataSongs = ref(null)
const errorMsg = ref('')

const searchArtist = async () => {
  loading.value = true
  errorMsg.value = ''
  evolutionData.value = null
  yearlyDataMinutes.value = null
  yearlyDataSongs.value = null
  
  try {
    const params = searchQuery.value.trim() 
      ? { [searchType.value]: searchQuery.value.trim() } 
      : {}
    const [evolutionRes, yearlyRes] = await Promise.all([
      api.getEvolution(params),
      api.getYearly(params)
    ])
    
    const data = evolutionRes.data
    const yearly = yearlyRes.data
    
    if (!data || data.length === 0) {
      errorMsg.value = "No se encontraron datos en tu historial para esta búsqueda."
      return
    }

    evolutionData.value = {
      labels: data.map(item => item.year_month),
      datasets: [
        {
          label: 'Horas Mensuales',
          backgroundColor: '#1ED760',
          borderColor: '#1ED760',
          data: data.map(item => item.hours_monthly)
        }
      ]
    }
    
    if (yearly && yearly.length > 0) {
      const colors = [
        '#1ED760', '#FF5A5F', '#00A699', '#FFB400', '#9C27B0', 
        '#03A9F4', '#E91E63', '#00E676', '#3F51B5', '#FF7043', 
        '#8BC34A', '#FFCA28'
      ]
      
      yearlyDataMinutes.value = {
        labels: yearly.map(item => item.year.toString()),
        datasets: [{
          label: 'Total Minutos',
          backgroundColor: colors,
          data: yearly.map(item => item.total_minutes),
          borderWidth: 0,
        }]
      }
      
      yearlyDataSongs.value = {
        labels: yearly.map(item => item.year.toString()),
        datasets: [{
          label: searchType.value === 'track' ? 'Total Veces Escuchada' : 'Total Canciones Únicas',
          backgroundColor: colors,
          data: yearly.map(item => item.total_songs),
          borderWidth: 0,
        }]
      }
    }
  } catch (error) {
    console.error("Error fetching evolution data", error)
    errorMsg.value = "Hubo un error al buscar los datos. Intenta nuevamente."
  } finally {
    loading.value = false
  }
}

let searchTimeout
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    searchArtist()
  }, 500)
}

import { onMounted } from 'vue'

onMounted(() => {
  searchArtist()
})
</script>

<template>
  <div class="evolution-view view-container">
    <header class="view-header">
      <h1>Evolución de Escucha</h1>
      <p class="text-gray subtitle">Busca una canción, álbum o artista para ver cómo ha evolucionado tu tiempo escuchándolo a lo largo de los meses.</p>
    </header>

    <div class="search-container">
      <div class="search-bar">
        <Search class="search-icon" :size="20" />
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Ej: Bad Bunny, Radiohead, Los Tres..."
          class="search-input"
          @input="debouncedSearch"
        />
      </div>
      
      <div class="type-filters">
        <button 
          class="type-btn" 
          :class="{ active: searchType === 'track' }"
          @click="searchType = 'track'; debouncedSearch()"
        >
          Canción
        </button>
        <button 
          class="type-btn" 
          :class="{ active: searchType === 'album' }"
          @click="searchType = 'album'; debouncedSearch()"
        >
          Álbum
        </button>
        <button 
          class="type-btn" 
          :class="{ active: searchType === 'artist' }"
          @click="searchType = 'artist'; debouncedSearch()"
        >
          Artista
        </button>
      </div>
    </div>

    <LoadingSpinner v-if="loading" />

    <div v-else-if="errorMsg" class="error-message card">
      {{ errorMsg }}
    </div>

    <div v-else-if="evolutionData" class="results-container">
      <div class="chart-wrapper card main-chart">
        <h2 v-if="searchQuery">
          Horas escuchadas de 
          <span v-if="searchType === 'track'">canción</span>
          <span v-else-if="searchType === 'album'">álbum</span>
          <span v-else>artista</span>
          "{{ searchQuery }}"
        </h2>
        <h2 v-else>Horas escuchadas generales</h2>
        <LineChart :chartData="evolutionData" />
      </div>

      <div class="yearly-charts-grid" v-if="yearlyDataMinutes && yearlyDataSongs">
        <div class="chart-wrapper card sub-chart">
          <h2 v-if="searchQuery">
            Total Minutos Anuales 
            <span v-if="searchType === 'track'">de la canción</span>
            <span v-else-if="searchType === 'album'">del álbum</span>
            <span v-else>del artista</span>
            "{{ searchQuery }}"
          </h2>
          <h2 v-else>Total Minutos Anuales Generales</h2>
          <BarChart :chartData="yearlyDataMinutes" />
        </div>
        
        <div class="chart-wrapper card sub-chart">
          <h2 v-if="searchQuery">
            <span v-if="searchType === 'track'">Total anual de veces escuchada la canción</span>
            <span v-else-if="searchType === 'album'">Total anual de canciones escuchadas del álbum</span>
            <span v-else>Total anual de canciones escuchadas del artista</span>
            "{{ searchQuery }}"
          </h2>
          <h2 v-else>Total Canciones Únicas Anuales Generales</h2>
          <DoughnutChart :chartData="yearlyDataSongs" />
        </div>
      </div>
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
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
}

.search-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.search-bar {
  display: flex;
  align-items: center;
  background-color: var(--spotify-light-gray);
  border-radius: 24px;
  padding: 8px 16px;
  flex-grow: 1;
  max-width: 600px;
}

.search-icon {
  color: var(--spotify-text-gray);
  margin-right: 12px;
}

.search-input {
  background: none;
  border: none;
  color: var(--spotify-white);
  font-size: 14px;
  width: 100%;
  outline: none;
}

.search-input::placeholder {
  color: var(--spotify-text-gray);
}

.type-filters {
  display: flex;
  gap: 12px;
}

.type-btn {
  background-color: transparent;
  color: var(--spotify-text-gray);
  border: 1px solid var(--spotify-text-gray);
  padding: 6px 16px;
  border-radius: 24px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.type-btn:hover {
  color: var(--spotify-white);
  border-color: var(--spotify-white);
}

.type-btn.active {
  background-color: var(--spotify-white);
  color: var(--spotify-black);
  border-color: var(--spotify-white);
}

.error-message {
  padding: 24px;
  text-align: center;
  color: var(--spotify-text-gray);
}

.results-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.yearly-charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.chart-wrapper {
  padding: 24px;
}

.chart-wrapper h2 {
  font-size: 20px;
  margin-bottom: 24px;
}

@media (max-width: 900px) {
  .yearly-charts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
