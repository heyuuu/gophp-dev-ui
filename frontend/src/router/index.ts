import { createRouter, createWebHistory } from 'vue-router'

// menu
type MenuItem = {
  path: string
  name: string
  children?: MenuItem[]
}

export const menuItems: MenuItem[] = [
  { path: '/', name: '首页' },
  {
    path: '/compile',
    name: 'compile',
    children: [
      { path: '/compile/run/code', name: 'Run' },
      { path: '/compile/test', name: 'Test' }
    ]
  },
  {
    path: '/execute',
    name: 'execute',
    children: [
      { path: '/execute/run/code', name: 'Run' },
      { path: '/execute/test', name: 'Test' }
    ]
  }
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
      path: '/:mode/run/code',
      name: 'run',
      component: () => import('../views/run/RunCodeView.vue'),
      props: true
    },
    {
      path: '/:mode/test',
      name: 'test',
      component: () => import('../views/test/TestIndexView.vue'),
      props: true
    },
    {
      path: '/:mode/test/list',
      name: 'test_list',
      component: () => import('../views/test/TestListView.vue'),
      props: true
    },
    {
      path: '/:mode/test/run',
      name: 'test_run',
      component: () => import('../views/test/TestRunView.vue'),
      props: true
    }
  ]
})

export default router
