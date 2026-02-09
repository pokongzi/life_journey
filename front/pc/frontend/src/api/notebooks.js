/**
 * 笔记本 API（本地 Go 后端）
 */
import request from './request'

export function listNotebooks() {
  return request.get('/notebooks')
}

export function createNotebook(data) {
  return request.post('/notebooks', data)
}

export function updateNotebook(id, data) {
  return request.put(`/notebooks/${id}`, data)
}

export function deleteNotebook(id) {
  return request.delete(`/notebooks/${id}`)
}
