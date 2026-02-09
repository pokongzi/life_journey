/**
 * 路由配置
 * 未登录时跳转登录页；登录后进入主布局
 */
const constantRouterMap = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录', noAuth: true },
  },
  {
    path: '/',
    component: () => import('@/layout/MainLayout.vue'),
    redirect: '/notes',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'notes',
        name: 'Notes',
        component: () => import('@/views/notes/Index.vue'),
        meta: { title: '笔记', icon: 'Document' },
      },
      {
        path: 'todos',
        name: 'Todos',
        component: () => import('@/views/todos/Index.vue'),
        meta: { title: '待办', icon: 'List' },
      },
      {
        path: 'tools',
        name: 'Tools',
        component: () => import('@/views/tools/Index.vue'),
        meta: { title: '工具', icon: 'Tools' },
      },
    ],
  },
]

export default constantRouterMap
