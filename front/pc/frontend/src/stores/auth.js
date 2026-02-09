/**
 * 认证状态管理
 */
import { defineStore } from 'pinia'
import { loginByPassword, loginByCode, logout as logoutApi, getMe } from '@/api/auth'

const TOKEN_KEY = 'lj_token'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || '',
    user: null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
  },
  actions: {
    async loginByPassword(email, password) {
      const res = await loginByPassword({ email, password })
      this.token = res?.token || res?.data?.token || ''
      this.user = res?.user || res?.data?.user || null
      if (this.token) localStorage.setItem(TOKEN_KEY, this.token)
      return res
    },
    async loginByCode(email, code) {
      const res = await loginByCode({ email, code })
      this.token = res?.token || res?.data?.token || ''
      this.user = res?.user || res?.data?.user || null
      if (this.token) localStorage.setItem(TOKEN_KEY, this.token)
      return res
    },
    async fetchMe() {
      try {
        const res = await getMe()
        this.user = res?.user || res?.data || res || null
        return this.user
      } catch {
        return null
      }
    },
    async logout() {
      try {
        await logoutApi()
      } catch {}
      this.token = ''
      this.user = null
      localStorage.removeItem(TOKEN_KEY)
    },
  },
})
