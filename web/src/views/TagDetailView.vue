<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Pagination from '../components/Pagination.vue'
import ArticleList from '../components/ArticleList.vue'

const route = useRoute()

const articles = ref([])
const page = ref(0)
const totalPages = ref(0)
const loading = ref(false)
const error = ref('')

async function fetchArticles(name, p) {
  loading.value = true
  error.value = ''
  try {
    const resp = await fetch(
      `/api/tag/list?name=${encodeURIComponent(name)}&page=${p}`,
    )
    const data = await resp.json()
    if (data.code !== 0) {
      throw new Error(data.msg || '请求失败')
    }
    articles.value = data.data.list
    totalPages.value = data.data.total_pages
    page.value = p
  } catch (e) {
    error.value = e.message || '加载失败'
    articles.value = []
  } finally {
    loading.value = false
  }
}

function goPage(p) {
  fetchArticles(route.params.name, p)
}

// 路由参数变化时（含首次进入）重新加载。
watch(
  () => route.params.name,
  (name) => {
    if (name) fetchArticles(name, 0)
  },
  { immediate: true },
)
</script>

<template>
  <main class="page">
    <div class="content">
      <h1 class="tag-title">{{ route.params.name }}</h1>

      <p v-if="loading" class="hint">加载中…</p>
      <p v-else-if="error" class="hint error">{{ error }}</p>
      <p v-else-if="articles.length === 0" class="hint">暂无文章</p>

      <ArticleList v-else :articles="articles" />

      <Pagination
        :page="page"
        :total-pages="totalPages"
        :disabled="loading"
        @change="goPage"
      />
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

.tag-title {
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--text-title);
  margin: 0 0 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border);
}

.hint {
  color: var(--text-hint);
  text-align: center;
  padding: 2rem 0;
}

.hint.error {
  color: var(--error);
}

</style>
