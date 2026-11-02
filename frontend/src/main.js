import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './tv-theme.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { clearSession, getToken } from './auth'

const app = createApp(App)

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 配置axios
axios.defaults.baseURL = '/api'
axios.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
axios.interceptors.response.use((response) => {
  const code = response?.data?.code
  const url = response.config?.url || ''
  const isLogin = url.includes('/admin/login')
  if (code === 401 && !isLogin) {
    clearSession()
    window.dispatchEvent(new Event('admin-auth-expired'))
  } else if (code === 403 && !isLogin) {
    ElMessage.error(response.data.message || '没有权限')
  }
  return response
})
app.config.globalProperties.$axios = axios

app.use(ElementPlus, { locale: zhCn })
app.mount('#app')
