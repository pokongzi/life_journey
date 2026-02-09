import { createRouter, createWebHashHistory } from 'vue-router'
import routerMap from './routerMap'
import { useAuthStore } from '@/stores/auth'

const Router = createRouter({
  history: createWebHashHistory(),
  routes: routerMap,
})

Router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some((r) => r.meta?.requiresAuth)
  const noAuth = to.meta?.noAuth

  if (noAuth) {
    if (authStore.isLoggedIn && to.path === '/login') {
      next('/')
    } else {
      next()
    }
    return
  }

  if (requiresAuth && !authStore.isLoggedIn) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  next()
})

export default Router
