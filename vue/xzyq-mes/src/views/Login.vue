<template>
  <div class="login-container">
    <el-card class="login-card">
      <div slot="header" class="login-header">
        <h2>新质元启(XZYQ) MES系统</h2>
      </div>
      <el-form :model="loginForm" :rules="rules" ref="loginForm" label-width="0px">
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            prefix-icon="el-icon-user"
            placeholder="用户名"
          ></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            prefix-icon="el-icon-lock"
            type="password"
            placeholder="密码"
            @keyup.enter.native="handleLogin"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >登录</el-button>
        </el-form-item>
        <el-form-item class="register-link">
          <el-link type="primary" @click="goToRegister">还没有账号？点击注册</el-link>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Login',
  data() {
    return {
      loginForm: {
        username: '',
        password: ''
      },
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      },
      loading: false
    }
  },
  methods: {
    handleLogin() {
      this.$refs.loginForm.validate(async valid => {
        if (valid) {
          this.loading = true
          try {
            await this.$store.dispatch('login', this.loginForm)
            this.$message.success('登录成功')
            this.$router.push('/dashboard')
          } catch (error) {
            this.$message.error(error.response?.data?.error || '登录失败')
          } finally {
            this.loading = false
          }
        }
      })
    },
    goToRegister() {
      this.$router.push('/register')
    }
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f3f3f3;
}

.login-card {
  width: 400px;
}

.login-header {
  text-align: center;
}

.login-button {
  width: 100%;
}

.register-link {
  text-align: center;
  margin-top: 10px;
}
</style>