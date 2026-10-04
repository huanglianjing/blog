<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { theme, toggleTheme } from '../theme'
import {
  flattenCrumbs,
  commonPrefixLength,
  stepDelay,
  crumbsToText,
} from '../typewriter'

// 窄屏隐藏面包屑的断点，与下面 @media (max-width: 480px) 保持一致。
const NARROW = '(max-width: 480px)'

const router = useRouter()
const route = useRoute()

// 面包屑层级由路由 meta.crumbs 给出（站点名不含在内，模板里单独渲染）。
const crumbs = computed(() => {
  const fn = route.meta.crumbs
  return typeof fn === 'function' ? fn(route) : []
})

// 终态的字符序列。
const target = computed(() => flattenCrumbs(crumbs.value))

/*
 * 屏幕上这些字符取自哪个序列，以及已显示到第几个。
 * 回退阶段 displayed 仍是旧序列，否则删到一半的字会突然换成新内容；
 * 删到公共前缀后才切到新序列，接着往下打。
 */
const displayed = ref([])
const shown = ref(0)
// 打字动画进行中（用于显示光标）。
const typing = ref(false)

const visible = computed(() => displayed.value.slice(0, shown.value))

/*
 * 把字符按 group 合并成一段段，模板按段渲染。
 * 不能逐字符包 span：那样每个字都成了独立的 inline 盒子，
 * 换行与 text-overflow 的表现都会变，而且 DOM 节点数没必要那么多。
 */
const segments = computed(() => {
  const result = []
  for (const item of visible.value) {
    const last = result[result.length - 1]
    if (last && last.group === item.group) {
      last.text += item.char
    } else {
      result.push({ group: item.group, text: item.char, sep: item.sep, to: item.to })
    }
  }
  return result
})

// 读屏软件读这段完整文本，动画中间态对 AT 隐藏。
const a11yText = computed(() => crumbsToText(crumbs.value))

let timer = 0

function stopTimer() {
  window.clearTimeout(timer)
  timer = 0
}

// 直接落到终态，不做动画：首屏、窄屏（面包屑不可见）、用户要求减弱动效时走这里。
function settle() {
  stopTimer()
  displayed.value = target.value
  shown.value = target.value.length
  typing.value = false
}

/**
 * 逐字回退到公共前缀，再逐字打出新内容。
 *
 * 每一步单独 setTimeout，延时按各阶段的字符数摊算，所以长文章名不会拖很久。
 * 动画中途再次切换路由时，从当前屏幕状态接着算，不用先播完上一段。
 */
function animate(to) {
  stopTimer()
  typing.value = true

  // 与屏幕上实际显示的内容比公共前缀（不是与上一个终态比），
  // 这样连续快速切换时也不会把已经删掉的字当成还留着。
  const keep = Math.min(commonPrefixLength(displayed.value, to), shown.value)
  const eraseStep = stepDelay(shown.value - keep, true)
  const typeStep = stepDelay(to.length - keep, false)

  /*
   * 必须显式记住处于哪个阶段：只用 `shown > keep` 判断会在打字打过 keep 之后
   * 又满足回退条件，于是删一个打一个来回死循环。
   */
  let erasing = shown.value > keep

  function tick() {
    if (erasing) {
      // 回退：一次删一个字符（分隔符也算一个）。
      shown.value--
      if (shown.value <= keep) erasing = false
      timer = window.setTimeout(tick, eraseStep)
      return
    }
    // 已回退到公共前缀，换成新序列继续往下打。
    // 前缀部分字符相同，切换不会闪；样式（当前页 ↔ 可点上级）在此一并生效。
    if (displayed.value !== to) displayed.value = to
    if (shown.value < to.length) {
      shown.value++
      timer = window.setTimeout(tick, typeStep)
      return
    }
    // 到达终态，停下。
    typing.value = false
    timer = 0
  }

  tick()
}

/*
 * 盯扁平后的文本而非 crumbs 数组：computed 每次都返回新数组，而内容常常没变
 * （如总页数从 0 变 1，两次都不显示页码），盯数组会白跑一次删光重打。
 */
const targetText = computed(() => target.value.map((c) => c.char).join(''))

watch(targetText, () => {
  // 窄屏下面包屑 display: none，动画跑了也看不见，直接落终态省掉定时器。
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (reduceMotion || window.matchMedia(NARROW).matches) {
    settle()
    return
  }
  animate(target.value)
})

