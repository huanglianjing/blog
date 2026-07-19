import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ArticleView from '../views/ArticleView.vue'
import ArticleDetailView from '../views/ArticleDetailView.vue'
import CategoryView from '../views/CategoryView.vue'
import CategoryDetailView from '../views/CategoryDetailView.vue'
import TagView from '../views/TagView.vue'
import TagDetailView from '../views/TagDetailView.vue'
import NotFoundView from '../views/NotFoundView.vue'

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/article', name: 'article', component: ArticleView },
  { path: '/article/:title', name: 'article-detail', component: ArticleDetailView },
  { path: '/category', name: 'category', component: CategoryView },
  { path: '/category/:name', name: 'category-detail', component: CategoryDetailView },
  { path: '/tag', name: 'tag', component: TagView },
  { path: '/tag/:name', name: 'tag-detail', component: TagDetailView },
  { path: '/:pathMatch(.*)*', name: 'notfound', component: NotFoundView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
