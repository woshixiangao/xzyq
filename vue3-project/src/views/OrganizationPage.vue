<template>
  <div class="organization-page">
    <div class="header">
      <h2>组织管理</h2>
      <el-button type="primary" @click="showCreateDialog">新建组织</el-button>
    </div>

    <!-- 组织列表 -->
    <el-table :data="organizations" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="组织名称">
        <template #default="{ row }">
          <el-link type="primary" @click="goToDetail(row.id)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="父级租户" width="180">
        <template #default="{ row }">
          <span v-if="row.parent_org">{{ row.parent_org.name }}</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="user_count" label="用户数量" width="100">
        <template #default="{ row }">
          <el-tag>{{ row.user_count }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button size="small" @click="showEditDialog(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="组织名称">
          <el-input v-model="form.name" placeholder="请输入组织名称" />
        </el-form-item>
        <el-form-item label="父级租户">
          <el-select v-model="form.parent_id" placeholder="请选择父级租户" clearable>
            <el-option
              v-for="org in allOrganizations"
              :key="org.id"
              :label="org.name"
              :value="org.id"
              :disabled="org.id === form.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入组织描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import axios from 'axios'

export default {
  name: 'OrganizationPage',
  setup() {
    const router = useRouter()
    const organizations = ref([])
    const allOrganizations = ref([]) // 存储所有组织，用于父级租户选择
    const dialogVisible = ref(false)
    const dialogTitle = ref('')
    const form = ref({
      id: null,
      name: '',
      description: '',
      parent_id: null
    })

    // 获取组织列表
    const fetchOrganizations = async () => {
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get('/api/organizations', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        organizations.value = response.data
      } catch (error) {
        if (error.response && error.response.status === 401) {
          ElMessage.error('登录已过期，请重新登录')
          router.push('/login')
        } else {
          ElMessage.error('获取组织列表失败')
        }
      }
    }

    // 获取所有组织列表（用于父级租户选择）
    const fetchAllOrganizations = async () => {
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get('/api/organizations/all', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        allOrganizations.value = response.data
      } catch (error) {
        console.error('获取所有组织列表失败:', error)
      }
    }

    // 显示创建对话框
    const showCreateDialog = () => {
      form.value = {
        id: null,
        name: '',
        description: '',
        parent_id: null
      }
      dialogTitle.value = '新建组织'
      dialogVisible.value = true
      // 获取所有组织列表用于父级租户选择
      fetchAllOrganizations()
    }

    // 显示编辑对话框
    const showEditDialog = (row) => {
      form.value = { 
        ...row,
        parent_id: row.parent_org?.id || null
      }
      dialogTitle.value = '编辑组织'
      dialogVisible.value = true
      // 获取所有组织列表用于父级租户选择
      fetchAllOrganizations()
    }

    // 提交表单
    const handleSubmit = async () => {
      try {
        const token = localStorage.getItem('token')
        if (form.value.id) {
          // 编辑
          await axios.put(`/api/organizations/${form.value.id}`, form.value, {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          })
          ElMessage.success('更新成功')
        } else {
          // 创建
          const response = await axios.post('/api/organizations', form.value, {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          })
          ElMessage.success('创建成功')
          // 显示管理员账号信息
          ElMessage.info(`管理员账号：${response.data.admin_user.username}\n初始密码：${response.data.admin_user.password}`)
        }
        dialogVisible.value = false
        fetchOrganizations()
      } catch (error) {
        if (error.response && error.response.status === 401) {
          ElMessage.error('登录已过期，请重新登录')
          router.push('/login')
        } else {
          ElMessage.error(form.value.id ? '更新失败' : '创建失败')
        }
      }
    }

    // 删除组织
    const handleDelete = async (id) => {
      try {
        const row = organizations.value.find(org => org.id === id)
        if (row.user_count > 0) {
          ElMessage.warning(`该组织下还有 ${row.user_count} 个用户，不能删除`)
          return
        }

        await ElMessageBox.confirm('确定要删除这个组织吗？', '提示', {
          type: 'warning'
        })
        const token = localStorage.getItem('token')
        try {
          await axios.delete(`/api/organizations/${id}`, {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          })
          ElMessage.success('删除成功')
          fetchOrganizations()
        } catch (error) {
          console.error('Delete error:', error.response)
          if (error.response?.status === 401) {
            ElMessage.error('登录已过期，请重新登录')
            router.push('/login')
          } else if (error.response?.data?.error) {
            ElMessage.error(error.response.data.error)
          } else {
            ElMessage.error('删除组织失败，请稍后重试')
          }
        }
      } catch (error) {
        // 用户取消删除操作，不需要提示错误
        if (error !== 'cancel') {
          ElMessage.error('操作取消')
        }
      }
    }

    // 跳转到组织详情页
    const goToDetail = (id) => {
      router.push(`organizations/${id}`)
    }

    onMounted(() => {
      fetchOrganizations()
    })

    return {
      organizations,
      allOrganizations,
      dialogVisible,
      dialogTitle,
      form,
      showCreateDialog,
      showEditDialog,
      handleSubmit,
      handleDelete,
      goToDetail
    }
  }
}
</script>

<style scoped>
.organization-page {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>