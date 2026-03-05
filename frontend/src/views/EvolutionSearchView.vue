<script setup>
import { ref } from 'vue'
import api from '../services/api'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import LineChart from '../components/ui/LineChart.vue'
import { Search } from 'lucide-vue-next'

const searchQuery = ref('')
const loading = ref(false)
const evolutionData = ref(null)
const errorMsg = ref('')

const searchArtist = async () => {
  loading.value = true
  errorMsg.value = ''
  evolutionData.value = null
  
  try {
    const response = await api.getEvolution({ search: searchQuery.value })
    const data = response.data
    
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
      <p class="text-gray subtitle">Busca un artista para ver cómo ha evolucionado tu tiempo escuchándolo a lo largo de los meses.</p>
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
    </div>

    <LoadingSpinner v-if="loading" />

    <div v-else-if="errorMsg" class="error-message card">
      {{ errorMsg }}
    </div>

    <div v-else-if="evolutionData" class="chart-wrapper card">
      <h2 v-if="searchQuery">Horas escuchadas de "{{ searchQuery }}"</h2>
      <h2 v-else>Horas escuchadas generales</h2>
      <LineChart :chartData="evolutionData" />
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

.error-message {
  padding: 24px;
  text-align: center;
  color: var(--spotify-text-gray);
}

.chart-wrapper {
  padding: 24px;
}

.chart-wrapper h2 {
  font-size: 20px;
  margin-bottom: 24px;
}
</style>
