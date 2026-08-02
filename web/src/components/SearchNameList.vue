<script setup>
import HighlightText from './HighlightText.vue'

// 分类与标签的命中结果结构一致（name + count），共用一套 chip 展示，
// 只有跳转路径前缀不同。
defineProps({
  items: { type: Array, required: true },
  keyword: { type: String, default: '' },
  type: { type: String, required: true }, // category 或 tag
})
</script>

<template>
  <div class="name-list">
    <RouterLink
      v-for="item in items"
      :key="item.name"
      class="name-link"
      :to="`/${type}/${encodeURIComponent(item.name)}`"
    >
      <HighlightText :text="item.name" :keyword="keyword" /><span class="count">({{ item.count }})</span>
    </RouterLink>
  </div>
</template>

<style scoped>
.name-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.name-link {
  display: inline-flex;
  align-items: baseline;
  padding: 0.35rem 0.75rem;
  font-size: 0.95rem;
  color: var(--text);
  background: var(--bg-chip);
  border-radius: 4px;
  text-decoration: none;
  transition: color 0.2s ease, background 0.2s ease;
}

.name-link:hover {
  color: var(--text-strong);
  background: var(--bg-chip-hover);
}

.count {
  margin-left: 0.3rem;
  color: var(--text-muted);
  font-size: 0.85rem;
}
</style>
