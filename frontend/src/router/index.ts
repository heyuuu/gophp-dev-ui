import { createRouter, createWebHistory } from 'vue-router'

// menu
type MenuItem = {
  path: string
  name: string
  children?: MenuItem[]
}

export const menuItems: MenuItem[] = [
  { path: '/', name: '首页' },
  { path: '/run', name: '运行代码' },
  { path: '/test', name: '运行Case' }
]

// router
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue')
    },
    {
      path: '/run/',
      name: 'run',
      component: () => import('../views/run/RunView.vue')
    },
    {
      path: '/test/',
      name: 'test',
      component: () => import('../views/test/IndexView.vue')
    },
    {
      path: '/test/list',
      name: 'test_list',
      component: () => import('../views/test/ListView.vue')
    },
    {
      path: '/test/run',
      name: 'test_run',
      component: () => import('../views/test/RunView.vue')
    }
  ]
})

export default router
