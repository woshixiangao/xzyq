<template>
  <el-menu
    :default-active="activeMenu"
    class="side-menu"
    @select="handleSelect"
  >
    <el-menu-item index="/users">
      <el-icon><User /></el-icon>
      <span>用户管理</span>
    </el-menu-item>

    <el-menu-item index="/organizations">
      <el-icon><Folder /></el-icon>
      <span>组织管理</span>
    </el-menu-item>

    <el-menu-item index="/object-classes">
      <el-icon><Box /></el-icon>
      <span>对象类</span>
    </el-menu-item>
  </el-menu>
</template>

<script>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { User, Folder, Box } from '@element-plus/icons-vue'

export default {
  name: 'SideMenu',
  components: {
    User,
    Folder,
    Box
  },
  setup() {
    const route = useRoute()
    const router = useRouter()

    const activeMenu = computed(() => route.path)

    const handleSelect = (index) => {
      router.push(index)
    }

    return {
      activeMenu,
      handleSelect
    }
  }
}
</script>

<style scoped>
.side-menu {
  height: 100%;
  border-right: none;
}

.el-menu-item {
  display: flex;
  align-items: center;
}

.el-icon {
  margin-right: 8px;
}
</style> 