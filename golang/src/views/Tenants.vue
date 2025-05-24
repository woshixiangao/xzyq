<template>
  <div class="tenants-container">
    <el-card>
      <div slot="header" class="header">
        <span>租户管理</span>
        <el-button type="primary" size="small" @click="handleAdd">添加租户</el-button>
      </div>
      
      <el-table :data="tenants" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="租户名称"></el-table-column>
        <el-table-column prop="description" label="描述"></el-table-column>
        <el-table-column label="操作" width="200">
          <template slot-scope="scope">
            <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="mini" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 添加/编辑租户对话框 -->
      <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="30%">
        <el-form :model="form" ref="form" :rules="rules">
          <el-form-item label="租户名称" prop="name">
            <el-input v-model="form.name"></el-input>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input type="textarea" v-model="form.description"></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </span>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'TenantManagement',
  data() {
    return {
      tenants: [],
      dialogVisible: false,
      isEdit: false,
      form: {
        id: '',
        name: '',
        description: ''
      },
      rules: {
        name: [{ required: true, message: '请输入租户名称', trigger: 'blur' }]
      }
    }
  },
  computed: {
    dialogTitle() {
      return this.isEdit ? '编辑租户' : '添加租户'
    }
  },
  created() {
    this.fetchTenants()
  },
  methods: {
    async fetchTenants() {
      try {
        const response = await this.$axios.get('/api/tenants')
        this.tenants = response.data.tenants || []
      } catch (error) {
        this.$message.error('获取租户列表失败')
      }
    },
    handleAdd() {
      this.isEdit = false
      this.form = {
        id: '',
        name: '',
        description: ''
      }
      this.dialogVisible = true
    },
    handleEdit(row) {
      this.isEdit = true
      this.form = { ...row }
      this.dialogVisible = true
    },
    async submitForm() {
      try {
        await this.$refs.form.validate()
        if (this.isEdit) {
          await this.$axios.put(`/api/tenants/${this.form.id}`, this.form)
          this.$message.success('更新租户成功')
        } else {
          await this.$axios.post('/api/tenants', this.form)
          this.$message.success('添加租户成功')
        }
        this.dialogVisible = false
        this.fetchTenants()
        this.$refs.form.resetFields()
      } catch (error) {
        this.$message.error(error.response?.data?.error || `${this.isEdit ? '更新' : '添加'}租户失败`)
      }
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该租户吗？')
        await this.$axios.delete(`/api/tenants/${row.id}`)
        this.$message.success('删除租户成功')
        this.fetchTenants()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error(error.response?.data?.error || '删除租户失败')
        }
      }
    }
  }
}
</script>

<style scoped>
.tenants-container {
  padding: 20px;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>