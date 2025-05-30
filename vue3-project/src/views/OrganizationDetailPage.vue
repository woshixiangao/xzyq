<template>
  <div class="organization-detail-page">
    <div class="header">
      <h2>组织详情</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 详细信息标签页 -->
      <el-tab-pane label="详细信息" name="info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="组织ID">{{ organization.id }}</el-descriptions-item>
          <el-descriptions-item label="组织名称">{{ organization.name }}</el-descriptions-item>
          <el-descriptions-item label="父级组织">{{ organization.parent_org?.name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="用户数量">
            <el-tag>{{ organization.user_count }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ organization.created_at }}</el-descriptions-item>
          <el-descriptions-item label="描述">{{ organization.description || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-tab-pane>

      <!-- 用户列表标签页 -->
      <el-tab-pane label="用户列表" name="users">
        <el-table :data="users" style="width: 100%" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="username" label="用户名" width="150" />
          <el-table-column prop="email" label="邮箱" width="200">
            <template #default="{ row }">
              {{ row.email || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="phone" label="手机号" width="120">
            <template #default="{ row }">
              {{ row.phone || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="role" label="角色" width="100">
            <template #default="{ row }">
              <el-tag :type="row.role === 'admin' ? 'danger' : 'success'">
                {{ row.role }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="is_active" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'info'">
                {{ row.is_active ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-popconfirm
                title="确定要删除该用户吗？"
                @confirm="handleDeleteUser(row.id)"
              >
                <template #reference>
                  <el-button size="small" type="danger">删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="暂无用户数据" />
          </template>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'OrganizationDetailPage',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const organization = ref({})
    const users = ref([])
    const loading = ref(false)
    const activeTab = ref('info')

    // 获取组织详情
    const fetchOrganizationDetail = async () => {
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get(`/api/organizations/${route.params.id}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        organization.value = response.data
      } catch (error) {
        if (error.response?.status === 401) {
          ElMessage.error('登录已过期，请重新登录')
          router.push('/login')
        } else {
          ElMessage.error('获取组织详情失败')
        }
      }
    }

    // 获取组织用户列表
    const fetchOrganizationUsers = async () => {
      loading.value = true
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get(`/api/organizations/${route.params.id}/users`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        users.value = response.data || []
        console.log('Fetched users:', users.value)
      } catch (error) {
        console.error('Error fetching users:', error)
        if (error.response?.status === 401) {
          ElMessage.error('登录已过期，请重新登录')
          router.push('/login')
        } else {
          ElMessage.error(error.response?.data?.error || '获取用户列表失败')
        }
        users.value = []
      } finally {
        loading.value = false
      }
    }

    // 删除用户
    const handleDeleteUser = async (userId) => {
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`/api/users/${userId}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        ElMessage.success('删除成功')
        // 重新获取组织详情（更新用户数量）和用户列表
        await Promise.all([
          fetchOrganizationDetail(),
          fetchOrganizationUsers()
        ])
      } catch (error) {
        console.error('Error deleting user:', error)
        if (error.response?.status === 401) {
          ElMessage.error('登录已过期，请重新登录')
          router.push('/login')
        } else {
          ElMessage.error(error.response?.data?.error || '删除用户失败')
        }
      }
    }

    // 返回上一页
    const goBack = () => {
      router.back()
    }

    onMounted(() => {
      fetchOrganizationDetail()
      fetchOrganizationUsers()
    })

    return {
      organization,
      users,
      loading,
      activeTab,
      handleDeleteUser,
      goBack
    }
  }
}
</script>

<style scoped>
.organization-detail-page {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.el-descriptions {
  margin: 20px 0;
}
</style> 