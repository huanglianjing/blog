<script setup>
import { computed } from 'vue'

const props = defineProps({
  // 当前页，从 0 开始。
  page: { type: Number, required: true },
  totalPages: { type: Number, required: true },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['change'])

// 生成要展示的页码序列（1 开始），中间用 '...' 表示缩略。
// 规则：必定包含第一页、最后一页、当前页及其前后页；
// 相邻两个展示页码之间缺少的页数 >= 2 时用一个缩略代替，
// 恰好缺 1 页时把该页补上（连起来）。
const items = computed(() => {
  const total = props.totalPages
  const cur = props.page + 1 // 转成 1 开始便于展示

  // 总页数不多时全部展示，不做缩略。
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => ({ type: 'page', value: i + 1 }))
  }

  // 必定展示的页码集合。
  const shown = new Set()
  for (const p of [1, total, cur - 1, cur, cur + 1]) {
    if (p >= 1 && p <= total) shown.add(p)
  }

  const sorted = [...shown].sort((a, b) => a - b)
  const result = []
  for (let i = 0; i < sorted.length; i++) {
    if (i > 0) {
      const gap = sorted[i] - sorted[i - 1]
      if (gap === 2) {
        // 中间只缺一页，补上连起来。
        result.push({ type: 'page', value: sorted[i - 1] + 1 })
      } else if (gap > 2) {
        result.push({ type: 'ellipsis' })
      }
    }
    result.push({ type: 'page', value: sorted[i] })
  }
  return result
})

function go(value) {
  const target = value - 1 // 转回 0 开始
  if (props.disabled || target === props.page) return
  emit('change', target)
}
</script>

<template>
  <nav v-if="totalPages > 1" class="pagination" aria-label="分页">
    <template v-for="(item, i) in items" :key="i">
      <span v-if="item.type === 'ellipsis'" class="ellipsis">…</span>
      <button
        v-else
        class="page-btn"
        :class="{ active: item.value === page + 1 }"
        :disabled="disabled"
        :aria-current="item.value === page + 1 ? 'page' : undefined"
        @click="go(item.value)"
      >
        {{ item.value }}
      </button>
    </template>
  </nav>
</template>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin-top: 2rem;
}

.page-btn {
  min-width: 2rem;
  padding: 0.35rem 0.6rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  color: #333;
  font-variant-numeric: tabular-nums;
  cursor: pointer;
  transition: border-color 0.2s ease, color 0.2s ease;
}

.page-btn:hover:not(:disabled):not(.active) {
  border-color: #999;
  color: #000;
}

.page-btn.active {
  border-color: #333;
  background: #333;
  color: #fff;
  cursor: default;
}

.page-btn:disabled {
  color: #bbb;
  cursor: not-allowed;
}

.ellipsis {
  padding: 0 0.25rem;
  color: #999;
  user-select: none;
}
</style>
