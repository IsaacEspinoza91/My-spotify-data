<script setup>
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { 
  Chart as ChartJS, 
  Title, 
  Tooltip, 
  Legend, 
  ArcElement
} from 'chart.js'

ChartJS.register(ArcElement, Title, Tooltip, Legend)

const props = defineProps({
  chartData: {
    type: Object,
    required: true
  },
  title: {
    type: String,
    default: ''
  }
})

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'right',
      labels: {
        color: '#b3b3b3',
        padding: 20
      }
    },
    title: {
      display: !!props.title,
      text: props.title,
      color: '#FFFFFF'
    },
    tooltip: {
      backgroundColor: 'rgba(0,0,0,0.8)',
      titleColor: '#1ED760',
      bodyColor: '#FFFFFF',
      cornerRadius: 4
    }
  },
  borderWidth: 0
}))
</script>

<template>
  <div class="chart-container">
    <Doughnut :data="chartData" :options="chartOptions" />
  </div>
</template>

<style scoped>
.chart-container {
  position: relative;
  height: 300px;
  width: 100%;
}
</style>
