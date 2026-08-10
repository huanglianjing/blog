<script setup>
import { ref, onMounted } from 'vue'
import ArticleList from '../components/ArticleList.vue'

const articles = ref([])
const loading = ref(false)
const error = ref('')

async function fetchRecent() {
  loading.value = true
  error.value = ''
  try {
    const resp = await fetch('/api/article/recent')
    const data = await resp.json()
    if (data.code !== 0) {
      throw new Error(data.msg || '请求失败')
    }
    articles.value = data.data.list
  } catch (e) {
    error.value = e.message || '加载失败'
    articles.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchRecent)
</script>

<template>
  <main class="page">
    <div class="content">
      <section class="block">
        <div class="block-head">
          <h2 class="block-title">最近文章</h2>
          <RouterLink class="more" to="/article">更多</RouterLink>
        </div>

        <p v-if="loading" class="hint">加载中…</p>
        <p v-else-if="error" class="hint error">{{ error }}</p>
        <p v-else-if="articles.length === 0" class="hint">暂无文章</p>

        <ArticleList v-else :articles="articles" />
      </section>
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

/* 标题与「更多」同一行，两端对齐 */
.block-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 0.9rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border);
}

.block-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text);
}

.hint {
  color: var(--text-hint);
  text-align: center;
  padding: 2rem 0;
}

.hint.error {
  color: var(--error);
}

.more {
  flex: none;
  font-size: 0.9rem;
  color: var(--text-secondary);
  text-decoration: none;
  transition: color 0.2s ease;
}

.more:hover {
  color: var(--text-strong);
}
</style>
