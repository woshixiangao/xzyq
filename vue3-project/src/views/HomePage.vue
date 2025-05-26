<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="logo">系统名称</div>
      <div class="user-info">
        <el-dropdown>
          <span class="el-dropdown-link">
            <el-avatar :size="small" :src="user.avatar" />
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
          <el-menu-item index="1">
            <el-icon><icon-menu /></el-icon>
            <span>导航一</span>
          </el-menu-item>
          <el-menu-item index="2">
            <el-icon><icon-menu /></el-icon>
            <span>导航二</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  Menu as IconMenu,
  ArrowDown
} from '@element-plus/icons-vue'

const router = useRouter()
const user = ref({
  name: '管理员',
  avatar: ''
})

const editProfile = () => {
  // 修改资料逻辑
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