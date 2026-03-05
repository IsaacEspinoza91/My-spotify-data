<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'
import StatCard from '../components/ui/StatCard.vue'
import LoadingSpinner from '../components/ui/LoadingSpinner.vue'
import BarChart from '../components/ui/BarChart.vue'

const loading = ref(true)
const stats = ref(null)
const timeHabitsData = ref(null)
const dowHabitsData = ref(null)

const dayMap = {
  0: 'Dom', 1: 'Lun', 2: 'Mar', 3: 'Mié', 4: 'Jue', 5: 'Vie', 6: 'Sáb'
}

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
    <header class="view-header">
      <h1>Resumen General</h1>
    </header>

    <LoadingSpinner v-if="loading" />

    <div v-else class="dashboard-content">
      <!-- General Stats Cards -->
      <section class="stats-grid">
        <StatCard 
          title="Horas Escuchadas" 
          :value="formatNumber(Math.round(stats?.total_hours))"
          subtitle="Horas totales en tu historial"
        />
        <StatCard 
          title="Minutos Totales" 
          :value="formatNumber(Math.round(stats?.total_minutes))" 
        />
        <StatCard 
          title="Promedio Diario" 
          :value="formatNumber(stats?.average_daily_hours) + ' hrs'" 
        />
        <StatCard 
          title="Artistas Únicos" 
          :value="formatNumber(stats?.unique_artists)" 
        />
        <StatCard 
          title="Canciones Únicas" 
          :value="formatNumber(stats?.unique_songs)" 
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
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
