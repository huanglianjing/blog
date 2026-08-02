<script setup>
import TopBar from './components/TopBar.vue'
</script>

<template>
  <div class="app">
    <TopBar />
    <RouterView />
    <footer class="beian">
      <a
        class="beian-link"
        href="https://beian.miit.gov.cn/"
        target="_blank"
        rel="noopener noreferrer"
        >粤ICP备2021020498号</a
      >
      <a
        class="beian-link beian-police"
        href="https://beian.mps.gov.cn/#/query/webSearch?code=44030402004549"
        target="_blank"
        rel="noopener noreferrer"
        >粤公网安备 44030402004549号</a
      >
    </footer>
  </div>
</template>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.beian {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 0.25rem 1.25rem;
  padding: 0.75rem 1rem;
  background: var(--bg);
  font-size: 0.75rem;
  color: var(--text-muted);
}

.beian-link {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  color: var(--text-muted);
  text-decoration: none;
  transition: color 0.2s ease;
}

.beian-link:hover {
  color: var(--text-secondary);
}

</style>
