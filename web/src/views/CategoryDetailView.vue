<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Pagination from '../components/Pagination.vue'

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
      `/category/list?name=${encodeURIComponent(name)}&page=${p}`,
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
      <h1 class="category-title">{{ route.params.name }}</h1>

      <p v-if="loading" class="hint">加载中…</p>
      <p v-else-if="error" class="hint error">{{ error }}</p>
      <p v-else-if="articles.length === 0" class="hint">暂无文章</p>

      <ul v-else class="article-list">
        <li v-for="article in articles" :key="article.id" class="article-item">
          <RouterLink
            class="article-link"
            :to="`/article/${encodeURIComponent(article.title)}`"
          >
            {{ article.title }}
          </RouterLink>
          <div class="article-meta">
            <span class="date">{{ article.date }}</span>
            <RouterLink
              v-if="article.category_name"
              class="category"
              :to="`/category/${encodeURIComponent(article.category_name)}`"
            >{{ article.category_name }}</RouterLink>
            <RouterLink
              v-for="tag in article.tags"
              :key="tag"
              class="tag"
              :to="`/tag/${encodeURIComponent(tag)}`"
            >{{ tag }}</RouterLink>
          </div>
          <p v-if="article.summary" class="article-summary">{{ article.summary }}</p>
        </li>
      </ul>

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
  background: #ffffff;
}

.content {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.category-title {
  font-size: 1.6rem;
  font-weight: 700;
  color: #222;
  margin: 0 0 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}

.hint {
  color: #808080;
  text-align: center;
  padding: 2rem 0;
}

.hint.error {
  color: #c0392b;
}

.article-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.article-item {
  padding: 1.2rem 0.25rem;
  border-bottom: 1px solid #eee;
}

.article-link {
  display: inline-block;
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  text-decoration: none;
  transition: color 0.2s ease;
}

.article-link:hover {
  color: #000;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.5rem;
  font-size: 0.8rem;
}

.article-meta .date {
  color: #999;
}

.article-meta .category {
  color: #fff;
  background: #666;
  padding: 0.1rem 0.5rem;
  border-radius: 3px;
  text-decoration: none;
}

.article-meta .category:hover {
  background: #444;
}

.article-meta .tag {
  color: #555;
  background: #f0f0f0;
  padding: 0.1rem 0.5rem;
  border-radius: 3px;
  text-decoration: none;
}

.article-meta .tag:hover {
  background: #e4e4e4;
}

.article-summary {
  margin: 0.6rem 0 0;
  color: #666;
  font-size: 0.9rem;
  line-height: 1.6;
  /* 最多展示四行，超出省略 */
  display: -webkit-box;
  -webkit-line-clamp: 4;
  line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

</style>
