<script setup>
import { computed } from 'vue'
import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-vue-next'

const props = defineProps({
  page: { type: Number, required: true },
  totalPages: { type: Number, required: true },
  total: { type: Number, required: true }
})

const emit = defineEmits(['update:page'])

const hasPrev = computed(() => props.page > 1)
const hasNext = computed(() => props.page < props.totalPages)

const prevPage = () => {
  if (hasPrev.value) emit('update:page', props.page - 1)
}

const nextPage = () => {
  if (hasNext.value) emit('update:page', props.page + 1)
}

const goToPage = (p) => {
  if (p !== props.page && p >= 1 && p <= props.totalPages) {
    emit('update:page', p)
  }
}

// Algorithm to calculate which page numbers to show
const visiblePages = computed(() => {
  const current = props.page
  const last = props.totalPages
  const delta = 2
  const left = current - delta
  const right = current + delta + 1
  const range = []
  const rangeWithDots = []
  let l

  for (let i = 1; i <= last; i++) {
    if (i === 1 || i === last || i >= left && i < right) {
      range.push(i)
    }
  }

  for (let i of range) {
    if (l) {
      if (i - l === 2) {
        rangeWithDots.push(l + 1)
      } else if (i - l !== 1) {
        rangeWithDots.push('...')
      }
    }
    rangeWithDots.push(i)
    l = i
  }

  return rangeWithDots
})

</script>

<template>
  <div class="pagination-container" v-if="totalPages > 1">
    <span class="pagination-info text-gray">
      {{ total }} resultados
    </span>
    
    <div class="pagination-controls">
      <button 
        class="icon-btn" 
        :disabled="!hasPrev" 
        @click="prevPage"
        title="Página anterior"
      >
        <ChevronLeft :size="20" />
      </button>
      
      <div class="page-numbers">
        <template v-for="(p, index) in visiblePages" :key="index">
          <span v-if="p === '...'" class="dots">
            <MoreHorizontal :size="16" />
          </span>
          <button 
            v-else 
            class="page-btn" 
            :class="{ active: p === page }"
            @click="goToPage(p)"
          >
            {{ p }}
          </button>
        </template>
      </div>

      <button 
        class="icon-btn" 
        :disabled="!hasNext" 
        @click="nextPage"
        title="Página siguiente"
      >
        <ChevronRight :size="20" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.pagination-container {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 16px;
  padding: 16px 0;
}

.pagination-info {
  font-size: 14px;
}

.pagination-controls {
  display: flex;
  gap: 8px;
}

.icon-btn {
  background-color: transparent;
  color: var(--spotify-text-gray);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.icon-btn:not(:disabled):hover {
  color: var(--spotify-white);
  background-color: var(--spotify-light-gray);
  transform: none;
}

.icon-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-btn {
  background-color: transparent;
  color: var(--spotify-text-gray);
  border: none;
  border-radius: 50%;
  min-width: 32px;
  height: 32px;
  padding: 0 8px;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.page-btn:hover {
  color: var(--spotify-white);
  background-color: var(--spotify-light-gray);
}

.page-btn.active {
  background-color: var(--spotify-green);
  color: var(--spotify-black);
  font-weight: 700;
}

.dots {
  color: var(--spotify-text-gray);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding-bottom: 4px;
}
</style>
