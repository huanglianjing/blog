// mermaid 图表渲染：把正文里 goldmark 输出的 <pre><code class="language-mermaid">
// 代码块换成渲染好的 svg。mermaid 体积较大（约 500KB gzip 后 150KB+），
// 故用动态 import 按需加载——只有文章里真的出现 mermaid 代码块时才拉取这份 chunk。
import { theme } from './theme'

let mermaidPromise = null
// 同一页面内多个图表的 id 要唯一，mermaid 用它生成 svg 内部的 DOM id。
let seq = 0

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
  // 配色是渲染时烧进 svg 的，故每次渲染前按当前主题重设一遍配置。
  mermaid.initialize(configFor(theme.value === 'dark'))

  for (const { target, source } of items) {
    seq += 1
    const id = `mermaid-svg-${seq}`
    try {
      const { svg } = await mermaid.render(id, source)
      const wrapper = document.createElement('div')
      wrapper.className = 'mermaid-block'
      // 记下源码，供主题切换后重新渲染。
      wrapper.dataset.source = source
      wrapper.innerHTML = svg
      target.replaceWith(wrapper)
    } catch (e) {
      console.warn('mermaid 渲染失败，保留代码块展示：', e)
      // mermaid.render 失败时会往 body 里留下一个临时容器，清掉避免页面底部出现残留。
      document.getElementById(`d${id}`)?.remove()
    }
  }
}
