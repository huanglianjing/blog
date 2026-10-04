// mermaid 图表渲染：把正文里 goldmark 输出的 <pre><code class="language-mermaid">
// 代码块换成渲染好的 svg。mermaid 体积较大（约 500KB gzip 后 150KB+），
// 故用动态 import 按需加载——只有文章里真的出现 mermaid 代码块时才拉取这份 chunk。
import { theme } from './theme'

let mermaidPromise = null
// 同一页面内多个图表的 id 要唯一，mermaid 用它生成 svg 内部的 DOM id。
let seq = 0
// 每次调用 renderMermaid 自增；连续切换主题时旧的那次渲染作废，避免两轮渲染交错替换 DOM。
let generation = 0

// 主题切换时新旧图表之间的颜色过渡时长，与 theme.js 的圆形扩散动画一致。
const MORPH_MS = 500
// 参与过渡的样式：线条（stroke）、填充（fill，含箭头 / 圆点等 marker）、
// 渐变色标、foreignObject 里 html 标签的文字与背景色。
const MORPH_PROPS = ['fill', 'stroke', 'color', 'backgroundColor', 'stopColor', 'opacity']

// 按当前主题拼出 mermaid 配置。
function configFor(dark) {
  return {
    // 关掉自动扫描，全部由本模块显式渲染，避免与 Vue 的 DOM 更新抢时序。
    startOnLoad: false,
    // 文章内容由自己维护，可信；strict 会把中文标签里的部分字符转义掉。
    securityLevel: 'loose',
    // 内置主题的明暗是写死的，故随站点主题切换：neutral 灰白、dark 深灰底。
    theme: dark ? 'dark' : 'neutral',
    // 手绘风：节点带斜线纹理、边框与连线有抖动。
    look: 'handDrawn',
    // 手绘抖动来自 roughjs 的随机数，默认值 0 表示每次渲染都换一个随机种子，
    // 于是切换主题重绘后同一张图的笔触会变。给个非 0 定值使其稳定（取值本身无意义）。
    // 注意该项在配置顶层，不在 themeVariables 里。
    handDrawnSeed: 1,
    themeVariables: {
      fontFamily: 'inherit',
    },
  }
}

// 记录 svg 内每个元素当前的计算样式。若上一轮过渡还在进行，取到的是动画中途值，
// 下一轮从这里接着变，连续快速切换也不会跳色。
function snapshot(svg) {
  return [svg, ...svg.querySelectorAll('*')].map((el) => {
    const cs = getComputedStyle(el)
    const style = {}
    for (const p of MORPH_PROPS) style[p] = cs[p]
    return { tag: el.tagName, style }
  })
}

/**
 * 新 svg 已挂进 DOM 后，让每个元素从旧样式过渡到新样式。
 *
 * 同一份源码 + 固定的 handDrawnSeed，换主题前后 svg 结构与几何完全一致、只有配色不同，
 * 按文档顺序配对即可。结构对不上（mermaid 版本差异等）时放弃逐元素过渡，退化为整体淡入。
 */
function morph(before, svg) {
  const els = [svg, ...svg.querySelectorAll('*')]
  if (els.length !== before.length || els.some((el, i) => el.tagName !== before[i].tag)) {
    svg.animate({ opacity: [0, 1] }, { duration: MORPH_MS, easing: 'ease-in-out' })
    return
  }
  els.forEach((el, i) => {
    const old = before[i]
    const cs = getComputedStyle(el)
    const frames = {}
    for (const p of MORPH_PROPS) {
      if (old.style[p] !== cs[p]) frames[p] = [old.style[p], cs[p]]
    }
    // 不留 fill-mode：动画结束后自然回落到新 svg 自身的样式。
    if (Object.keys(frames).length) {
      el.animate(frames, { duration: MORPH_MS, easing: 'ease-in-out' })
    }
  })
}

// 首次调用时加载 mermaid，后续复用同一个 Promise（配置在每次渲染前重设）。
function loadMermaid() {
  if (!mermaidPromise) {
    mermaidPromise = import('mermaid').then(({ default: mermaid }) => mermaid)
  }
  return mermaidPromise
}

/**
 * 渲染 root 内所有 mermaid 代码块。
 *
 * 保留原始源码在 data-source 上，切换主题时可以用同一份源码重新渲染。
 * 渲染失败的块保持原样（退回代码块展示），只在控制台留一条警告。
 */
export async function renderMermaid(root) {
  if (!root) return
  const gen = ++generation

  // 已渲染过的块带 .mermaid-block 容器，重新渲染时直接复用。
  const pending = [...root.querySelectorAll('code.language-mermaid')].map((code) => ({
    // pre 是要被替换掉的节点；code 里的文本即图表源码。
    target: code.closest('pre') ?? code,
    source: code.textContent,
  }))
  const rerender = [...root.querySelectorAll('.mermaid-block[data-source]')].map((el) => ({
    target: el,
    source: el.dataset.source,
  }))

  const items = [...pending, ...rerender]
  if (!items.length) return

  const mermaid = await loadMermaid()
  if (gen !== generation) return
  // 配色是渲染时烧进 svg 的，故每次渲染前按当前主题重设一遍配置。
  mermaid.initialize(configFor(theme.value === 'dark'))

  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

  for (const { target, source } of items) {
    seq += 1
    const id = `mermaid-svg-${seq}`
    try {
      const { svg } = await mermaid.render(id, source)
      // 渲染期间又切了一次主题：交给更新的那轮处理，这里不再动 DOM。
      if (gen !== generation) return
      // 重新渲染的块先记下旧图样式，替换后用来做过渡。
      const oldSvg = target.classList.contains('mermaid-block') && target.querySelector('svg')
      const before = oldSvg && !reduceMotion ? snapshot(oldSvg) : null
      const wrapper = document.createElement('div')
      wrapper.className = 'mermaid-block'
      // 记下源码，供主题切换后重新渲染。
      wrapper.dataset.source = source
      wrapper.innerHTML = svg
      target.replaceWith(wrapper)
      const newSvg = wrapper.querySelector('svg')
      if (before && newSvg && wrapper.isConnected) morph(before, newSvg)
    } catch (e) {
      console.warn('mermaid 渲染失败，保留代码块展示：', e)
      // mermaid.render 失败时会往 body 里留下一个临时容器，清掉避免页面底部出现残留。
      document.getElementById(`d${id}`)?.remove()
    }
  }
}
