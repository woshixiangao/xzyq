<template>
  <el-container class="dashboard-container">
    <el-aside width="200px">
      <div class="logo">
        <h3>XZYQ MES</h3>
      </div>
      <el-menu
        :default-active="$route.path"
        class="el-menu-vertical"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard/organizations">
          <i class="el-icon-office-building"></i>
          <span>组织管理</span>
        </el-menu-item>
        <el-menu-item index="/dashboard/projects">
          <i class="el-icon-folder"></i>
          <span>项目管理</span>
        </el-menu-item>
        <el-menu-item index="/dashboard/products">
          <i class="el-icon-goods"></i>
          <span>产品管理</span>
        </el-menu-item>
        <el-menu-item index="/dashboard/users">
          <i class="el-icon-user"></i>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/dashboard/roles">
          <i class="el-icon-s-check"></i>
          <span>角色管理</span>
        </el-menu-item>
        <el-menu-item index="/dashboard/logs">
          <i class="el-icon-document"></i>
          <span>日志管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header>
        <div class="header-right">
          <el-dropdown trigger="click" @command="handleCommand">
            <span class="el-dropdown-link">
              {{ currentUser.username }}<i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
export default {
  name: 'Dashboard',
  computed: {
    currentUser() {
      return this.$store.state.user || {}
    }
  },
  methods: {
    handleCommand(command) {
      if (command === 'logout') {
        this.$store.dispatch('logout')
        this.$router.push('/login')
      } else if (command === 'profile') {
        // 处理个人信息
      }
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
}

.el-aside {
  background-color: #304156;
  color: #fff;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: #fff;
  background-color: #2b2f3a;
}

.el-menu {
  border-right: none;
}

.el-header {
  background-color: #fff;
  color: #333;
  line-height: 60px;
  border-bottom: 1px solid #e6e6e6;
}

.header-right {
  float: right;
  cursor: pointer;
}

.el-dropdown-link {
  color: #409EFF;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>