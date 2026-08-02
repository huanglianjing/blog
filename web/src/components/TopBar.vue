<script setup>
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { theme, toggleTheme } from '../theme'

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

// 以按钮中心作为新主题扩散的圆心。取按钮几何中心而非鼠标落点，
// 键盘回车触发时 event 里没有有效坐标，也能从同一处扩散。
function onThemeClick(e) {
  const r = e.currentTarget.getBoundingClientRect()
  toggleTheme({ x: r.left + r.width / 2, y: r.top + r.height / 2 })
}

// 收起时清空输入，下次展开不残留上一次的关键词。
function close() {
  open.value = false
  keyword.value = ''
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
        <!-- 主题切换：浅色下显示月亮（点击转深色），深色下显示太阳 -->
        <button
          class="theme-btn"
          type="button"
          :aria-label="theme === 'dark' ? '切换到浅色模式' : '切换到深色模式'"
          :title="theme === 'dark' ? '切换到浅色模式' : '切换到深色模式'"
          @click="onThemeClick"
        >
          <svg
            v-if="theme === 'dark'"
            viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"
            fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
          >
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2M12 20v2M2 12h2M20 12h2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
          </svg>
          <svg
            v-else
            viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"
            fill="none" stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round"
          >
            <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" />
          </svg>
        </button>
        <div ref="searchRef" class="search" :class="{ open }">
          <!-- 输入框常驻 DOM（不用 v-show），否则展开/收起没有过渡动画 -->
          <input
            ref="inputRef"
            v-model="keyword"
            class="search-input"
            type="text"
            placeholder="搜索"
            :tabindex="open ? 0 : -1"
            :aria-hidden="open ? 'false' : 'true'"
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
  background: var(--bg);
  border-bottom: 1px solid var(--border-subtle);
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
  color: var(--text);
  text-decoration: none;
}

.brand:hover {
  color: var(--text-strong);
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.nav-link {
  font-size: 0.95rem;
  color: var(--text-secondary);
  text-decoration: none;
  transition: color 0.2s ease;
}

.nav-link:hover {
  color: var(--text-strong);
}

.nav-link.router-link-active {
  color: var(--text-strong);
  font-weight: 600;
}

/* 主题切换按钮：与搜索图标同规格的纯图标按钮 */
.theme-btn {
  display: flex;
  align-items: center;
  padding: 0;
  border: none;
  background: none;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s ease;
}

.theme-btn:hover {
  color: var(--text-strong);
}

.search {
  position: relative; /* 输入框以图标为基准向左浮出 */
  display: flex;
  align-items: center;
}

/*
 * 输入框绝对定位在图标左侧，浮在「文章 / 分类 / 标签」与主题按钮之上
 * （靠自身与顶栏同色的背景遮挡），
 * 不参与 flex 布局，因此展开时不会把导航链接往左顶。
 * 收起态宽度为 0 且透明，配合 transition 形成展开 / 收起动画。
 */
.search-input {
  position: absolute;
  top: 50%;
  right: calc(100% + 0.4rem);
  z-index: 2;
  transform: translateY(-50%);
  width: 0;
  padding: 0.3rem 0;
  font-size: 0.9rem;
  color: var(--text);
  background: var(--bg);
  border: 1px solid transparent;
  border-radius: 4px;
  outline: none;
  opacity: 0;
  pointer-events: none;
  transition:
    width 0.25s ease,
    padding 0.25s ease,
    opacity 0.2s ease,
    border-color 0.2s ease;
}

.search.open .search-input {
  /* 左边缘越过「文章」再多两个字宽，把整组导航链接都盖住 */
  width: 14.5rem;
  padding: 0.3rem 0.6rem;
  border-color: var(--border-strong);
  opacity: 1;
  pointer-events: auto;
}

.search-input:focus {
  border-color: var(--border-hover);
}

.search-btn {
  display: flex;
  align-items: center;
  padding: 0;
  border: none;
  background: none;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s ease;
}

.search-btn:hover {
  color: var(--text-strong);
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

  /*
   * 窄屏导航 gap 收紧到 1rem、字号也更小，盖住「文章」多两个字只需 12rem；
   * 320px 下这已是上限，再宽就会顶到站点名。
   */
  .search.open .search-input {
    width: 12rem;
  }
}
</style>
