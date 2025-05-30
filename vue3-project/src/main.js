import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import { ElMessage } from 'element-plus'
import 'element-plus/dist/index.css'
import axios from 'axios'

// 配置axios默认值
axios.defaults.baseURL = ''  // 使用相对路径，让代理接管请求

// 添加请求拦截器
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 添加响应拦截器
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // 清除已过期的token
      localStorage.removeItem('token')
      // 跳转到登录页
      router.push('/login')
      ElMessage.error('登录已过期，请重新登录')
    }
    return Promise.reject(error)
  }
)

const app = createApp(App)

app.use(router)
app.use(ElementPlus)

app.config.globalProperties.$axios = axios

app.mount('#app')
