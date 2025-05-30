<template>
  <el-container class="home-page">
    <el-aside width="200px">
      <div class="logo">
        <h2>{{ systemName }}</h2>
      </div>
      <side-menu />
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-right">
          <span>{{ username }}</span>
          <el-button type="text" @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import SideMenu from '../components/SideMenu.vue'
import axios from 'axios'

export default {
  name: 'HomePage',
  components: {
    SideMenu
  },
  setup() {
    const router = useRouter()
    const systemName = ref('系统名称')
    const username = ref(localStorage.getItem('username') || '用户')

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

    return {
      systemName,
      username,
      handleLogout
    }
  }
}
</script>

<style scoped>
.home-page {
  height: 100vh;
}

.logo {
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgb(48, 65, 86);
  color: white;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.el-aside {
  background-color: rgb(48, 65, 86);
  color: white;
}

.el-header {
  background-color: white;
  border-bottom: 1px solid #ddd;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.el-main {
  padding: 0;
  background-color: #f0f2f5;
}
</style>