const open = ref(false)
const keyword = ref('')
const searchRef = ref(null)
const inputRef = ref(null)
const searchButtonRef = ref(null)

// 点击按钮和快捷键共用展开、聚焦逻辑。
async function focusSearch() {
  open.value = true
  await nextTick()
  inputRef.value?.focus()
}

// 图标既是展开入口，展开后再点又是提交按钮。
function onIconClick() {
  if (open.value) {
    submit()
    return
  }
  focusSearch()
}

// 输入文字或使用组合键时保留原有按键行为，避免截走输入框中的斜杠。
function onSearchShortcut(e) {
  if (e.key !== '/' || e.defaultPrevented || e.isComposing || e.repeat || e.ctrlKey || e.metaKey || e.altKey) return
  const target = e.target
  if (target instanceof HTMLElement && (target.isContentEditable || target.closest('input, textarea, select'))) return
  e.preventDefault()
  focusSearch()
}

// 以按钮中心作为新主题扩散的圆心。取按钮几何中心而非鼠标落点，
// 键盘回车触发时 event 里没有有效坐标，也能从同一处扩散。
function onThemeClick(e) {
  const r = e.currentTarget.getBoundingClientRect()
  toggleTheme({ x: r.left + r.width / 2, y: r.top + r.height / 2 })
}

// 收起时清空输入，下次展开不残留上一次的关键词。
function close() {
  // 输入框收起后仍留在 DOM 中，需主动移走焦点。
  inputRef.value?.blur()
  open.value = false
  keyword.value = ''
}

function onSearchEscape() {
  close()
  searchButtonRef.value?.focus({ preventScroll: true })
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

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onSearchShortcut)
  // 首屏不做动画：一进站就看着标题一个个打出来太吵，也会拖慢首屏观感。
  settle()
})

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onSearchShortcut)
  stopTimer()
})
</script>

<template>
  <header class="topbar">
    <nav class="topbar-inner">
      <div class="brand-area">
        <RouterLink class="brand" to="/">Moondo</RouterLink>
        <!--
          面包屑：站点名之后的各级，末级为当前页不带链接。
          逐字打出，中间态对读屏软件隐藏，另给一段完整文本（.sr-only）供其朗读。

          这里不加 v-if：元素在与不在会改变 .brand-area 的高度（见样式里的说明），
          切换路由时站点名就会跟着上下抖。空内容时它宽 0，常驻不占位。
        -->
        <span class="crumbs" :class="{ typing }" aria-hidden="true">
          <template v-for="seg in segments" :key="seg.group">
            <span v-if="seg.sep" class="sep">{{ seg.text }}</span>
            <RouterLink v-else-if="seg.to" class="crumb" :to="seg.to">{{ seg.text }}</RouterLink>
            <span v-else class="crumb current">{{ seg.text }}</span>
          </template>
        </span>
        <span v-if="a11yText" class="sr-only">面包屑：{{ a11yText }}</span>
      </div>
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
            aria-label="搜索关键词"
            aria-describedby="search-shortcut-tip"
            :tabindex="open ? 0 : -1"
            :aria-hidden="open ? 'false' : 'true'"
            @keyup.enter="submit"
            @keyup.esc="onSearchEscape"
          />
          <button
            ref="searchButtonRef"
            class="search-btn"
            type="button"
            aria-label="搜索"
            aria-keyshortcuts="/"
            aria-describedby="search-shortcut-tip"
            @click="onIconClick"
          >
            <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
              <circle cx="11" cy="11" r="7" fill="none" stroke="currentColor" stroke-width="2" />
              <line
                x1="16.2" y1="16.2" x2="21" y2="21"
                stroke="currentColor" stroke-width="2" stroke-linecap="round"
              />
            </svg>
          </button>
          <span id="search-shortcut-tip" class="search-tip" role="tooltip">按 <kbd>/</kbd> 搜索</span>
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

/* 站点名 + 面包屑：整体可压缩，过长时由面包屑末级省略，站点名不被挤掉 */
.brand-area {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  min-width: 0;
  overflow: hidden;
}

.brand {
  flex: none;
  font-size: 1.25rem;
  font-weight: 800;
  letter-spacing: 0.05em;
  color: var(--text);
  text-decoration: none;
}

