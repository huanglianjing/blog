<script setup>
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import hljs from 'highlight.js/lib/common'
// 浅色为基线主题，深色主题的选择器带 html[data-theme='dark'] 前缀覆盖其上。
import 'highlight.js/styles/github.css'
import '../assets/hljs-dark.css'

const route = useRoute()

const detail = ref(null)
const loading = ref(false)
const error = ref('')
const bodyRef = ref(null)
const toc = ref([])
const activeHeadingId = ref('')
const showBackTop = ref(false)

let headingElements = []
let scrollFrame = 0

// 复制按钮的两个图标：默认双层方框，复制成功后切换为对勾（由 .copied 控制显隐）。
const COPY_ICON_SVG = `
<svg class="code-copy-icon icon-copy" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
  <rect x="8" y="8" width="14" height="14" rx="2" />
  <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
</svg>`
const CHECK_ICON_SVG = `
<svg class="code-copy-icon icon-check" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
  <path d="M20 6 9 17l-5-5" />
</svg>`

async function fetchDetail(title) {
  loading.value = true
  error.value = ''
  detail.value = null
  toc.value = []
  activeHeadingId.value = ''
  headingElements = []
  try {
    const resp = await fetch(`/api/article/detail?title=${encodeURIComponent(title)}`)
    const data = await resp.json()
    if (data.code !== 0) {
      throw new Error(data.msg || '请求失败')
    }
    detail.value = data.data
    // 先结束 loading，让正文（v-else-if）渲染进 DOM，再提取标题生成目录。
    loading.value = false
    await nextTick()
    buildToc()
    enhanceCodeBlocks()
  } catch (e) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

// 从渲染后的正文中提取各级标题，生成目录。
function buildToc() {
  if (!bodyRef.value) return
  const headings = bodyRef.value.querySelectorAll('h1, h2, h3, h4, h5, h6')
  const items = []
  headings.forEach((el, i) => {
    // 标题若无 id，补一个以便锚点跳转。
    if (!el.id) el.id = `heading-${i}`
    items.push({
      id: el.id,
      text: el.textContent.trim(),
      level: Number(el.tagName.slice(1)),
    })
  })
  toc.value = items
  headingElements = Array.from(headings)
  updateActiveHeading()

  // 支持直接访问带有 #heading 的文章链接。
  if (route.hash) {
    let id = route.hash.slice(1)
    try {
      id = decodeURIComponent(id)
    } catch {
      // hash 不是合法的 URI 编码时，直接使用原值查找。
    }
    nextTick(() => scrollTo(id, false))
  }
}

// 为正文中的每个代码块添加语法高亮、行号、语言标签和复制按钮。
function enhanceCodeBlocks() {
  if (!bodyRef.value) return
  const blocks = bodyRef.value.querySelectorAll('pre')
  blocks.forEach((pre) => {
    // 避免重复包裹（例如同一文章多次调用）。
    if (pre.parentElement?.classList.contains('code-block')) return

    const code = pre.querySelector('code')
    // 从 code 的 class="language-xxx" 中解析语言名。
    let lang = ''
    if (code) {
      const cls = [...code.classList].find((c) => c.startsWith('language-'))
      if (cls) lang = cls.slice('language-'.length)
    }

    // 语法高亮：指定语言且被 hljs 支持时按该语言高亮，否则自动识别。
    if (code) {
      const original = code.textContent
      try {
        const result =
          lang && hljs.getLanguage(lang)
            ? hljs.highlight(original, { language: lang })
            : hljs.highlightAuto(original)
        code.innerHTML = result.value
        code.classList.add('hljs')
        // 自动识别出的语言用于工具栏标签展示。
        if (!lang && result.language) lang = result.language
      } catch {
        // 高亮失败时保留原始文本，不影响展示。
      }

      // 行号列：代码不换行，物理行即逻辑行，用独立 gutter 保证对齐。
      const lineCount = original.replace(/\n$/, '').split('\n').length
      const gutter = document.createElement('span')
      gutter.className = 'code-gutter'
      gutter.setAttribute('aria-hidden', 'true')
      gutter.textContent = Array.from({ length: lineCount }, (_, i) => i + 1).join('\n')
      pre.insertBefore(gutter, code)
    }

    // 用一个容器包裹 pre，工具栏绝对定位在右上角。
    const wrapper = document.createElement('div')
    wrapper.className = 'code-block'

    const bar = document.createElement('div')
    bar.className = 'code-bar'

    const label = document.createElement('span')
    label.className = 'code-lang'
    label.textContent = lang || 'text'

    const btn = document.createElement('button')
    btn.type = 'button'
    btn.className = 'code-copy'
    btn.title = '复制'
    btn.setAttribute('aria-label', '复制代码')
    btn.innerHTML = COPY_ICON_SVG + CHECK_ICON_SVG

    bar.append(label, btn)

    pre.parentNode.insertBefore(wrapper, pre)
    wrapper.append(bar, pre)
  })
}

// 点击复制按钮时，把对应代码块的文本写入剪贴板（事件委托）。
async function handleCopyClick(e) {
  const btn = e.target.closest('.code-copy')
  if (!btn) return
  const pre = btn.closest('.code-block')?.querySelector('pre')
  if (!pre) return

  // 只取 code 的文本，排除行号列（.code-gutter）。
  const text = (pre.querySelector('code') ?? pre).textContent
  try {
    await navigator.clipboard.writeText(text)
  } catch {
    // 剪贴板 API 不可用时降级到 execCommand。
    const ta = document.createElement('textarea')
    ta.value = text
    ta.style.position = 'fixed'
    ta.style.opacity = '0'
    document.body.appendChild(ta)
    ta.select()
    try {
      document.execCommand('copy')
    } catch {
      // 忽略失败。
    }
    document.body.removeChild(ta)
  }

  btn.title = '已复制'
  btn.setAttribute('aria-label', '已复制')
  btn.classList.add('copied')
  window.setTimeout(() => {
    btn.title = '复制'
    btn.setAttribute('aria-label', '复制代码')
    btn.classList.remove('copied')
  }, 1500)
}

// 点击目录项，平滑滚动到对应标题。
function scrollTo(id, smooth = true) {
  const el = document.getElementById(id)
  if (!el) return

  activeHeadingId.value = id
  el.scrollIntoView({ behavior: smooth ? 'smooth' : 'auto', block: 'start' })
  window.history.replaceState(null, '', `#${encodeURIComponent(id)}`)
}

// 根据滚动位置高亮当前章节。用 requestAnimationFrame 限制滚动事件的计算频率。
function updateActiveHeading() {
  if (!headingElements.length) return

  const headerOffset = 96
  let current = headingElements[0]
  for (const heading of headingElements) {
    if (heading.getBoundingClientRect().top > headerOffset) break
    current = heading
  }
  activeHeadingId.value = current.id
}

function handleScroll() {
  if (scrollFrame) return
  scrollFrame = window.requestAnimationFrame(() => {
    updateActiveHeading()
    // 向下滚动超过一屏时展示回到顶部按钮。
    showBackTop.value = window.scrollY > window.innerHeight
    scrollFrame = 0
  })
}

// 平滑滚动回页面顶部。
function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

window.addEventListener('scroll', handleScroll, { passive: true })

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
  if (scrollFrame) window.cancelAnimationFrame(scrollFrame)
})

