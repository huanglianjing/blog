<script setup>
// 主题切换图标：太阳 ⇄ 月亮的变形动画。
//
// 用一个进度值 t 驱动整张图（0 = 太阳，1 = 月亮）：
// - 主圆半径 4 → 9；
// - 一个半径 7 的「切口圆」从右上角外侧滑入：用 mask 挖掉主圆描边，
//   同时把切口圆自身的描边裁剪在主圆以内，作为月牙内侧那段弧线；
// - 光线在前段进度里缩短并淡出（反向播放即为月亮变太阳时光线逐渐出现）。
//
// 终态与原先两个静态图标重合：太阳为 r=4 的圆 + 8 根 r 从 8 到 10 的光线；
// 月亮路径 "M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" 的外弧圆心 (12,12)、半径 9，
// 内弧圆心约 (16.8,7.2)、半径 7。
//
// 用 requestAnimationFrame 直接改属性，而不是 CSS 过渡 r / cx / cy：
// 后者在 Safari 上支持不完整，且动画中途反向时 JS 能从当前进度平滑接续。
import { ref, computed, watch, onUnmounted, useId } from 'vue'

const props = defineProps({
  // true 显示月亮（浅色模式下，点击转深色），false 显示太阳
  moon: { type: Boolean, required: true },
})

// 与 theme.js 的圆形扩散动画时长一致
const DURATION = 500
// 光线在进度的前 60% 内完成消失，余下时间留给月牙成形
const RAY_END = 0.6
// 切口圆心沿 45° 方向到中心的距离：起点在太阳状态下与主圆不相交，终点即月牙内弧圆心
const CUT_FROM = 13
const CUT_TO = Math.hypot(16.8 - 12, 7.2 - 12)
const RAY_ANGLES = [0, 45, 90, 135, 180, 225, 270, 315]

const uid = useId()
const maskId = `theme-icon-mask-${uid}`
const clipId = `theme-icon-clip-${uid}`

const t = ref(props.moon ? 1 : 0)
let frame = 0

const easeInOut = (p) => (p < 0.5 ? 4 * p * p * p : 1 - (-2 * p + 2) ** 3 / 2)
const lerp = (a, b, k) => a + (b - a) * k

watch(
  () => props.moon,
  (moon) => {
    cancelAnimationFrame(frame)
    const from = t.value
    const to = moon ? 1 : 0
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      t.value = to
      return
    }
    // 中途反向时只走剩下的距离，时长按比例缩短，速度观感一致
    const duration = DURATION * Math.abs(to - from)
    const start = performance.now()
    const step = (now) => {
      const p = Math.min((now - start) / duration, 1)
      t.value = lerp(from, to, easeInOut(p))
      if (p < 1) frame = requestAnimationFrame(step)
    }
    frame = requestAnimationFrame(step)
  },
)

onUnmounted(() => cancelAnimationFrame(frame))

const radius = computed(() => lerp(4, 9, t.value))

const cut = computed(() => {
  const d = lerp(CUT_FROM, CUT_TO, t.value) / Math.SQRT2
  return { x: 12 + d, y: 12 - d }
})

// 主圆与切口圆的交点，即月牙的两个尖角。描边被 mask / clip 截断处是直角，
// 在交点补一个描边宽度的实心圆点，模拟原图标 stroke-linejoin="round" 的圆角。
const tips = computed(() => {
  const R = radius.value
  const dx = cut.value.x - 12
  const dy = cut.value.y - 12
  const D = Math.hypot(dx, dy)
  if (D >= R + 7 || D <= Math.abs(R - 7)) return []
  const a = (R * R - 49 + D * D) / (2 * D)
  const h = Math.sqrt(R * R - a * a)
  const mx = 12 + (a * dx) / D
  const my = 12 + (a * dy) / D
  return [
    { cx: mx + (h * dy) / D, cy: my - (h * dx) / D },
    { cx: mx - (h * dy) / D, cy: my + (h * dx) / D },
  ]
})

// s：光线消失进度，0 为完整光线，1 为完全消失
const rays = computed(() => {
  const s = Math.min(t.value / RAY_END, 1)
  const inner = 8
  const outer = inner + 2 * (1 - s)
  return {
    opacity: 1 - s,
    lines: RAY_ANGLES.map((deg) => {
      const rad = (deg * Math.PI) / 180
      const cos = Math.cos(rad)
      const sin = Math.sin(rad)
      return {
        x1: 12 + inner * cos,
        y1: 12 + inner * sin,
        x2: 12 + outer * cos,
        y2: 12 + outer * sin,
      }
    }),
  }
})
</script>

<template>
  <svg
    viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"
    fill="none" stroke="currentColor" stroke-width="2"
    stroke-linecap="round" stroke-linejoin="round"
  >
    <defs>
      <!-- 白色区域保留、黑色区域挖掉；显式 stroke="none"，否则会继承外层的描边 -->
      <mask :id="maskId" maskUnits="userSpaceOnUse" x="0" y="0" width="24" height="24">
        <rect width="24" height="24" fill="white" stroke="none" />
        <circle :cx="cut.x" :cy="cut.y" r="7" fill="black" stroke="none" />
      </mask>
      <clipPath :id="clipId">
        <circle cx="12" cy="12" :r="radius" />
      </clipPath>
    </defs>
    <!-- 主圆：被切口圆挖掉的部分不画 -->
    <circle cx="12" cy="12" :r="radius" :mask="`url(#${maskId})`" />
    <!-- 切口圆描边：只保留落在主圆内的一段，即月牙内弧 -->
    <circle :cx="cut.x" :cy="cut.y" r="7" :clip-path="`url(#${clipId})`" />
    <circle v-for="(p, i) in tips" :key="i" v-bind="p" r="1" fill="currentColor" stroke="none" />
    <g v-if="rays.opacity > 0" :opacity="rays.opacity">
      <line v-for="(l, i) in rays.lines" :key="i" v-bind="l" />
    </g>
  </svg>
</template>
