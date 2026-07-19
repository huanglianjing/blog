<script setup>
import { ref, onMounted } from 'vue'

const tags = ref([])
const loading = ref(false)
const error = ref('')

async function fetchTags() {
  loading.value = true
  error.value = ''
  try {
    const resp = await fetch('/tag/overview')
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
  background: #ffffff;
}

.content {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.hint {
  color: #808080;
  text-align: center;
  padding: 2rem 0;
}

.hint.error {
  color: #c0392b;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.tag-link {
  display: inline-flex;
  align-items: baseline;
  padding: 0.35rem 0.75rem;
  font-size: 0.95rem;
  color: #333;
  background: #f0f0f0;
  border-radius: 4px;
  text-decoration: none;
  transition: color 0.2s ease, background 0.2s ease;
}

.tag-link:hover {
  color: #000;
  background: #e4e4e4;
}

.count {
  margin-left: 0.3rem;
  color: #999;
  font-size: 0.85rem;
}
</style>
