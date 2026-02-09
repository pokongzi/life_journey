# Life Journey PC 端

基于 [electron-egg](https://github.com/dromara/electron-egg) 的 PC 客户端，架构：**前端 ↔ Go 本地后端 ↔ 远程服务端**。

## 环境要求

- Node.js 18+（推荐 20+，Vite 5 需 Node 18）
- npm / pnpm

## 目录结构

```
front/pc/
├── electron/          # Electron 主进程
├── frontend/          # Vue 3 前端
│   └── src/
│       ├── api/       # API 层（与 Go 后端通信）
│       ├── layout/     # 布局（MainLayout/Sidebar/TopBar）
│       ├── stores/     # Pinia store（auth）
│       ├── views/      # 页面（登录、笔记、待办、工具）
│       └── router/     # 路由与登录守卫
└── public/            # 打包后静态资源
```

## 开发

1. **启动 Go 后端**（端口 13245）：
   ```bash
   cd ../../server && go run main.go
   ```

2. **仅前端开发**（vite 代理 /api 到 Go）：
   ```bash
   cd frontend && npm run dev
   ```
   访问 http://localhost:8080

3. **Electron 完整开发**（含主进程）：
   ```bash
   npm run dev
   ```

## 已实现

- 登录页（邮箱+密码 / 邮箱+验证码，验证码接口待 Go 实现）
- 主布局：侧边栏（笔记、待办、工具）、顶栏（用户、登出）
- 笔记页：列表占位、新建/编辑对话框占位
- 待办页：列表、添加、完成/未完成、删除（需 Go /api/todos）
- 工具页：图片压缩占位（需 Go /api/tools/image/compress）
- 路由守卫：未登录跳转登录页
- API 层：axios + token 携带 + 错误处理

## 待 Go 后端实现

- `POST /api/auth/login`、`POST /api/auth/login-by-code`、`POST /api/auth/register`
- `GET/POST/PUT/DELETE /api/notebooks`
- `GET/POST/PUT/DELETE /api/notes`
- `GET/POST/PUT/DELETE /api/todos`
- `POST /api/tools/image/compress`

Go 后端需监听 `127.0.0.1:13245`，开启 CORS 允许前端源。
