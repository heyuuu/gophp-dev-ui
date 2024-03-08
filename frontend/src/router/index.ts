import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue')
    },
    {
      path:'/tests/',
      name: 'tests',
      component: () => import('../views/tests/Index.vue')
    },
    {
      path:'/tests/list',
      name: 'tests_list',
      component: () => import('../views/tests/List.vue')
    },
    {
      path:'/tests/run',
      name: 'tests_run',
      component: () => import('../views/tests/Run.vue')
    },
  ]
})

export default router
