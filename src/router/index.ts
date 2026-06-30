import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ArticlesLayout from '../views/ArticlesLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/articles',
      component: ArticlesLayout,
      children: [
        {
          path: '',
          name: 'articles',
          component: () => import('../views/ArticlesView.vue')
        },
        {
          path: ':id',
          name: 'article',
          component: () => import('../views/ArticleView.vue')
        }
      ]
    }
  ],
})

router.beforeEach((to, from) => {
  if (from.name === 'articles') {
    sessionStorage.setItem('scroll_articles', window.scrollY.toString())
  }
})

export default router