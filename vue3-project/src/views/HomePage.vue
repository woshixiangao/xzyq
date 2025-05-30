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
          <user-profile />
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
import SideMenu from '../components/SideMenu.vue'
import UserProfile from '../components/UserProfile.vue'

export default {
  name: 'HomePage',
  components: {
    SideMenu,
    UserProfile
  },
  setup() {
    const systemName = ref('系统名称')

    return {
      systemName
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
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
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
  background-color: #f0f2f5;
  padding: 20px;
}
</style>