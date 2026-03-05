<script setup>
import { ref, watch } from 'vue'
import { Search, SlidersHorizontal } from 'lucide-vue-next'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'filter-change'])

const localFilters = ref({ ...props.modelValue })
const showAdvancedFilters = ref(false)

const toggleAdvancedFilters = () => {
  showAdvancedFilters.value = !showAdvancedFilters.value
}

let searchTimeout
const debouncedSearch = (e) => {
  const value = e.target.value
  localFilters.value.search = value
  
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    emitUpdate()
  }, 500)
}

const validateHour = (field) => {
  let val = localFilters.value[field]
  if (val !== '' && val !== null) {
    let num = parseInt(val)
    if (isNaN(num)) localFilters.value[field] = ''
    else if (num < 0) localFilters.value[field] = 0
    else if (num > 24) localFilters.value[field] = 24
  }
}

const emitUpdate = () => {
  emit('update:modelValue', localFilters.value)
  emit('filter-change')
}

// Watch for external changes (like resetting filters)
watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

</script>

<template>
  <div class="filters-container">
    <div class="main-filters">
      <div class="search-bar">
        <Search :size="20" class="search-icon" />
        <input 
          type="text" 
          v-model="localFilters.search" 
          @input="debouncedSearch"
          placeholder="Buscar un artista, álbum o canción..."
          class="search-input"
        />
      </div>
      <button 
        class="advanced-filters-btn" 
        :class="{ active: showAdvancedFilters }"
        @click="toggleAdvancedFilters"
      >
        <ScalarsHorizontal :size="20" />
        Filtros Avanzados
      </button>
    </div>

    <div v-if="showAdvancedFilters" class="advanced-filters-panel">
      <div class="filter-group">
        <label>Rango de Fechas</label>
        <div class="filter-inputs">
          <input 
            type="date" 
            v-model="localFilters.start_date" 
            @change="emitUpdate"
            class="filter-input date-input"
          />
          <span class="separator">-</span>
          <input 
            type="date" 
            v-model="localFilters.end_date" 
            @change="emitUpdate"
            class="filter-input date-input"
          />
        </div>
      </div>

      <div class="filter-group">
        <label>Rango de Horas (0-24)</label>
        <div class="filter-inputs">
          <input 
            type="number" 
            min="0" 
            max="24" 
            placeholder="Desde"
            v-model="localFilters.start_hour" 
            @input="validateHour('start_hour')"
            @change="emitUpdate"
            class="filter-input hour-input"
          />
          <span class="separator">-</span>
          <input 
            type="number" 
            min="0" 
            max="24" 
            placeholder="Hasta"
            v-model="localFilters.end_hour" 
            @input="validateHour('end_hour')"
            @change="emitUpdate"
            class="filter-input hour-input"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.filters-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  background-color: var(--spotify-dark-gray);
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 24px;
}

.main-filters {
  display: flex;
  gap: 16px;
  align-items: center;
}

.search-bar {
  display: flex;
  align-items: center;
  background-color: var(--spotify-light-gray);
  border-radius: 24px;
  padding: 8px 16px;
  flex-grow: 1;
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

.advanced-filters-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: transparent;
  color: var(--spotify-text-gray);
  border: 1px solid var(--spotify-text-gray);
  padding: 8px 16px;
  border-radius: 24px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.advanced-filters-btn:hover {
  color: var(--spotify-white);
  border-color: var(--spotify-white);
}

.advanced-filters-btn.active {
  background-color: var(--spotify-white);
  color: var(--spotify-black);
  border-color: var(--spotify-white);
}

.advanced-filters-panel {
  display: flex;
  gap: 32px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  font-size: 12px;
  font-weight: 600;
  color: var(--spotify-text-gray);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.filter-inputs {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-input {
  background-color: var(--spotify-light-gray);
  border: 1px solid transparent;
  color: var(--spotify-white);
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
  color-scheme: dark;
}

.filter-input:focus {
  outline: none;
  border-color: var(--spotify-text-gray);
}

.hour-input {
  width: 80px;
}

.separator {
  color: var(--spotify-text-gray);
  font-weight: bold;
}
</style>
