<script setup></script>

<template>
  <header class="topbar">
    <nav class="topbar-inner">
      <RouterLink class="brand" to="/">Moondo</RouterLink>
      <div class="nav-links">
        <RouterLink class="nav-link" to="/article">文章</RouterLink>
        <RouterLink class="nav-link" to="/category">分类</RouterLink>
        <RouterLink class="nav-link" to="/tag">标签</RouterLink>
      </div>
    </nav>
  </header>
</template>

<style scoped>
.topbar {
  position: sticky;
  top: 0;
  z-index: 10;
  background: #ffffff;
  border-bottom: 1px solid #f0f0f0;
}

.topbar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 960px;
  margin: 0 auto;
  padding: 0 1rem;
  height: 56px;
}

.brand {
  font-size: 1.25rem;
  font-weight: 800;
  letter-spacing: 0.05em;
  color: #333333;
  text-decoration: none;
}

.brand:hover {
  color: #000000;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.nav-link {
  font-size: 0.95rem;
  color: #666666;
  text-decoration: none;
  transition: color 0.2s ease;
}

.nav-link:hover {
  color: #000000;
}

.nav-link.router-link-active {
  color: #000000;
  font-weight: 600;
}

/* 窄屏收紧间距，避免超小屏导航拥挤 */
@media (max-width: 480px) {
  .topbar-inner {
    padding: 0 0.75rem;
  }

  .brand {
    font-size: 1.1rem;
    letter-spacing: 0.03em;
  }

  .nav-links {
    gap: 1rem;
  }

  .nav-link {
    font-size: 0.9rem;
  }
}
</style>
