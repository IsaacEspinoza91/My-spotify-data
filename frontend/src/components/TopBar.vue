<script setup>
import { ref, computed } from 'vue'
import { ChevronLeft, ChevronRight, User, Search, Library } from 'lucide-vue-next'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const goBack = () => router.back()
const goForward = () => router.forward()

const searchQuery = ref('')
const showDropdown = ref(false)

const availableViews = [
  { name: 'Inicio', path: '/' },
  { name: 'Top Artistas', path: '/artists' },
  { name: 'Top Canciones', path: '/songs' },
  { name: 'Top Álbumes', path: '/albums' },
  { name: 'Evolución', path: '/evolution' },
  { name: 'Wrapped', path: '/wrapped' }
]

const filteredViews = computed(() => {
  if (!searchQuery.value) return availableViews
  const query = searchQuery.value.toLowerCase()
  return availableViews.filter(view => view.name.toLowerCase().includes(query))
})

const navigateTo = (path) => {
  router.push(path)
  searchQuery.value = ''
  showDropdown.value = false
}

const handleBlur = () => {
  setTimeout(() => {
    showDropdown.value = false
  }, 200)
}
</script>

<template>
  <div class="topbar">
    <div class="navigation-controls">
      <button class="nav-btn" @click="goBack" title="Volver atrás">
        <ChevronLeft :size="24" />
      </button>
      <button class="nav-btn" @click="goForward" title="Ir adelante">
        <ChevronRight :size="24" />
      </button>
    </div>

    <div class="center-controls">

      <div class="global-search-container" :class="{ focused: showDropdown }">
        <Search class="search-icon" :size="20" />
        <input 
          type="text" 
          v-model="searchQuery" 
          class="search-input" 
          placeholder="¿A dónde quieres ir?" 
          @focus="showDropdown = true" 
          @blur="handleBlur"
        />
        <div class="divider"></div>
        <button class="browse-btn" title="Explorar">
          <Library :size="20" />
        </button>

        <div v-if="showDropdown && filteredViews.length > 0" class="dropdown-menu">
          <div 
            v-for="view in filteredViews" 
            :key="view.path" 
            class="dropdown-item"
            @click.stop="navigateTo(view.path)"
          >
            <span>{{ view.name }}</span>
          </div>
          <div v-if="filteredViews.length === 0" class="dropdown-empty text-gray">
            No se encontraron secciones
          </div>
        </div>
      </div>
    </div>
    
    <div class="profile-section">
      <div class="profile-pill">
        <div class="avatar">
          <User :size="16" />
        </div>
        <span class="username">Isaac Espinoza</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  position: relative;
}

.navigation-controls {
  display: flex;
  gap: 8px;
  flex: 1; /* Helps keep center-controls strictly centered */
}

.nav-btn {
  background-color: rgba(0, 0, 0, 0.7);
  color: var(--spotify-white);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.1s ease, background-color 0.2s ease;
}

.nav-btn:hover {
  background-color: rgba(0, 0, 0, 0.9);
  transform: scale(1.04);
}

.nav-btn:active {
  transform: scale(1);
}

.center-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
}

.home-btn {
  background-color: #242424;
  color: var(--spotify-text-gray);
  border: none;
  border-radius: 50%;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.1s ease, background-color 0.2s ease, color 0.2s ease;
}

.home-btn:hover {
  background-color: #2a2a2a;
  color: var(--spotify-white);
  transform: scale(1.04);
}

.home-btn.active {
  color: var(--spotify-white);
}

.home-btn:active {
  transform: scale(1);
}

.global-search-container {
  display: flex;
  align-items: center;
  background-color: #242424;
  border-radius: 24px;
  height: 48px;
  width: 480px;
  padding: 0 16px;
  position: relative;
  border: 1px solid transparent;
  transition: background-color 0.2s ease, border-color 0.2s ease;
}

.global-search-container:hover {
  background-color: #2a2a2a;
}

.global-search-container.focused {
  background-color: #2a2a2a;
  border-color: #727272;
}

.search-icon {
  color: var(--spotify-text-gray);
  flex-shrink: 0;
}

.global-search-container.focused .search-icon {
  color: var(--spotify-white);
}

.search-input {
  background: transparent;
  border: none;
  color: var(--spotify-white);
  width: 100%;
  height: 100%;
  padding-left: 12px;
  font-size: 16px;
  font-weight: 500;
  outline: none;
}

.search-input::placeholder {
  color: #b3b3b3;
}

.divider {
  height: 24px;
  width: 1px;
  background-color: #727272;
  margin: 0 12px;
  flex-shrink: 0;
}

.browse-btn {
  background: transparent;
  border: none;
  color: var(--spotify-text-gray);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 4px;
  transition: color 0.2s ease;
  flex-shrink: 0;
}

.browse-btn:hover {
  color: var(--spotify-white);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  width: 100%;
  background-color: #282828;
  border-radius: 8px;
  padding: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.6);
  z-index: 100;
  max-height: 400px;
  overflow-y: auto;
}

.dropdown-item {
  padding: 12px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  color: var(--spotify-white);
  transition: background-color 0.2s ease;
}

.dropdown-item:hover {
  background-color: #3e3e3e;
}

.dropdown-empty {
  padding: 16px;
  text-align: center;
}

.profile-section {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex: 1; /* Helps keep center-controls strictly centered */
}

.profile-pill {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: rgba(0, 0, 0, 0.7);
  padding: 2px 12px 2px 2px;
  border-radius: 20px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.profile-pill:hover {
  background-color: var(--spotify-light-gray);
}

.avatar {
  background-color: var(--spotify-dark-gray);
  border-radius: 50%;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--spotify-white);
}

.username {
  font-size: 14px;
  font-weight: 700;
  color: var(--spotify-white);
}
</style>
