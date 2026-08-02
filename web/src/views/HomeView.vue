<script setup></script>

<template>
  <main class="home">
    <h3 class="title">Moondo</h3>
  </main>
</template>

<style scoped>
.home {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  flex: 1;
  background: var(--bg);
}

.title {
  font-size: clamp(1.5rem, 5vw, 3rem);
  font-weight: 800;
  letter-spacing: 0.05em;
  color: var(--text-hint);
  user-select: none;
}
</style>