// 路由参数变化时（含首次进入）重新加载。
watch(
  () => route.params.title,
  (title) => {
    if (title) fetchDetail(title)
  },
  { immediate: true },
)
</script>

<template>
  <main class="page">
    <div class="layout">
      <div class="content">
        <p v-if="loading" class="hint">加载中…</p>
        <p v-else-if="error" class="hint error">{{ error }}</p>

        <article v-else-if="detail">
          <h1 class="title">{{ detail.title }}</h1>
          <div class="meta">
            <span class="date">{{ detail.date }}</span>
            <RouterLink
              v-if="detail.category_name"
              class="category"
              :to="`/category/${encodeURIComponent(detail.category_name)}`"
            >{{ detail.category_name }}</RouterLink>
            <RouterLink
              v-for="tag in detail.tags"
              :key="tag"
              class="tag"
              :to="`/tag/${encodeURIComponent(tag)}`"
            >{{ tag }}</RouterLink>
          </div>
          <div ref="bodyRef" class="body" v-html="detail.content" @click="handleCopyClick"></div>
        </article>
      </div>

      <aside v-if="toc.length" class="toc">
        <nav class="toc-inner" aria-label="文章目录">
          <p class="toc-title">目录</p>
          <ul class="toc-list">
            <li
              v-for="item in toc"
              :key="item.id"
              :class="`toc-item toc-level-${item.level}`"
            >
              <a
                :href="`#${encodeURIComponent(item.id)}`"
                :class="{ active: activeHeadingId === item.id }"
                :aria-current="activeHeadingId === item.id ? 'location' : undefined"
                :title="item.text"
                @click.prevent="scrollTo(item.id)"
              >{{ item.text }}</a>
            </li>
          </ul>
        </nav>
      </aside>
    </div>

    <Transition name="fade">
      <button
        v-if="showBackTop"
        class="back-top"
        type="button"
        aria-label="回到顶部"
        @click="scrollToTop"
      >
        <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
          <path
            d="M12 5l-7 7h4v7h6v-7h4z"
            fill="currentColor"
          />
        </svg>
      </button>
    </Transition>
  </main>
