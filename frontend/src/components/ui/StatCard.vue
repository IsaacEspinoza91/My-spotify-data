<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  subtitle: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <div class="stat-card card">
    <h3 class="stat-title">{{ title }}</h3>
    <div class="stat-value text-green">{{ value }}</div>
    <div v-if="subtitle" class="stat-subtitle text-gray">{{ subtitle }}</div>
  </div>
</template>

<style scoped>
.stat-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
}

.stat-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--spotify-white);
}

.stat-value {
  font-size: 32px;
  font-weight: 800;
  letter-spacing: -1px;
}

.stat-subtitle {
  font-size: 13px;
  font-weight: 500;
}
</style>
