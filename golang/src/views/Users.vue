<template>
  <div class="users-container">
    <el-card>
      <div slot="header" class="header">
        <span>用户管理</span>
        <el-button type="primary" size="small" @click="dialogVisible = true">添加用户</el-button>
      </div>
      
      <el-table :data="users" border style="width: 100%">
        <el-table-column prop="username" label="用户名"></el-table-column>
        <el-table-column label="操作" width="200">
          <template slot-scope="scope">
            <el-button size="mini" @click="handleEdit(scope.row)">修改密码</el-button>
            <el-button size="mini" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 添加用户对话框 -->
      <el-dialog title="添加用户" :visible.sync="dialogVisible" width="30%">
        <el-form :model="form" ref="form" :rules="rules">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="form.username"></el-input>
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input type="password" v-model="form.password"></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAdd">确定</el-button>
        </span>
      </el-dialog>

      <!-- 修改密码对话框 -->
      <el-dialog title="修改密码" :visible.sync="editDialogVisible" width="30%">
        <el-form :model="editForm" ref="editForm" :rules="editRules">
          <el-form-item label="新密码" prop="password">
            <el-input type="password" v-model="editForm.password"></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUpdate">确定</el-button>
        </span>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'UserManagement',
  data() {
    return {
      users: [],
      dialogVisible: false,
      editDialogVisible: false,
      form: {
        username: '',
        password: ''
      },
      editForm: {
        username: '',
        password: ''
      },
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      },
      editRules: {
        password: [{ required: true, message: '请输入新密码', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.fetchUsers()
  },
  methods: {
    async fetchUsers() {
      try {
        const response = await this.$axios.get('/api/users')
        this.users = response.data.users || []
      } catch (error) {
        this.$message.error('获取用户列表失败')
      }
    },
    async handleAdd() {
      try {
        await this.$refs.form.validate()
        await this.$axios.post('/api/users', this.form)
        this.$message.success('添加用户成功')
        this.dialogVisible = false
        this.fetchUsers()
        this.$refs.form.resetFields()
      } catch (error) {
        this.$message.error(error.response?.data?.error || '添加用户失败')
      }
    },
    handleEdit(row) {
      this.editForm.username = row.username
      this.editDialogVisible = true
    },
    async handleUpdate() {
      try {
        await this.$refs.editForm.validate()
        await this.$axios.put(`/api/users/${this.editForm.username}/password`, {
          password: this.editForm.password
        })
        this.$message.success('修改密码成功')
        this.editDialogVisible = false
        this.$refs.editForm.resetFields()
      } catch (error) {
        this.$message.error(error.response?.data?.error || '修改密码失败')
      }
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该用户吗？')
        await this.$axios.delete(`/api/users/${row.username}`)
        this.$message.success('删除用户成功')
        this.fetchUsers()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error(error.response?.data?.error || '删除用户失败')
        }
      }
    }
  }
}
</script>

<style scoped>
.users-container {
  padding: 20px;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>