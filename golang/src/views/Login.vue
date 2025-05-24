<template>
  <div class="login-container">
    <el-card class="login-box">
      <h2 class="text-center mb-4">登录</h2>
      <el-form @submit.native.prevent="handleLogin">
        <el-form-item label="用户名：">
          <el-input v-model="username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码：">
          <el-input type="password" v-model="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <div class="d-flex justify-content-between align-items-center">
          <el-button type="primary" native-type="submit">登录</el-button>
          <router-link to="/register" class="text-primary">注册账号</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'UserLogin',
  data() {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    async handleLogin() {
      try {
        const response = await this.$axios.post('/api/login', {
          username: this.username,
          password: this.password
        })
        if (response.data.message === '登录成功') {
          // 保存用户名到 sessionStorage
          sessionStorage.setItem('username', this.username)
          this.$router.push('/home')
        }
      } catch (error) {
        this.$message.error(error.response?.data?.error || '登录失败')
      }
    }
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

.login-box {
  width: 100%;
  max-width: 400px;
}
</style>