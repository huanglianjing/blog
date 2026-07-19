<script setup>
import { ref, onMounted } from 'vue'

const categories = ref([])
const loading = ref(false)
const error = ref('')

async function fetchCategories() {
  loading.value = true
  error.value = ''
  try {
    const resp = await fetch('/category/overview')
    const data = await resp.json()
    if (data.code !== 0) {
      throw new Error(data.msg || '请求失败')
    }
    categories.value = data.data.list
  } catch (e) {
    error.value = e.message || '加载失败'
    categories.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchCategories)
</script>

<template>
  <main class="page">
    <div class="content">
      <p v-if="loading" class="hint">加载中…</p>
      <p v-else-if="error" class="hint error">{{ error }}</p>
      <p v-else-if="categories.length === 0" class="hint">暂无分类</p>

      <ul v-else class="category-list">
        <li v-for="c in categories" :key="c.name" class="category-item">
          <RouterLink
            class="category-link"
            :to="`/category/${encodeURIComponent(c.name)}`"
          >
            {{ c.name }}<span class="count">({{ c.count }})</span>
          </RouterLink>
        </li>
      </ul>
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

.category-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.category-item {
  padding: 0.9rem 0.25rem;
  border-bottom: 1px solid #eee;
}

.category-link {
  font-size: 1.05rem;
  color: #333;
  text-decoration: none;
  transition: color 0.2s ease;
}

.category-link:hover {
  color: #000;
}

.count {
  margin-left: 0.35rem;
  color: #999;
  font-size: 0.9rem;
}
</style>
