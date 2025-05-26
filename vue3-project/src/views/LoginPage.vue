<template>
  <div class="login-container">
    <el-form :model="loginForm" :rules="rules" ref="loginForm" class="login-form">
      <h2 class="title">系统登录</h2>
      <el-form-item prop="username">
        <el-input v-model="loginForm.username" placeholder="用户名"></el-input>
      </el-form-item>
      <el-form-item prop="password">
        <el-input type="password" v-model="loginForm.password" placeholder="密码"></el-input>
      </el-form-item>
      <el-form-item>
        <el-checkbox v-model="rememberUsername">记住账号</el-checkbox>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm" :loading="loading">登录</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const router = useRouter()
const rememberUsername = ref(false)
const loginForm = ref({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const loading = ref(false)

// 页面加载时检查是否有保存的用户名
onMounted(() => {
  const savedUsername = localStorage.getItem('rememberedUsername')
  if (savedUsername) {
    loginForm.value.username = savedUsername
    rememberUsername.value = true
  }
})

const submitForm = () => {
  loading.value = true
  // 保存或清除记住的用户名
  if (rememberUsername.value) {
    localStorage.setItem('rememberedUsername', loginForm.value.username)
  } else {
    localStorage.removeItem('rememberedUsername')
  }

  axios.post('http://localhost:8080/api/login', loginForm.value)
    .then(response => {
      localStorage.setItem('token', response.data.token)
      router.push('/')
    })
    .catch(error => {
      ElMessage.error(error.response?.data?.message || '登录失败')
    })
    .finally(() => {
      loading.value = false
    })
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
  padding: 30px;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.title {
  margin-bottom: 30px;
  text-align: center;
  color: #333;
}
</style>