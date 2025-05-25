<template>
  <div class="register-container">
    <el-card class="register-card">
      <div slot="header" class="register-header">
        <h2>新质元启(XZYQ) MES系统 - 用户注册</h2>
      </div>
      <el-form :model="registerForm" :rules="rules" ref="registerForm" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="registerForm.username"
            prefix-icon="el-icon-user"
            placeholder="请输入用户名"
          ></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="registerForm.password"
            prefix-icon="el-icon-lock"
            type="password"
            placeholder="请输入密码"
          ></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="registerForm.email"
            prefix-icon="el-icon-message"
            placeholder="请输入邮箱"
          ></el-input>
        </el-form-item>
        <el-form-item label="姓名" prop="full_name">
          <el-input
            v-model="registerForm.full_name"
            prefix-icon="el-icon-user"
            placeholder="请输入姓名"
          ></el-input>
        </el-form-item>
        <el-form-item label="手机" prop="phone">
          <el-input
            v-model="registerForm.phone"
            prefix-icon="el-icon-phone"
            placeholder="请输入手机号码"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="register-button"
            @click="handleRegister"
          >注册</el-button>
          <el-button class="back-button" @click="goToLogin">返回登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Register',
  data() {
    return {
      registerForm: {
        username: '',
        password: '',
        email: '',
        full_name: '',
        phone: ''
      },
      rules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        full_name: [
          { required: true, message: '请输入姓名', trigger: 'blur' }
        ],
        phone: [
          { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
        ]
      },
      loading: false
    }
  },
  methods: {
    handleRegister() {
      this.$refs.registerForm.validate(async valid => {
        if (valid) {
          this.loading = true
          try {
            const response = await axios.post('http://localhost:8080/api/auth/register', this.registerForm)
            this.$message.success('注册成功')
            this.$router.push('/login')
          } catch (error) {
            this.$message.error(error.response?.data?.error || '注册失败')
          } finally {
            this.loading = false
          }
        }
      })
    },
    goToLogin() {
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.register-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f3f3f3;
}

.register-card {
  width: 500px;
}

.register-header {
  text-align: center;
}

.register-button {
  width: 100%;
  margin-bottom: 10px;
}

.back-button {
  width: 100%;
}
</style>