<script setup>
import { ChevronLeft, ChevronRight, User } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => router.back()
const goForward = () => router.forward()
</script>

<template>
  <div class="topbar">
    <div class="navigation-controls">
      <button class="nav-btn" @click="goBack" title="Go back">
        <ChevronLeft :size="24" />
      </button>
      <button class="nav-btn" @click="goForward" title="Go forward">
        <ChevronRight :size="24" />
      </button>
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
}

.navigation-controls {
  display: flex;
  gap: 8px;
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

.profile-section {
  display: flex;
  align-items: center;
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
