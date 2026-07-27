<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import SearchNameList from '../components/SearchNameList.vue'
import SearchArticleList from '../components/SearchArticleList.vue'

// 概览每类最多展示的条数，与后端 service.OverviewSize 一致。
const OVERVIEW_SIZE = 5

const route = useRoute()
const keyword = computed(() => route.query.q || '')

const result = ref(null)
const loading = ref(false)
const error = ref('')

// 四个大类按固定顺序展示：分类 → 标签 → 标题 → 正文，无命中的大类不展示。
const groups = computed(() => {
  if (!result.value) return []
  return [
    { type: 'category', label: '分类', data: result.value.categories },
    { type: 'tag', label: '标签', data: result.value.tags },
    { type: 'title', label: '标题', data: result.value.titles },
    { type: 'content', label: '正文', data: result.value.contents },
  ].filter((g) => g.data.total > 0)
})

const empty = computed(() => groups.value.length === 0)

async function fetchOverview(q) {
  if (!q) {
    result.value = null
    return
  }
  loading.value = true
  error.value = ''
  try {
    const resp = await fetch(`/api/search/overview?keyword=${encodeURIComponent(q)}`)
    const data = await resp.json()
    if (data.code !== 0) {
      throw new Error(data.msg || '请求失败')
    }
    result.value = data.data
  } catch (e) {
    error.value = e.message || '加载失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

watch(keyword, fetchOverview, { immediate: true })
</script>

<template>
  <main class="page">
    <div class="content">
      <p class="search-for">
        搜索<template v-if="keyword">：<span class="keyword">{{ keyword }}</span></template>
      </p>

      <p v-if="!keyword" class="hint">请输入搜索关键词</p>
      <p v-else-if="loading" class="hint">加载中…</p>
      <p v-else-if="error" class="hint error">{{ error }}</p>
      <p v-else-if="empty" class="hint">没有匹配的结果</p>

      <template v-else>
        <section v-for="g in groups" :key="g.type" class="group">
          <div class="group-head">
            <h2 class="group-title">{{ g.label }}（{{ g.data.total }}）</h2>
            <RouterLink
              v-if="g.data.total > OVERVIEW_SIZE"
              class="more"
              :to="{ name: 'search-type', params: { type: g.type }, query: { q: keyword } }"
            >更多</RouterLink>
          </div>

          <SearchNameList
            v-if="g.type === 'category' || g.type === 'tag'"
            :items="g.data.list"
            :keyword="keyword"
            :type="g.type"
          />
          <SearchArticleList v-else :articles="g.data.list" :keyword="keyword" />
        </section>
      </template>
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

.search-for {
  margin: 0 0 1.5rem;
  color: #666;
  font-size: 0.95rem;
}

.keyword {
  color: #333;
  font-weight: 600;
}

.hint {
  color: #808080;
  text-align: center;
  padding: 2rem 0;
}

.hint.error {
  color: #c0392b;
}

.group {
  margin-bottom: 2rem;
}

.group-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 0.9rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #eee;
}

.group-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #333;
}

.more {
  flex: none;
  font-size: 0.85rem;
  color: #666;
  text-decoration: none;
  transition: color 0.2s ease;
}

.more:hover {
  color: #000;
}
</style>
