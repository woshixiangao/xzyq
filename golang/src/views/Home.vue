<template>
  <el-container class="home-container">
    <!-- 顶部栏 -->
    <el-header>
      <div class="header-content">
        <div class="logo">后台管理系统</div>
        <div class="user-info">
          <span>{{ username }}</span>
          <el-button type="text" @click="handleLogout">退出登录</el-button>
        </div>
      </div>
    </el-header>
    
    <el-container>
      <!-- 侧边栏 -->
      <el-aside width="200px">
        <el-menu
          :router="true"
          background-color="#304156"
          text-color="#fff"
          active-text-color="#409EFF">
          <el-menu-item index="/home">
            <i class="el-icon-s-home"></i>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/logs">
            <i class="el-icon-document"></i>
            <span>日志管理</span>
          </el-menu-item>
          <el-submenu index="2">
            <template slot="title">
              <i class="el-icon-user"></i>
              <span>用户管理</span>
            </template>
            <el-menu-item index="/users">
              <span>用户列表</span>
            </el-menu-item>
          </el-submenu>
          <el-submenu index="3">
            <template slot="title">
              <i class="el-icon-office-building"></i>
              <span>租户管理</span>
            </template>
            <el-menu-item index="/tenants">
              <span>租户列表</span>
            </el-menu-item>
          </el-submenu>
        </el-menu>
      </el-aside>
      
      <!-- 主要内容区 -->
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
export default {
  name: 'MainHome',
  data() {
    return {
      username: ''
    }
  },
  created() {
    // 获取当前登录用户信息
    this.username = sessionStorage.getItem('username') || ''
  },
  methods: {
    async handleLogout() {
      try {
        await this.$axios.post('/api/logout')
        this.$router.push('/login')
      } catch (error) {
        this.$message.error('退出登录失败')
      }
    }
  }
}
</script>

<style scoped>
.home-container {
  height: 100vh;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.el-aside {
  background-color: #304156;
  color: #fff;
}

.el-menu {
  border-right: none;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>