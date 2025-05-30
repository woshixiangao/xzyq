<template>
  <div class="organization-page">
    <div class="header">
      <h2>组织管理</h2>
      <el-button type="primary" @click="showCreateDialog">新建组织</el-button>
    </div>

    <!-- 组织列表 -->
    <el-table :data="organizations" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="组织名称" />
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
    const dialogVisible = ref(false)
    const dialogTitle = ref('')
    const form = ref({
      id: null,
      name: '',
      description: ''
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

    // 显示创建对话框
    const showCreateDialog = () => {
      form.value = {
        id: null,
        name: '',
        description: ''
      }
      dialogTitle.value = '新建组织'
      dialogVisible.value = true
    }

    // 显示编辑对话框
    const showEditDialog = (row) => {
      form.value = { ...row }
      dialogTitle.value = '编辑组织'
      dialogVisible.value = true
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
        await ElMessageBox.confirm('确定要删除这个组织吗？', '提示', {
          type: 'warning'
        })
        const token = localStorage.getItem('token')
        await axios.delete(`/api/organizations/${id}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        ElMessage.success('删除成功')
        fetchOrganizations()
      } catch (error) {
        if (error !== 'cancel') {
          if (error.response && error.response.status === 401) {
            ElMessage.error('登录已过期，请重新登录')
            router.push('/login')
          } else {
            ElMessage.error('删除失败')
          }
        }
      }
    }

    onMounted(() => {
      fetchOrganizations()
    })

    return {
      organizations,
      dialogVisible,
      dialogTitle,
      form,
      showCreateDialog,
      showEditDialog,
      handleSubmit,
      handleDelete
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