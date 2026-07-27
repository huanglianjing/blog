<script setup>
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const open = ref(false)
const keyword = ref('')
const searchRef = ref(null)
const inputRef = ref(null)

// 图标既是展开入口，展开后再点又是提交按钮。
async function onIconClick() {
  if (open.value) {
    submit()
    return
  }
  open.value = true
  await nextTick()
  inputRef.value?.focus()
}

function close() {
  open.value = false
}

function submit() {
  const q = keyword.value.trim()
  if (!q) return
  close()
  router.push({ name: 'search', query: { q } })
}

// 点击搜索区域以外的任意位置收起输入框。
function onDocumentClick(e) {
  if (open.value && searchRef.value && !searchRef.value.contains(e.target)) {
    close()
  }
}

onMounted(() => document.addEventListener('click', onDocumentClick))
onUnmounted(() => document.removeEventListener('click', onDocumentClick))
</script>

<template>
  <header class="topbar">
    <nav class="topbar-inner">
      <RouterLink class="brand" to="/">Moondo</RouterLink>
      <div class="nav-links">
        <RouterLink class="nav-link" to="/article">文章</RouterLink>
        <RouterLink class="nav-link" to="/category">分类</RouterLink>
        <RouterLink class="nav-link" to="/tag">标签</RouterLink>
        <div ref="searchRef" class="search" :class="{ open }">
          <input
            v-show="open"
            ref="inputRef"
            v-model="keyword"
            class="search-input"
            type="text"
            placeholder="搜索"
            @keyup.enter="submit"
            @keyup.esc="close"
          />
          <button class="search-btn" type="button" aria-label="搜索" @click="onIconClick">
            <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
              <circle cx="11" cy="11" r="7" fill="none" stroke="currentColor" stroke-width="2" />
              <line
                x1="16.2" y1="16.2" x2="21" y2="21"
                stroke="currentColor" stroke-width="2" stroke-linecap="round"
              />
            </svg>
          </button>
        </div>
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
  position: relative; /* 作为窄屏展开态搜索框的定位参照 */
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

.search {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.search-input {
  width: 10rem;
  padding: 0.3rem 0.6rem;
  font-size: 0.9rem;
  color: #333;
  border: 1px solid #ddd;
  border-radius: 4px;
  outline: none;
}

.search-input:focus {
  border-color: #999;
}

.search-btn {
  display: flex;
  align-items: center;
  padding: 0;
  border: none;
  background: none;
  color: #666;
  cursor: pointer;
  transition: color 0.2s ease;
}

.search-btn:hover {
  color: #000;
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

  /* 窄屏展开时输入框覆盖整条顶栏，避免挤压导航链接导致换行错位 */
  .search.open {
    position: absolute;
    inset: 0;
    z-index: 1;
    gap: 0.5rem;
    padding: 0 0.75rem;
    background: #ffffff;
  }

  .search.open .search-input {
    flex: 1;
    width: auto;
  }
}
</style>
