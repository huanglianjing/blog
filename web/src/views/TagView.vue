<script setup>
import { ref, onMounted } from 'vue'

const tags = ref([])
const loading = ref(false)
const error = ref('')

// 按文章数分级，热门标签字号更大
function tagLevel(count) {
  if (count >= 40) return 'level-4'
  if (count >= 20) return 'level-3'
  if (count >= 10) return 'level-2'
  if (count >= 5) return 'level-1'
  return 'level-0'
}

async function fetchTags() {
  loading.value = true
  error.value = ''
  try {
    const resp = await fetch('/api/tag/overview')
    const data = await resp.json()
    if (data.code !== 0) {
      throw new Error(data.msg || '请求失败')
    }
    tags.value = data.data.list
  } catch (e) {
    error.value = e.message || '加载失败'
    tags.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchTags)
</script>

<template>
  <main class="page">
    <div class="content">
      <p v-if="loading" class="hint">加载中…</p>
      <p v-else-if="error" class="hint error">{{ error }}</p>
      <p v-else-if="tags.length === 0" class="hint">暂无标签</p>

      <div v-else class="tag-list">
        <RouterLink
          v-for="t in tags"
          :key="t.name"
          class="tag-link"
          :class="tagLevel(t.count)"
          :to="`/tag/${encodeURIComponent(t.name)}`"
        >
          {{ t.name }}<span class="count">({{ t.count }})</span>
        </RouterLink>
      </div>
    </div>
  </main>
</template>

<style scoped>
.page {
  flex: 1;
  width: 100%;
  background: var(--bg);
}

.content {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.hint {
  color: var(--text-hint);
  text-align: center;
  padding: 2rem 0;
}

.hint.error {
  color: var(--error);
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 0.75rem;
}

.tag-link {
  display: inline-flex;
  align-items: baseline;
  padding: 0.35em 0.75em;
  font-size: 0.95rem;
  color: var(--text);
  background: var(--bg-chip);
  border-radius: 4px;
  text-decoration: none;
  transition: color 0.2s ease, background 0.2s ease;
}

/* 文章数越多字号越大：<5 / 5+ / 10+ / 20+ / 40+ */
.tag-link.level-0 {
  font-size: 1rem;
}

.tag-link.level-1 {
  font-size: 1.4rem;
}

.tag-link.level-2 {
  font-size: 1.6rem;
}

.tag-link.level-3 {
  font-size: 1.7rem;
}

.tag-link.level-4 {
  font-size: 2.1rem;
}

.tag-link:hover {
  color: var(--text-strong);
  background: var(--bg-chip-hover);
}

.count {
  margin-left: 0.3em;
  color: var(--text-muted);
  font-size: 0.85em;
}

@media (max-width: 600px) {
  .tag-link.level-2 {
    font-size: 1.2rem;
  }

  .tag-link.level-3 {
    font-size: 1.35rem;
  }

  .tag-link.level-4 {
    font-size: 1.55rem;
  }
}
</style>
