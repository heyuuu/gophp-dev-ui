import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      // component: () => import('../views/HomeView.vue'),
      redirect: '/test/'
    },
    {
      path: '/test/',
      name: 'test',
      component: () => import('../views/test/Index.vue')
    },
    {
      path: '/test/list',
      name: 'test_list',
      component: () => import('../views/test/List.vue')
    },
    {
      path: '/test/run',
      name: 'test_run',
      component: () => import('../views/test/Run.vue')
    }
  ]
})

export default router
