import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ArticleView from '../views/ArticleView.vue'
import ArticleDetailView from '../views/ArticleDetailView.vue'
import CategoryView from '../views/CategoryView.vue'
import CategoryDetailView from '../views/CategoryDetailView.vue'
import TagView from '../views/TagView.vue'
import TagDetailView from '../views/TagDetailView.vue'
import SearchView from '../views/SearchView.vue'
import SearchTypeView from '../views/SearchTypeView.vue'
import NotFoundView from '../views/NotFoundView.vue'

// 站点名，用作首页标题与兜底标题。
const SITE_NAME = 'Moondo'

// 搜索大类的中文名，与 SearchTypeView 里的 TYPE_LABELS 一致。
const SEARCH_TYPE_LABELS = {
  category: '分类',
  tag: '标签',
  title: '标题',
  content: '正文',
}

// 各页面网页标题：meta.title 为字符串或由 route 计算的函数，
// 统一在 afterEach 里写入 document.title。
const routes = [
  { path: '/', name: 'home', component: HomeView, meta: { title: SITE_NAME } },
  { path: '/article', name: 'article', component: ArticleView, meta: { title: '文章' } },
  {
    path: '/article/:title',
    name: 'article-detail',
    component: ArticleDetailView,
    // 文章页标题就是文章标题（路由参数已由 vue-router 解码）。
    meta: { title: (route) => route.params.title },
  },
  { path: '/category', name: 'category', component: CategoryView, meta: { title: '分类' } },
  {
    path: '/category/:name',
    name: 'category-detail',
    component: CategoryDetailView,
    meta: { title: (route) => `分类 - ${route.params.name}` },
  },
  { path: '/tag', name: 'tag', component: TagView, meta: { title: '标签' } },
  {
    path: '/tag/:name',
    name: 'tag-detail',
    component: TagDetailView,
    meta: { title: (route) => `标签 - ${route.params.name}` },
  },
  {
    path: '/search',
    name: 'search',
    component: SearchView,
    meta: { title: (route) => (route.query.q ? `搜索 - ${route.query.q}` : '搜索') },
  },
  {
    path: '/search/:type',
    name: 'search-type',
    component: SearchTypeView,
    meta: {
      title: (route) => {
        const parts = ['搜索']
        if (route.query.q) parts.push(route.query.q)
        const label = SEARCH_TYPE_LABELS[route.params.type]
        if (label) parts.push(label)
        return parts.join(' - ')
      },
    },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'notfound',
    component: NotFoundView,
    meta: { title: '404' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.afterEach((to) => {
  const title = typeof to.meta.title === 'function' ? to.meta.title(to) : to.meta.title
  document.title = title || SITE_NAME
})

export default router
