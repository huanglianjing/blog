<script setup>
import { computed } from 'vue'

const props = defineProps({
  text: { type: String, default: '' },
  keyword: { type: String, default: '' },
})

// 把文本按关键词切成若干段，命中段标记 hit，由模板渲染成 <mark>。
// 用 indexOf 循环而非正则，免去正则元字符转义；
// 大小写不敏感，与后端 SQLite 的 LIKE 行为保持一致。
const segments = computed(() => {
  const text = props.text || ''
  const keyword = props.keyword || ''
  if (!keyword) return [{ text, hit: false }]

  const haystack = text.toLowerCase()
  const needle = keyword.toLowerCase()
  const result = []
  let from = 0
  let at = haystack.indexOf(needle)
  while (at !== -1) {
    if (at > from) result.push({ text: text.slice(from, at), hit: false })
    result.push({ text: text.slice(at, at + needle.length), hit: true })
    from = at + needle.length
    at = haystack.indexOf(needle, from)
  }
  if (from < text.length) result.push({ text: text.slice(from), hit: false })
  return result
})
</script>

<template>
  <span>
    <template v-for="(seg, i) in segments" :key="i">
      <mark v-if="seg.hit">{{ seg.text }}</mark>
      <template v-else>{{ seg.text }}</template>
    </template>
  </span>
</template>
