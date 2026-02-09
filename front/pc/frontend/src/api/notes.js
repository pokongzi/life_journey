/**
 * 笔记 API（本地 Go 后端）
 */
import request from './request'

export function listNotes(params) {
  return request.get('/notes', { params })
}

export function getNote(id) {
  return request.get(`/notes/${id}`)
}

export function createNote(data) {
  return request.post('/notes', data)
}

export function updateNote(id, data) {
  return request.put(`/notes/${id}`, data)
}

export function deleteNote(id) {
  return request.delete(`/notes/${id}`)
}
