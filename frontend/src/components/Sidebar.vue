<script setup>
import { 
  Home, 
  Mic2, 
  Music, 
  Disc, 
  LineChart, 
  Gift 
} from 'lucide-vue-next'
import { useRoute } from 'vue-router'

const route = useRoute()

const navItems = [
  { name: 'Dashboard', path: '/', icon: Home },
  { name: 'Top Artists', path: '/artists', icon: Mic2 },
  { name: 'Top Songs', path: '/songs', icon: Music },
  { name: 'Top Albums', path: '/albums', icon: Disc },
  { name: 'Evolution', path: '/evolution', icon: LineChart },
  { name: 'Wrapped', path: '/wrapped', icon: Gift },
]
</script>

<template>
  <div class="sidebar">
    <div class="logo-area">
      <img src="../assets/Spotify_Logo_White.png" alt="Spotify Logo" class="spotify-logo" />
      <span class="logo-text">Stats</span>
    </div>

    <nav class="nav-links">
      <router-link 
        v-for="item in navItems" 
        :key="item.name" 
        :to="item.path"
        class="nav-item"
        :class="{ active: route.path === item.path }"
      >
        <component :is="item.icon" class="icon" :size="24" :stroke-width="route.path === item.path ? 2.5 : 2" />
        <span class="label">{{ item.name }}</span>
      </router-link>
    </nav>
  </div>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.logo-area {
  padding: 0 12px 24px 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.spotify-logo {
  height: 32px;
  width: auto;
  object-fit: contain;
}

.logo-text {
  font-size: 24px;
  color: var(--spotify-white);
  font-weight: 700;
  letter-spacing: -1px;
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  border-radius: 4px;
  color: var(--spotify-text-gray);
  font-weight: 600;
  transition: color 0.2s ease, background-color 0.2s ease;
}

.nav-item:hover {
  color: var(--spotify-white);
  text-decoration: none;
}

.nav-item.active {
  color: var(--spotify-white);
  background-color: var(--spotify-light-gray);
}

.icon {
  flex-shrink: 0;
}

.label {
  font-size: 14px;
}
</style>
