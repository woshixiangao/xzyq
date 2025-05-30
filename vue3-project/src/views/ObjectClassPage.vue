<template>
  <div class="object-class-page">
    <div class="header">
      <h2>对象类管理</h2>
      <el-button type="primary" @click="showCreateDialog">新建对象类</el-button>
    </div>

    <!-- 对象类列表 -->
    <el-table :data="objectClasses" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" width="150" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column prop="organization.name" label="所属组织" width="150" />
      <el-table-column prop="created_by_user.username" label="创建人" width="120" />
      <el-table-column prop="updated_at" label="更新时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.updated_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="showEditDialog(row)">编辑</el-button>
          <el-popconfirm
            title="确定要删除该对象类吗？"
            @confirm="handleDelete(row.id)"
          >
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'create' ? '新建对象类' : '编辑对象类'"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="form" label-width="100px" ref="formRef">
        <el-form-item label="名称" prop="name" :rules="[{ required: true, message: '请输入名称' }]">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'ObjectClassPage',
  setup() {
    const loading = ref(false)
    const submitting = ref(false)
    const dialogVisible = ref(false)
    const dialogType = ref('create')
    const currentId = ref(null)
    const objectClasses = ref([])
    const organizations = ref([])
    const formRef = ref(null)
    const form = ref({
      name: '',
      description: '',
    })

    // 获取对象类列表
    const fetchObjectClasses = async () => {
      loading.value = true
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get('/api/object-classes', {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        objectClasses.value = response.data
      } catch (error) {
        console.error('Error fetching object classes:', error)
        ElMessage.error('获取对象类列表失败')
      } finally {
        loading.value = false
      }
    }

    // 获取组织列表
    const fetchOrganizations = async () => {
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get('/api/organizations', {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        organizations.value = response.data
      } catch (error) {
        console.error('Error fetching organizations:', error)
        ElMessage.error('获取组织列表失败')
      }
    }

    // 显示创建对话框
    const showCreateDialog = () => {
      dialogType.value = 'create'
      currentId.value = null
      form.value = {
        name: '',
        description: '',
      }
      dialogVisible.value = true
    }

    // 显示编辑对话框
    const showEditDialog = (row) => {
      dialogType.value = 'edit'
      currentId.value = row.id
      form.value = {
        name: row.name,
        description: row.description,
      }
      dialogVisible.value = true
    }

    // 提交表单
    const handleSubmit = async () => {
      if (!formRef.value) return
      
      await formRef.value.validate(async (valid) => {
        if (!valid) return

        submitting.value = true
        try {
          const token = localStorage.getItem('token')
          const headers = { 'Authorization': `Bearer ${token}` }
          
          if (dialogType.value === 'create') {
            await axios.post('/api/object-classes', form.value, { headers })
            ElMessage.success('创建成功')
          } else {
            await axios.put(`/api/object-classes/${currentId.value}`, form.value, { headers })
            ElMessage.success('更新成功')
          }
          
          dialogVisible.value = false
          fetchObjectClasses()
        } catch (error) {
          console.error('Error submitting form:', error)
          ElMessage.error(error.response?.data?.error || '操作失败')
        } finally {
          submitting.value = false
        }
      })
    }

    // 删除对象类
    const handleDelete = async (id) => {
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`/api/object-classes/${id}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        ElMessage.success('删除成功')
        fetchObjectClasses()
      } catch (error) {
        console.error('Error deleting object class:', error)
        ElMessage.error(error.response?.data?.error || '删除失败')
      }
    }

    // 格式化日期
    const formatDate = (date) => {
      if (!date) return '-'
      return new Date(date).toLocaleString()
    }

    onMounted(() => {
      fetchObjectClasses()
      fetchOrganizations()
    })

    return {
      loading,
      submitting,
      dialogVisible,
      dialogType,
      currentId,
      objectClasses,
      organizations,
      form,
      formRef,
      showCreateDialog,
      showEditDialog,
      handleSubmit,
      handleDelete,
      formatDate
    }
  }
}
</script>

<style scoped>
.object-class-page {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.el-table {
  margin-top: 20px;
}
</style> 