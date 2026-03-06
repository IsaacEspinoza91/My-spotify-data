<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../services/api'
import StatCard from '../components/ui/StatCard.vue'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import BarChart from '../components/ui/BarChart.vue'

import imgHours from '../assets/cover_hours.png'
import imgMinutes from '../assets/cover_minutes.png'
import imgDaily from '../assets/cover_daily.png'
import imgArtists from '../assets/cover_artists.png'
import imgSongs from '../assets/cover_songs.png'
import imgPlays from '../assets/cover_plays.png'

const loading = ref(true)
const stats = ref(null)
const timeHabitsData = ref(null)
const dowHabitsData = ref(null)

const dayMap = {
  0: 'Dom', 1: 'Lun', 2: 'Mar', 3: 'Mié', 4: 'Jue', 5: 'Vie', 6: 'Sáb'
}

const greeting = computed(() => {
  const currentHour = new Date().getHours()
  if (currentHour >= 5 && currentHour < 12) {
    return 'Buenos días'
  } else if (currentHour >= 12 && currentHour < 20) {
    return 'Buenas tardes'
  } else {
    return 'Buenas noches'
  }
})

const fetchData = async () => {
  loading.value = true
  try {
    const [statsRes, timeRes, dowRes] = await Promise.all([
      api.getStats(),
      api.getHabits('time'),
      api.getHabits('dow')
    ])

    stats.value = statsRes.data

    // Prepare chart data for Time of Day
    if (timeRes.data) {
      timeHabitsData.value = {
        labels: timeRes.data.map(item => item.label),
        datasets: [{
          label: 'Reproducciones',
          backgroundColor: '#1ED760',
          borderRadius: 4,
          data: timeRes.data.map(item => item.count)
        }]
      }
    }

    // Prepare chart data for Day of Week
    if (dowRes.data) {
      // API returns num_day 0-6
      const sortedDow = [...dowRes.data].sort((a, b) => a.num_day - b.num_day)
      dowHabitsData.value = {
        labels: sortedDow.map(item => dayMap[item.num_day]),
        datasets: [{
          label: 'Reproducciones',
          backgroundColor: '#1ED760',
          borderRadius: 4,
          data: sortedDow.map(item => item.count)
        }]
      }
    }

  } catch (error) {
    console.error("Error fetching dashboard data", error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const formatNumber = (num) => {
  return new Intl.NumberFormat('es-CL').format(num || 0)
}
</script>

<template>
  <div class="dashboard-view view-container">
    <div class="dashboard-hero-bg"></div>
    
    <div class="dashboard-content-wrapper relative-z">
      <header class="view-header">
        <h1>{{ greeting }}</h1>
      </header>

      <LoadingSpinner v-if="loading" />

      <div v-else class="dashboard-content">
        <!-- General Stats Cards -->
        <section class="stats-grid">
          <StatCard 
            title="Horas Escuchadas" 
            :value="formatNumber(Math.round(stats?.total_hours))"
            :image="imgHours"
          />
          <StatCard 
            title="Minutos Totales" 
            :value="formatNumber(Math.round(stats?.total_minutes))" 
            :image="imgMinutes"
          />
          <StatCard 
            title="Promedio Diario" 
            :value="formatNumber(stats?.average_daily_hours) + ' hrs'" 
            :image="imgDaily"
          />
          <StatCard 
            title="Días en el Historial" 
            :value="formatNumber(Math.round(stats?.total_hours / stats?.average_daily_hours))" 
            :image="imgPlays"
          />
          <StatCard 
            title="Artistas Únicos" 
            :value="formatNumber(stats?.unique_artists)" 
            :image="imgArtists"
          />
          <StatCard 
            title="Canciones Únicas" 
            :value="formatNumber(stats?.unique_songs)" 
            :image="imgSongs"
          />
        </section>

      <!-- Habits Charts -->
      <section class="charts-section">
        <div class="chart-card card">
          <h3>Reproducciones por momento del día</h3>
          <BarChart v-if="timeHabitsData" :chartData="timeHabitsData" />
        </div>
        
        <div class="chart-card card">
          <h3>Reproducciones por día de la semana</h3>
          <BarChart v-if="dowHabitsData" :chartData="dowHabitsData" />
        </div>
      </section>
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

.dashboard-hero-bg {
  position: absolute;
  top: -24px;
  left: -24px;
  right: -24px;
  height: 200px;
  background: linear-gradient(180deg, rgba(83, 83, 100, 1) 0%, rgba(24, 24, 24, 1) 100%);
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

.view-header {
  margin-top: 32px;
}

.view-header h1 {
  font-size: 32px;
  font-weight: 800;
  color: var(--spotify-white);
  letter-spacing: -1px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px 24px;
  margin-bottom: 48px;
}

.charts-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 16px;
}

.chart-card {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.chart-card h3 {
  font-size: 18px;
  margin-bottom: 8px;
}
</style>
