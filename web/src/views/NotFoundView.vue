<script setup></script>

<template>
  <main class="notfound">
    <h1 class="code">404</h1>
  </main>
</template>

<style scoped>
.notfound {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  flex: 1;
  background: #ffffff;
}

.code {
  font-size: clamp(3rem, 15vw, 8rem);
  font-weight: 800;
  letter-spacing: 0.05em;
  color: #808080;
  user-select: none;
}
</style>
