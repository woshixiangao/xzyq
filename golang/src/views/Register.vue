<template>
  <div class="register-container">
    <div class="register-box">
      <h2>注册</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>用户名：</label>
          <input type="text" v-model="username" required>
        </div>
        <div class="form-group">
          <label>密码：</label>
          <input type="password" v-model="password" required>
        </div>
        <div class="form-group">
          <label>确认密码：</label>
          <input type="password" v-model="confirmPassword" required>
        </div>
        <div class="form-actions">
          <button type="submit">注册</button>
          <router-link to="/login">返回登录</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'UserRegister',  // 将 'Register' 改为 'UserRegister'
  data() {
    return {
      username: '',
      password: '',
      confirmPassword: ''
    }
  },
  methods: {
    async handleRegister() {
      if (this.password !== this.confirmPassword) {
        alert('两次输入的密码不一致')
        return
      }
      try {
        const response = await axios.post('/api/register', {
          username: this.username,
          password: this.password
        })
        if (response.data.message === '注册成功') {
          this.$router.push('/login')
        }
      } catch (error) {
        alert(error.response?.data?.error || '注册失败')
      }
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.register-box {
  padding: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
}

input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1rem;
}

button {
  padding: 0.5rem 1rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}

a {
  color: #2196F3;
  text-decoration: none;
}
</style>