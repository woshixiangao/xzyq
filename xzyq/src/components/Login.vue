<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo-container">
        <img src="../assets/logo.svg" alt="新质元启" class="logo">
      </div>
      <h2 class="login-title">新质元启</h2>
      <form @submit.prevent="handleLogin">
        <div class="input-group">
          <input type="text" placeholder="用户名" v-model="username" class="login-input" />
        </div>
        <div class="input-group">
          <input type="password" placeholder="密码" v-model="password" class="login-input" />
        </div>
        <button type="submit" class="login-btn">登 录</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LoginPage',
  data() {
    return {
      username: '',
      password: ''
    };
  },
  methods: {
    handleLogin() {
      // 这里可以添加登录逻辑
      console.log('用户名:', this.username, '密码:', this.password);
    }
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a2a6c, #3a7bd5);
}

.login-card {
  background: rgba(255, 255, 255, 0.9);
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  width: 380px;
  text-align: center;
}

.logo-container {
  margin-bottom: 20px;
}

.logo {
  width: 80px;
  height: 80px;
}

.login-title {
  color: #1a2a6c;
  margin-bottom: 30px;
  font-size: 24px;
  font-weight: 600;
}

.input-group {
  margin-bottom: 20px;
}

.login-input {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 14px;
  transition: all 0.3s;
}

.login-input:focus {
  border-color: #3a7bd5;
  box-shadow: 0 0 0 2px rgba(58, 123, 213, 0.2);
}

.login-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(to right, #1a2a6c, #3a7bd5);
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}
</style>