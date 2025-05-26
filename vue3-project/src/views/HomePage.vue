<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="logo">系统名称</div>
      <div class="user-info">
        <el-dropdown>
          <span class="el-dropdown-link">
            <el-avatar :size="32" :src="user.avatar" />
            {{ user.name }}
            <el-icon class="el-icon--right">
              <arrow-down />
            </el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="editProfile">修改资料</el-dropdown-item>
              <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    
    <el-container>
      <el-aside width="200px" class="aside">
        <el-menu
          default-active="1"
          class="el-menu-vertical"
          background-color="#545c64"
          text-color="#fff"
          active-text-color="#ffd04b">
          <el-menu-item index="/users" @click="router.push('/users')">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/organizations" @click="router.push('/organizations')">
            <el-icon><User /></el-icon>
            <span>组织管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>

    <!-- 修改资料对话框 -->
    <el-dialog v-model="dialogVisible" title="修改资料" width="30%">
      <el-form :model="profileForm" label-width="80px" :rules="rules" ref="formRef">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="profileForm.username"></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="profileForm.password" type="password" placeholder="不修改请留空"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="profileForm.email"></el-input>
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="profileForm.phone"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveProfile">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </el-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import {
  ArrowDown,
  User
} from '@element-plus/icons-vue'

const router = useRouter()
const user = ref({
  name: '管理员',
  avatar: ''
})

// 对话框显示状态
const dialogVisible = ref(false)

// 表单数据
const profileForm = ref({
  username: user.value.name,
  password: '',
  email: '',
  phone: ''
})

// 表单验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { min: 6, message: '密码长度不能小于6个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

// 表单引用
const formRef = ref()

// 获取用户信息
const getUserInfo = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('http://localhost:8080/api/user/profile', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    user.value.name = response.data.username
    profileForm.value = response.data
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

// 页面加载时获取用户信息
onMounted(() => {
  getUserInfo()
})

// 打开修改资料对话框
const editProfile = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('http://localhost:8080/api/user/profile', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    profileForm.value = response.data
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '获取用户信息失败')
  }
}

// 保存修改后的资料
const saveProfile = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    const token = localStorage.getItem('token')
    const response = await axios.put('http://localhost:8080/api/user/profile', profileForm.value, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    ElMessage.success('保存成功')
    dialogVisible.value = false
    // 更新显示的用户名
    user.value.name = response.data.username
    // 清空密码字段
    profileForm.value.password = ''
  } catch (error) {
    if (error.response?.data?.error) {
      ElMessage.error(error.response.data.error)
    } else if (error.message) {
      ElMessage.error(error.message)
    } else {
      ElMessage.error('保存失败')
    }
  }
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #409EFF;
  color: #fff;
}

.logo {
  font-size: 20px;
  font-weight: bold;
}

.user-info {
  cursor: pointer;
}

.aside {
  background-color: #545c64;
}

.el-menu-vertical {
  border-right: none;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.main {
  padding: 20px;
  background-color: #f0f2f5;
}

.el-dropdown-link {
  display: flex;
  align-items: center;
  color: #fff;
}

.el-dropdown-link .el-avatar {
  margin-right: 8px;
}
</style>