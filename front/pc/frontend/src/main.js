import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import App from './App.vue'
import './assets/global.less'
import components from './components/global'
import Router from './router/index'

const app = createApp(App)
const pinia = createPinia()

for (const i in components) {
  app.component(i, components[i])
}

app.use(pinia).use(Router).use(ElementPlus, { locale: zhCn }).mount('#app')