</template>

<style scoped>
.page {
  flex: 1;
  width: 100%;
  background: var(--bg);
}

/* 回到顶部按钮：固定在右下角 */
.back-top {
  position: fixed;
  right: 2rem;
  bottom: 2rem;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  padding: 0;
  color: var(--text-secondary);
  background: var(--bg-elevated);
  border: 1px solid var(--border-strong);
  border-radius: 50%;
  box-shadow: 0 2px 8px var(--shadow);
  cursor: pointer;
  transition: color 0.2s ease, border-color 0.2s ease;
}

.back-top:hover {
  color: var(--text-strong);
  border-color: var(--border-hover);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 900px) {
  .back-top {
    right: 1rem;
    bottom: 1rem;
  }
}

.layout {
  display: flex;
  gap: 2rem;
  max-width: 1100px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
  align-items: flex-start;
}

.content {
  flex: 1;
  min-width: 0;
}

/* 右侧目录：随页面滚动固定 */
.toc {
  flex: 0 0 220px;
  position: sticky;
  top: 76px;
  max-height: calc(100vh - 96px);
  overflow-y: auto;
}

.toc-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text);
  margin: 0 0 0.75rem;
}

.toc-list {
  list-style: none;
  margin: 0;
  padding: 0;
  border-left: 1px solid var(--border);
}

.toc-item a {
  display: block;
  padding: 0.25rem 0 0.25rem 0.75rem;
  font-size: 0.85rem;
  line-height: 1.4;
  color: var(--text-secondary);
  text-decoration: none;
  border-left: 2px solid transparent;
  margin-left: -1px;
  cursor: pointer;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toc-item a:hover {
  color: var(--text-strong);
  border-left-color: var(--text-secondary);
}

.toc-item a.active {
  color: var(--text-emphasis);
  font-weight: 600;
  border-left-color: var(--text);
}

/* 按标题级别缩进 */
.toc-level-1 a { padding-left: 0.75rem; }
.toc-level-2 a { padding-left: 1.5rem; }
.toc-level-3 a { padding-left: 2.25rem; }
.toc-level-4 a { padding-left: 3rem; }
.toc-level-5 a { padding-left: 3.75rem; }
.toc-level-6 a { padding-left: 4.5rem; }

/* 窄屏隐藏目录 */
@media (max-width: 900px) {
  .toc {
    display: none;
  }
}

.hint {
  color: var(--text-hint);
  text-align: center;
  padding: 2rem 0;
}

.hint.error {
  color: var(--error);
}

.title {
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--text-title);
  margin: 0 0 0.75rem;
}

.meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.5rem;
  padding-bottom: 1rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid var(--border);
  font-size: 0.85rem;
}

.date {
  color: var(--text-muted);
}

.category {
  color: var(--category-text);
  background: var(--category-bg);
  padding: 0.1rem 0.5rem;
  border-radius: 3px;
  text-decoration: none;
}

.category:hover {
  background: var(--category-bg-hover);
}

.tag {
  color: var(--text-chip);
  background: var(--bg-chip);
  padding: 0.1rem 0.5rem;
  border-radius: 3px;
  text-decoration: none;
}

.tag:hover {
  background: var(--bg-chip-hover);
}

/* 正文样式：v-html 内容不受 scoped 约束，用 :deep() 命中 */
.body {
  color: var(--text);
  line-height: 1.75;
  /* 长链接 / 长单词允许断行，避免撑破视口 */
  overflow-wrap: break-word;
  word-break: break-word;
}

.body :deep(h1),
.body :deep(h2),
.body :deep(h3),
.body :deep(h4),
.body :deep(h5),
.body :deep(h6) {
  margin: 1.5rem 0 0.75rem;
  line-height: 1.3;
  scroll-margin-top: 76px;
}

.body :deep(p) {
  margin: 0.75rem 0;
}

.body :deep(pre) {
  background: var(--bg-code);
  padding: 1rem;
  border-radius: 6px;
  overflow-x: auto;
}

.body :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.9em;
}

/* hljs 高亮后的代码块：去掉 hljs 主题自带的内边距/背景，沿用本站样式 */
.body :deep(pre code.hljs) {
  display: block;
  padding: 0;
  background: transparent;
  overflow: visible;
}

/* 代码块内部用 flex 布局：左侧行号列 + 右侧代码 */
.body :deep(.code-block pre) {
  display: flex;
}

/* 行号列：不换行、不可选中，随代码等高对齐 */
.body :deep(.code-gutter) {
  flex: 0 0 auto;
  padding-right: 1rem;
  margin-right: 1rem;
  border-right: 1px solid var(--gutter-border);
  color: var(--gutter-text);
  text-align: right;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.9em;
  white-space: pre;
  user-select: none;
  -webkit-user-select: none;
}

.body :deep(.code-block pre code) {
  flex: 1 1 auto;
}

/* 代码块容器：工具栏 + pre，工具栏在右上角展示语言和复制按钮 */
.body :deep(.code-block) {
  position: relative;
  margin: 0.75rem 0;
}

.body :deep(.code-block) pre {
  margin: 0;
  /* 顶部留出工具栏高度，避免代码被遮挡 */
  padding-top: 2.5rem;
}

.body :deep(.code-bar) {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 2rem;
  padding: 0 0.75rem;
  font-size: 0.75rem;
}

.body :deep(.code-lang) {
  color: var(--text-muted);
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.body :deep(.code-copy) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  padding: 0;
  color: var(--text-secondary);
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  cursor: pointer;
  transition: color 0.2s ease, border-color 0.2s ease, background 0.2s ease;
}

.body :deep(.code-copy-icon) {
  width: 0.875rem;
  height: 0.875rem;
  display: block;
}

.body :deep(.code-copy .icon-check) {
  display: none;
}

.body :deep(.code-copy.copied .icon-copy) {
  display: none;
}

.body :deep(.code-copy.copied .icon-check) {
  display: block;
}

.body :deep(.code-copy:hover) {
  color: var(--text-strong);
  border-color: var(--code-copy-hover-border);
  background: var(--code-copy-hover-bg);
}

.body :deep(.code-copy.copied) {
  color: var(--copied);
  border-color: transparent;
  background: transparent;
}

.body :deep(img) {
  max-width: 100%;
}

/* 表格外层可横向滚动，避免宽表格撑破页面 */
.body :deep(table) {
  display: block;
  max-width: 100%;
  overflow-x: auto;
  border-collapse: collapse;
}

.body :deep(th),
.body :deep(td) {
  border: 1px solid var(--border-strong);
  padding: 0.4rem 0.6rem;
  white-space: nowrap;
}
</style>
