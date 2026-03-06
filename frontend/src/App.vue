<script setup>
import { RouterView } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import TopBar from './components/TopBar.vue'
</script>

<template>
  <div class="app-layout">
    <aside class="sidebar-area">
      <Sidebar />
    </aside>
    <main class="main-area">
      <div class="topbar-area">
        <TopBar />
      </div>
      <div class="content-area">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  grid-template-areas: "sidebar main";
  height: 100vh;
  width: 100vw;
  background-color: var(--spotify-black);
}

.sidebar-area {
  grid-area: sidebar;
  background-color: #000000;
  padding: 24px 12px;
  display: flex;
  flex-direction: column;
}

.main-area {
  grid-area: main;
  background-color: var(--spotify-dark-gray);
  border-radius: 8px;
  margin: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

.topbar-area {
  height: 64px;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  position: sticky;
  top: 0;
  background-color: rgba(18, 18, 18, 0.8);
  backdrop-filter: blur(8px);
  z-index: 10;
}

.content-area {
  flex: 1;
  padding: 0 24px 24px 24px;
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
  overflow-y: auto;
}
</style>
