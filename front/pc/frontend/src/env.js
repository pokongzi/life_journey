/**
 * 环境变量
 * 开发时：前端请求 /api 由 vite proxy 转发到 Go 后端 127.0.0.1:13245
 * 打包后：Electron 加载打包后的静态文件，/api 需由主进程代理或 Go 同端口提供
 */
export const API_BASE = import.meta.env.VITE_API_BASE || ''
export const IS_DEV = import.meta.env.DEV
