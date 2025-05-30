<template>
  <div class="user-profile">
    <el-dropdown trigger="click" @command="handleCommand">
      <span class="user-profile-link">
        <el-avatar :size="32" class="mr-2">{{ username.charAt(0).toUpperCase() }}</el-avatar>
        {{ username }}
        <el-icon class="el-icon--right"><arrow-down /></el-icon>
      </span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="profile">个人信息</el-dropdown-item>
          <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 个人信息对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="个人信息"
      width="400px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="所属组织">
          <el-input v-model="form.orgName" disabled />
        </el-form-item>
        <el-form-item label="修改密码">
          <el-button text type="primary" @click="showPasswordForm = true">修改密码</el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSave" :loading="saving">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showPasswordForm"
      title="修改密码"
      width="400px"
      append-to-body
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showPasswordForm = false">取消</el-button>
          <el-button type="primary" @click="handleChangePassword" :loading="changingPassword">
            确认修改
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '用户')
const dialogVisible = ref(false)
const showPasswordForm = ref(false)
const saving = ref(false)
const changingPassword = ref(false)
const formRef = ref()
const passwordFormRef = ref()

const form = ref({
  username: '',
  email: '',
  phone: '',
  orgName: ''
})

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取用户信息
const fetchUserProfile = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('/api/user/profile', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const { username: name, email, phone, org } = response.data
    form.value = {
      username: name,
      email: email || '',
      phone: phone || '',
      orgName: org?.name || '-'
    }
  } catch (error) {
    console.error('Error fetching user profile:', error)
    ElMessage.error('获取用户信息失败')
  }
}

// 保存用户信息
const handleSave = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    saving.value = true
    try {
      const token = localStorage.getItem('token')
      await axios.put('/api/user/profile', {
        username: form.value.username,
        email: form.value.email,
        phone: form.value.phone
      }, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      
      ElMessage.success('保存成功')
      username.value = form.value.username
      localStorage.setItem('username', form.value.username)
      dialogVisible.value = false
    } catch (error) {
      console.error('Error saving profile:', error)
      ElMessage.error(error.response?.data?.error || '保存失败')
    } finally {
      saving.value = false
    }
  })
}

// 修改密码
const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return

    changingPassword.value = true
    try {
      const token = localStorage.getItem('token')
      await axios.put('/api/user/change-password', {
        old_password: passwordForm.value.oldPassword,
        new_password: passwordForm.value.newPassword
      }, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      
      ElMessage.success('密码修改成功')
      showPasswordForm.value = false
      passwordForm.value = {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    } catch (error) {
      console.error('Error changing password:', error)
      if (error.response?.status === 401) {
        ElMessage.error('原密码错误')
      } else {
        ElMessage.error(error.response?.data?.error || '密码修改失败')
      }
    } finally {
      changingPassword.value = false
    }
  })
}

// 处理下拉菜单命令
const handleCommand = (command) => {
  if (command === 'profile') {
    dialogVisible.value = true
    fetchUserProfile()
  } else if (command === 'logout') {
    handleLogout()
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    const token = localStorage.getItem('token')
    await axios.post('/api/logout', {}, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    router.push('/login')
  } catch (error) {
    console.error('Logout error:', error)
    ElMessage.error('退出失败')
  }
}

onMounted(() => {
  fetchUserProfile()
})
</script>

<style scoped>
.user-profile {
  display: inline-block;
}

.user-profile-link {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.mr-2 {
  margin-right: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 