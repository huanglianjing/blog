<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import Pagination from '../components/Pagination.vue'
import SearchNameList from '../components/SearchNameList.vue'
import SearchArticleList from '../components/SearchArticleList.vue'

const TYPE_LABELS = {
  category: '分类',
  tag: '标签',
  title: '标题',
  content: '正文',
}

const route = useRoute()
const keyword = computed(() => route.query.q || '')
const type = computed(() => route.params.type)
const label = computed(() => TYPE_LABELS[type.value] || '')
// 分类与标签一次返回全部结果，标题与正文分页展示。
const paged = computed(() => type.value === 'title' || type.value === 'content')

const list = ref([])
const page = ref(0)
const totalPages = ref(0)
const loading = ref(false)
const error = ref('')

async function fetchList(p) {
  if (!keyword.value || !label.value) {
    list.value = []
    return
  }
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams({ keyword: keyword.value, type: type.value })
    if (paged.value) params.set('page', p)
    const resp = await fetch(`/api/search/list?${params}`)
    const data = await resp.json()
    if (data.code !== 0) {
      throw new Error(data.msg || '请求失败')
    }
    list.value = data.data.list
    totalPages.value = data.data.total_pages || 0
    page.value = p
  } catch (e) {
    error.value = e.message || '加载失败'
    list.value = []
  } finally {
    loading.value = false
  }
}

// 关键词或大类变化时都回到第一页重新查询。
watch([keyword, type], () => fetchList(0), { immediate: true })
</script>

<template>
  <main class="page">
    <div class="content">
      <p class="search-for">
        <RouterLink class="back" :to="{ name: 'search', query: { q: keyword } }">← 全部结果</RouterLink>
      </p>

      <p v-if="!label" class="hint">未知的搜索类别</p>
      <template v-else>
        <div class="group-head">
          <h2 class="group-title">{{ label }}匹配：<span class="keyword">{{ keyword }}</span></h2>
        </div>

        <p v-if="loading" class="hint">加载中…</p>
        <p v-else-if="error" class="hint error">{{ error }}</p>
        <p v-else-if="list.length === 0" class="hint">没有匹配的结果</p>

        <template v-else>
          <SearchNameList
            v-if="!paged"
            :items="list"
            :keyword="keyword"
            :type="type"
          />
          <template v-else>
            <SearchArticleList :articles="list" :keyword="keyword" />
            <Pagination
              :page="page"
              :total-pages="totalPages"
              :disabled="loading"
              @change="fetchList"
            />
          </template>
        </template>
      </template>
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

.search-for {
  margin: 0 0 1rem;
  font-size: 0.9rem;
}

.back {
  color: var(--text-secondary);
  text-decoration: none;
  transition: color 0.2s ease;
}

.back:hover {
  color: var(--text-strong);
}

.group-head {
  margin-bottom: 0.9rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border);
}

.group-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.keyword {
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
</style>
