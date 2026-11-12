<script setup>
import { ref, onMounted } from 'vue'

// 404 配图统一放在 src/assets/404/，每次进入随机展示一张；增删图片后需重新构建。
// 均为懒加载：只下载被选中的那张，不把所有图片打进主包。
//
// SVG 以源码内联渲染，图里的 var(--text-hint) 等主题变量才能跟随深浅色切换
// （用 <img> 引用的 SVG 读不到页面的 CSS 变量）。这些文件随仓库构建，不是用户输入，
// 内联没有 XSS 风险；SVG 内的 class / id 需各自加前缀，避免互相或与页面样式冲突。
const svgLoaders = import.meta.glob('../assets/404/*.svg', { query: '?raw', import: 'default' })
// 位图用 <img> 展示，不随主题变色。
const imgLoaders = import.meta.glob('../assets/404/*.{png,jpg,jpeg,gif,webp,avif}', {
  query: '?url',
  import: 'default',
})

const candidates = [
  ...Object.keys(svgLoaders).map((path) => ({ type: 'svg', load: svgLoaders[path] })),
  ...Object.keys(imgLoaders).map((path) => ({ type: 'img', load: imgLoaders[path] })),
]

// 选中的图片：{ type: 'svg', content } 或 { type: 'img', src }；
// 加载中为 null，没有图片或加载失败时退回原来的 404 大字。
const picked = ref(null)
const fallback = ref(candidates.length === 0)

onMounted(async () => {
  if (fallback.value) return
  const c = candidates[Math.floor(Math.random() * candidates.length)]
  try {
    const data = await c.load()
    picked.value = c.type === 'svg' ? { type: 'svg', content: data } : { type: 'img', src: data }
  } catch {
    fallback.value = true
  }
})
</script>

<template>
  <main class="notfound">
    <h1 v-if="fallback" class="code">404</h1>
    <div
      v-else-if="picked?.type === 'svg'"
      class="art"
      role="img"
      aria-label="404 页面不存在"
      v-html="picked.content"
    ></div>
    <img v-else-if="picked?.type === 'img'" class="art" :src="picked.src" alt="404 页面不存在" />
  </main>
</template>

<style scoped>
.notfound {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  flex: 1;
  padding: 2rem 1.5rem;
  background: var(--bg);
}

.code {
  font-size: clamp(3rem, 15vw, 8rem);
  font-weight: 800;
  letter-spacing: 0.05em;
  color: var(--text-hint);
  user-select: none;
}

/* 配图宽度随屏宽收缩，高度不超过半屏，避免横图在矮屏上溢出 */
.art {
  display: block;
  width: min(440px, 100%);
  max-height: 60vh;
  object-fit: contain;
  user-select: none;
}

.art :deep(svg) {
  display: block;
  width: 100%;
  height: auto;
  max-height: 60vh;
}
</style>
