<template>
  <div class="object-class-detail">
    <div class="header">
      <h2>对象类详情</h2>
      <el-button type="primary" @click="goBack">返回</el-button>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 基本信息标签页 -->
      <el-tab-pane label="详细信息" name="info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ objectClass.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ objectClass.name }}</el-descriptions-item>
          <el-descriptions-item label="描述">{{ objectClass.description }}</el-descriptions-item>
          <el-descriptions-item label="所属组织">{{ objectClass.organization?.name }}</el-descriptions-item>
          <el-descriptions-item label="创建人">{{ objectClass.created_by_user?.username }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(objectClass.updated_at) }}</el-descriptions-item>
          <el-descriptions-item label="父对象类">{{ objectClass.parent?.name || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-tab-pane>

      <!-- 子对象类标签页 -->
      <el-tab-pane label="子对象类" name="children">
        <div class="sub-header">
          <el-button type="primary" @click="showCreateDialog">新增子对象类</el-button>
        </div>

        <el-table :data="children" v-loading="loading" style="width: 100%">
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
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-popconfirm
                title="确定要删除该子对象类吗？"
                @confirm="handleDelete(row.id)"
              >
                <template #reference>
                  <el-button size="small" type="danger">删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 创建子对象类对话框 -->
    <el-dialog
      title="新增子对象类"
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
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'ObjectClassDetailPage',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const activeTab = ref('info')
    const loading = ref(false)
    const submitting = ref(false)
    const dialogVisible = ref(false)
    const objectClass = ref({})
    const children = ref([])
    const formRef = ref(null)
    const form = ref({
      name: '',
      description: '',
    })

    // 获取对象类详情
    const fetchObjectClass = async () => {
      loading.value = true
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get(`/api/object-classes/${route.params.id}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        objectClass.value = response.data
      } catch (error) {
        console.error('Error fetching object class:', error)
        ElMessage.error('获取对象类详情失败')
      } finally {
        loading.value = false
      }
    }

    // 获取子对象类列表
    const fetchChildren = async () => {
      loading.value = true
      try {
        const token = localStorage.getItem('token')
        const response = await axios.get(`/api/object-classes/${route.params.id}/children`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        children.value = response.data
      } catch (error) {
        console.error('Error fetching children:', error)
        ElMessage.error('获取子对象类列表失败')
      } finally {
        loading.value = false
      }
    }

    // 显示创建对话框
    const showCreateDialog = () => {
      form.value = {
        name: '',
        description: '',
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
          
          await axios.post(`/api/object-classes/${route.params.id}/children`, form.value, { headers })
          ElMessage.success('创建成功')
          dialogVisible.value = false
          fetchChildren()
        } catch (error) {
          console.error('Error submitting form:', error)
          ElMessage.error(error.response?.data?.error || '操作失败')
        } finally {
          submitting.value = false
        }
      })
    }

    // 删除子对象类
    const handleDelete = async (id) => {
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`/api/object-classes/${id}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        ElMessage.success('删除成功')
        fetchChildren()
      } catch (error) {
        console.error('Error deleting object class:', error)
        ElMessage.error(error.response?.data?.error || '删除失败')
      }
    }

    // 返回上一页
    const goBack = () => {
      router.back()
    }

    // 格式化日期
    const formatDate = (date) => {
      if (!date) return '-'
      return new Date(date).toLocaleString()
    }

    onMounted(() => {
      fetchObjectClass()
      fetchChildren()
    })

    return {
      activeTab,
      loading,
      submitting,
      dialogVisible,
      objectClass,
      children,
      form,
      formRef,
      showCreateDialog,
      handleSubmit,
      handleDelete,
      goBack,
      formatDate
    }
  }
}
</script>

<style scoped>
.object-class-detail {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.sub-header {
  margin-bottom: 20px;
}

.el-descriptions {
  margin: 20px 0;
}
</style> 