.brand:hover {
  color: var(--text-strong);
}

/*
 * 面包屑的垂直几何必须与「显示了哪些字符」无关，否则打字动画每删 / 补一个字，
 * 站点名与面包屑就会上下抖一下（约 0.5～1.5px）。有两处会引起抖动：
 *
 * 1. 行高用 normal 时，行盒高度取决于实际用到的字体：只有 `/` 时走 system-ui，
 *    出现汉字时回退到 PingFang SC，行盒从 17px 变 20px。故此处写死行高。
 * 2. 各级 .crumb 为了省略号带 overflow: hidden，因而是滚动容器，
 *    基线由盒子下边缘合成而非文字基线（见 css-align）。若这里按 baseline 对齐，
 *    基线组会随「当前有没有 .crumb 字符」在文字基线与盒底之间来回跳。
 *    改为 center：各 item 行高一致、盒高一致，居中即等于对齐，与内容无关。
 */
.crumbs {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
  font-size: 0.9rem;
  line-height: 1.4;
  white-space: nowrap;
}

/*
 * 零宽字符撑出一条恒定的文字基线：.crumbs 作为 .brand-area 的 flex item
 * 要参与基线对齐（与站点名对齐），而它的基线取自第一个子项——
 * 内容为空时那就是光标（居中、无基线），基线又变了。
 * 有了它，无论内容多少，基线都由这个子项决定。它自身宽 0，
 * 负 margin 抵掉它引入的那一个 gap，所以水平位置分毫不变。
 */
.crumbs::before {
  content: '\200b';
  flex: none;
  margin-right: -0.5rem;
}

.sep {
  flex: none;
  color: var(--text-muted);
  user-select: none;
}

/* 各级都可压缩，超长（如很长的文章名）用省略号收尾 */
.crumb {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

/*
 * 可点击的上级：浅色常规字重，无下划线；当前页（末级）：深色加粗。
 * 靠颜色 + 字重两重区分，hover 时再加实线下划线给出可点的即时反馈。
 */
a.crumb {
  color: var(--text-secondary);
  font-weight: 400;
  text-decoration: none;
  transition: color 0.2s ease;
}

a.crumb:hover {
  color: var(--text-strong);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.crumb.current {
  color: var(--text-emphasis);
  font-weight: 600;
}

/*
 * 打字光标：只在动画期间出现，不常驻闪烁（顶栏一直闪很烦）。
 * 用 ::after 而非独立元素，省掉一个随动画增删的节点。
 * 它是 .crumbs 的 flex item，靠容器的 align-items: center 居中。
 */
.crumbs.typing::after {
  content: '';
  flex: none;
  width: 1px;
  height: 1.05em;
  /* 光标也是 flex item，会吃到 gap；抵消掉 gap 只留 2px，贴住最后一个字符 */
  margin-left: calc(2px - 0.5rem);
  background: var(--text-emphasis);
  animation: crumb-caret 0.9s step-end infinite;
}

@keyframes crumb-caret {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}

/* 仅供读屏软件：视觉上不可见，但仍在无障碍树中 */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  margin: -1px;
  padding: 0;
  overflow: hidden;
  clip-path: inset(50%);
  white-space: nowrap;
  border: 0;
}

.nav-links {
  display: flex;
  flex: none;
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

/* 收起时悬停一秒后显示提示；移开鼠标或展开搜索框时立即隐藏。 */
.search-tip {
  position: absolute;
  top: calc(100% + 0.75rem);
  right: 0;
  z-index: 3;
  padding: 0.4rem 0.6rem;
  color: var(--text);
  background: var(--bg-elevated);
  border: 1px solid var(--border-strong);
  border-radius: 4px;
  box-shadow: 0 2px 8px var(--shadow);
  font-size: 0.8rem;
  white-space: nowrap;
  pointer-events: none;
  visibility: hidden;
}

.search:not(.open):hover .search-tip {
  visibility: visible;
  transition: visibility 0s 1s;
}

.search-tip kbd {
  padding: 0 0.25rem;
  border: 1px solid var(--border-strong);
  border-radius: 3px;
  font-family: inherit;
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

  /* 超小屏导航已占满整行，面包屑挤到只剩省略号，不如整段隐藏 */
  .crumbs {
    display: none;
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
