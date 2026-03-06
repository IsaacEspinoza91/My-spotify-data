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
  },
  image: {
    type: String,
    default: null
  }
})
</script>

<template>
  <div class="good-morning-card">
    <div class="card-image-wrapper" v-if="image">
      <img :src="image" alt="Stat Cover" class="card-image" />
    </div>
    <div class="card-content">
      <h3 class="stat-title">{{ title }}</h3>
      <div class="stat-value text-green">{{ value }}</div>
      <div v-if="subtitle" class="stat-subtitle text-gray">{{ subtitle }}</div>
    </div>
  </div>
</template>

<style scoped>
.good-morning-card {
  display: flex;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
  height: 80px;
  transition: background-color 0.3s ease;
  cursor: pointer;
}

.good-morning-card:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.card-image-wrapper {
  width: 80px;
  height: 80px;
  flex-shrink: 0;
  box-shadow: 4px 0 12px rgba(0, 0, 0, 0.2);
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 16px;
  flex-grow: 1;
}

.stat-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--spotify-white);
  margin-bottom: 2px;
}

.stat-value {
  font-size: 20px;
  font-weight: 800;
  letter-spacing: -0.5px;
}

.stat-subtitle {
  font-size: 11px;
  font-weight: 500;
  margin-top: 2px;
}
</style>
