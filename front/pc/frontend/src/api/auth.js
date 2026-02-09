/**
 * 认证相关 API（由 Go 本地后端转发到远程服务端）
 */
import request from './request'

/** 邮箱+密码登录 */
export function loginByPassword(data) {
  return request.post('/auth/login', data)
}

/** 邮箱+验证码登录 */
export function loginByCode(data) {
  return request.post('/auth/login-by-code', data)
}

/** 邮箱注册 */
export function register(data) {
  return request.post('/auth/register', data)
}

/** 获取当前用户信息 */
export function getMe() {
  return request.get('/auth/me')
}

/** 登出 */
export function logout() {
  return request.post('/auth/logout')
}
