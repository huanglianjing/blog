import { ref } from 'vue'

const STORAGE_KEY = 'theme'
// 圆形扩散动画时长；降级路径的整体渐变时长与 style.css 的 .theme-switching 保持一致。
const EXPAND_MS = 500
const FADE_MS = 250

// 初值取 <html data-theme>：index.html 的内联脚本已在首屏渲染前从 localStorage 写好，
// 这里不重复读存储，避免两处逻辑不一致。
const theme = ref(document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light')

let fadeTimer = 0

// 写入 DOM 与 localStorage。隐私模式下 localStorage 可能抛异常，失败不影响当次切换。
function applyTheme(next) {
  theme.value = next
  document.documentElement.dataset.theme = next
  try {
    localStorage.setItem(STORAGE_KEY, next)
  } catch {
    // 存储不可用时仅当次生效，刷新后回落到默认的 light。
  }
}

// 降级路径：整页颜色渐变。切换瞬间挂上过渡 class，结束后移除
// （不做成常驻规则，否则会给所有元素的 hover 意外加上过渡）。
function toggleWithFade(next) {
  const root = document.documentElement
  root.classList.add('theme-switching')
  applyTheme(next)

  window.clearTimeout(fadeTimer)
  fadeTimer = window.setTimeout(() => {
    root.classList.remove('theme-switching')
  }, FADE_MS)
}

/**
 * light ⇄ dark。
 *
 * origin 为点击处的视口坐标 {x, y}，新主题以它为圆心扩散覆盖旧主题。
 * 不支持 View Transition API（或用户要求减弱动效）时降级为整页渐变。
 */
export function toggleTheme(origin) {
  const next = theme.value === 'dark' ? 'light' : 'dark'
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

  // 浏览器不会为 prefers-reduced-motion 自动跳过 view transition，需自己判断。
  if (!document.startViewTransition || reduceMotion || !origin) {
    toggleWithFade(next)
    return
  }

  // 扩散动画期间禁掉降级用的整页 transition：两者叠加会让新主题快照的颜色
  // 在圆形内部继续渐变，看起来像是慢了一拍。
  const root = document.documentElement
  window.clearTimeout(fadeTimer)
  root.classList.remove('theme-switching')
  root.classList.add('theme-expanding')

  const transition = document.startViewTransition(() => applyTheme(next))

  transition.ready
    .then(() => {
      const { x, y } = origin
      // 圆要盖满整个视口，半径取到最远角的距离。
      const radius = Math.hypot(Math.max(x, innerWidth - x), Math.max(y, innerHeight - y))
      root.animate(
        { clipPath: [`circle(0px at ${x}px ${y}px)`, `circle(${radius}px at ${x}px ${y}px)`] },
        {
          duration: EXPAND_MS,
          easing: 'ease-in-out',
          // 只动新主题快照：它盖在旧主题之上，圆内露出新色、圆外仍是旧色。
          pseudoElement: '::view-transition-new(root)',
        },
      )
    })
    .catch(() => {
      // transition 被打断（如连续快速点击）时忽略，finished 仍会兜底清理。
    })

  transition.finished.finally(() => root.classList.remove('theme-expanding'))
}

export { theme }
