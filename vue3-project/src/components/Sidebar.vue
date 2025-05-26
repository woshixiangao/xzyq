<template>
  <el-menu
    :default-active="activeMenu"
    class="sidebar-menu"
    :router="true"
  >
    <!-- ... existing code ... -->
    
    <el-menu-item index="/organizations">
      <el-icon><Folder /></el-icon>
      <span>组织管理</span>
    </el-menu-item>
    
    <!-- ... existing code ... -->
  </el-menu>
</template>

<script>
import { computed } from 'vue'  // 移除未使用的ref
import { useRoute } from 'vue-router'
import { Folder } from '@element-plus/icons-vue'

export default {
  name: 'SidebarMenu',  // 改为多词组合的组件名
  components: {
    Folder
  },
  setup() {
    const route = useRoute()
    const activeMenu = computed(() => route.path)

    return {
      activeMenu
    }
  }
}
</script>