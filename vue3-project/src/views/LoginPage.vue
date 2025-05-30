<template>
  <div class="login-container">
    <el-form :model="loginForm" :rules="rules" ref="loginFormRef" class="login-form">
      <h2 class="title">系统登录</h2>
      <el-form-item prop="username">
        <el-input 
          v-model="loginForm.username" 
          placeholder="用户名"
          @keyup.enter="handleLogin"
        >
          <template #prefix>
            <el-icon><User /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item prop="password">
        <el-input 
          v-model="loginForm.password" 
          type="password" 
          placeholder="密码"
          show-password
          @keyup.enter="handleLogin"
        >
          <template #prefix>
            <el-icon><Lock /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-checkbox v-model="rememberUsername">记住账号</el-checkbox>
      </el-form-item>
      <el-form-item>
        <el-button 
          type="primary" 
          :loading="loading" 
          class="login-button"
          @click="handleLogin"
        >
          登录
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)
const rememberUsername = ref(false)
const loginForm = ref({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 页面加载时检查是否有保存的用户名
onMounted(() => {
  const savedUsername = localStorage.getItem('rememberedUsername')
  if (savedUsername) {
    loginForm.value.username = savedUsername
    rememberUsername.value = true
  }
})

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 清除之前的登录信息
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    
    // 保存或清除记住的用户名
    if (rememberUsername.value) {
      localStorage.setItem('rememberedUsername', loginForm.value.username)
    } else {
      localStorage.removeItem('rememberedUsername')
    }

    const response = await axios.post('/api/login', loginForm.value)
    
    // 确保响应中包含所需的数据
    if (!response.data.token || !response.data.user) {
      throw new Error('Invalid server response')
    }
    
    localStorage.setItem('token', response.data.token)
    localStorage.setItem('username', response.data.user.username)
    
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('Login error:', error)
    // 确保清除任何可能存在的无效token
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    
    if (error.response?.data?.error) {
      ElMessage.error(error.response.data.error)
    } else if (error.response?.status === 401) {
      ElMessage.error('用户名或密码错误')
    } else if (error.response?.status === 403) {
      ElMessage.error('该账号已被禁用')
    } else if (error.message === 'Invalid server response') {
      ElMessage.error('服务器响应异常，请稍后重试')
    } else {
      ElMessage.error('登录失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.login-form {
  width: 350px;
  padding: 35px;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.title {
  margin-bottom: 30px;
  text-align: center;
  color: #333;
}

.login-button {
  width: 100%;
}

:deep(.el-input__wrapper) {
  padding-left: 0;
}

:deep(.el-input__prefix) {
  margin-left: 8px;
  color: #909399;
}
</style>