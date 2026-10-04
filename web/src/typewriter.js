/**
 * 面包屑打字机效果的纯逻辑：把层级数组摊平成字符序列，并算出两个序列的公共前缀。
 * 与 Vue 无关，便于单独理解与复用。
 */

/*
 * 动画速度由下面四个常量和删除阶段的倍率控制。
 *
 * 删除与打字各自的总时长预算（ms）。按预算分摊到每个字符，而不是固定每字延时：
 * 文章名可能有三四十个字符，固定每字延时会让顶栏打到一秒多以后还没停。
 */
const ERASE_BUDGET = 200
const TYPE_BUDGET = 600
/*
 * 每字延时的上下限，用来兜住极端字数。
 * 注意：短面包屑（如翻页只删打 1 个字符）由 MAX_STEP 决定速度，预算根本用不满，
 * 所以调整整体速度时这两个值要跟着一起改，否则短的那几种不会有变化。
 */
const MIN_STEP = 18
const MAX_STEP = 80

/**
 * 把 `[{ text, to? }]` 摊平成逐字符数组。
 *
 * 每个字符记住它属于第几级（`group`）以及那一级的 `to`，模板据此重新分组渲染，
 * 保证动画中途的样式（链接 / 当前页）与终态一致。分隔符 `/` 也占一个字符，
 * 这样回退时会先吃掉分隔符，观感上和真的按 backspace 一样。
 */
export function flattenCrumbs(crumbs) {
  const chars = []
  crumbs.forEach((crumb, i) => {
    // 分隔符自成一组，避免它跟着上一级的链接样式走。
    chars.push({ char: '/', group: `sep-${i}`, sep: true, to: null })
    for (const char of String(crumb.text)) {
      chars.push({ char, group: `crumb-${i}`, sep: false, to: crumb.to ?? null })
    }
  })
  return chars
}

/**
 * 两个字符序列的公共前缀长度，只比字符本身，不比 `to`。
 *
 * 同一级的文字常常不变、只是从当前页变成了可点的上级（「分类」→「分类 / 数据库」、
 * 搜索关键词进到某一类结果），此时若把 `to` 也算进比较，前缀会退到最开头，
 * 于是把一模一样的字删掉再原样打一遍，白费半秒。
 * 代价是这一段的样式（加粗 ↔ 虚下划线）会在切换的瞬间变一下，比来回重打好得多。
 */
export function commonPrefixLength(a, b) {
  const max = Math.min(a.length, b.length)
  let i = 0
  while (i < max && a[i].char === b[i].char) i++
  return i
}

// 把总预算分摊到 count 个字符上，夹在上下限之间。count 为 0 时不会被用到。
export function stepDelay(count, erasing) {
  if (count <= 0) return erasing ? MIN_STEP / 1.5 : MIN_STEP
  const budget = erasing ? ERASE_BUDGET : TYPE_BUDGET
  const delay = Math.min(MAX_STEP, Math.max(MIN_STEP, Math.round(budget / count)))
  // 删除阶段提速到原来的 1.5 倍；输入阶段沿用原来的延时。
  return erasing ? delay / 1.5 : delay
}

/** 面包屑的无障碍文本，供读屏软件读取（动画中间态不适合朗读）。 */
export function crumbsToText(crumbs) {
  return crumbs.map((c) => c.text).join(' / ')
}
