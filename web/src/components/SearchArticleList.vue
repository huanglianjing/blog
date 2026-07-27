<script setup>
import HighlightText from './HighlightText.vue'

defineProps({
  articles: { type: Array, required: true },
  keyword: { type: String, default: '' },
})
</script>

<template>
  <ul class="article-list">
    <li v-for="article in articles" :key="article.id" class="article-item">
      <RouterLink
        class="article-link"
        :to="`/article/${encodeURIComponent(article.title)}`"
      >
        <HighlightText :text="article.title" :keyword="keyword" />
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
      <!-- 正文命中时展示带高亮的片段，仅标题命中时回落到文章摘要。 -->
      <p v-if="article.snippet" class="article-summary">
        <HighlightText :text="article.snippet" :keyword="keyword" />
      </p>
      <p v-else-if="article.summary" class="article-summary">{{ article.summary }}</p>
    </li>
  </ul>
</template>

<style scoped>
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
