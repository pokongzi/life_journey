/**
 * axios 请求封装，统一请求本地 Go 后端
 * 前端仅与本地 Go 通信，不直连远程服务端
 */
import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

// 请求拦截：携带 token
request.interceptors.request.use((config) => {
  const token = localStorage.getItem('lj_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// 响应拦截
request.interceptors.response.use(
  (res) => res.data,
  (err) => {
    const msg = err.response?.data?.message || err.message || '请求失败'
    return Promise.reject(new Error(msg))
  }
)

export default request
