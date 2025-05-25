<template>
  <div class="user-container">
    <div class="operation-bar">
      <el-button type="primary" @click="handleAdd">新增用户</el-button>
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-select v-model="searchForm.role_id" placeholder="选择角色" clearable>
            <el-option
              v-for="item in roles"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.is_active" placeholder="状态" clearable>
            <el-option label="启用" :value="true"></el-option>
            <el-option label="禁用" :value="false"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="loading"
      :data="users"
      border
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="username" label="用户名"></el-table-column>
      <el-table-column prop="email" label="邮箱"></el-table-column>
      <el-table-column prop="full_name" label="姓名"></el-table-column>
      <el-table-column prop="phone" label="电话"></el-table-column>
      <el-table-column prop="role_names" label="角色">
        <template slot-scope="scope">
          <el-tag
            v-for="role in scope.row.roles"
            :key="role.id"
            size="small"
            style="margin-right: 5px"
          >{{ role.name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_active" label="状态">
        <template slot-scope="scope">
          <el-tag :type="scope.row.is_active ? 'success' : 'danger'">
            {{ scope.row.is_active ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_login_at" label="最后登录时间"></el-table-column>
      <el-table-column label="操作" width="250">
        <template slot-scope="scope">
          <el-button
            size="mini"
            @click="handleEdit(scope.row)"
          >编辑</el-button>
          <el-button
            size="mini"
            :type="scope.row.is_active ? 'warning' : 'success'"
            @click="handleToggleStatus(scope.row)"
          >{{ scope.row.is_active ? '禁用' : '启用' }}</el-button>
          <el-button
            size="mini"
            type="danger"
            @click="handleDelete(scope.row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="currentPage"
      :page-sizes="[10, 20, 50, 100]"
      :page-size="pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
    ></el-pagination>

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible">
      <el-form :model="form" :rules="rules" ref="form" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!form.id">
          <el-input v-model="form.password" type="password"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email"></el-input>
        </el-form-item>
        <el-form-item label="姓名" prop="full_name">
          <el-input v-model="form.full_name"></el-input>
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="form.phone"></el-input>
        </el-form-item>
        <el-form-item label="角色" prop="role_ids">
          <el-select v-model="form.role_ids" multiple placeholder="请选择">
            <el-option
              v-for="item in roles"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="form.is_active"
            active-text="启用"
            inactive-text="禁用"
          ></el-switch>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSubmit">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'UserList',
  data() {
    return {
      loading: false,
      users: [],
      roles: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      searchForm: {
        role_id: '',
        is_active: ''
      },
      form: {
        username: '',
        password: '',
        email: '',
        full_name: '',
        phone: '',
        role_ids: [],
        is_active: true
      },
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        full_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
        role_ids: [{ required: true, message: '请选择角色', trigger: 'change' }]
      }
    }
  },
  created() {
    this.fetchData()
    this.fetchRoles()
  },
  methods: {
    async fetchData() {
      this.loading = true
      try {
        const response = await axios.get('/users', {
          params: {
            page: this.currentPage,
            per_page: this.pageSize,
            role_id: this.searchForm.role_id,
            is_active: this.searchForm.is_active
          }
        })
        this.users = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取用户列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchRoles() {
      try {
        const response = await axios.get('/roles')
        this.roles = response.data.items
      } catch (error) {
        this.$message.error('获取角色列表失败')
      }
    },
    handleSearch() {
      this.currentPage = 1
      this.fetchData()
    },
    handleReset() {
      this.searchForm = {
        role_id: '',
        is_active: ''
      }
      this.handleSearch()
    },
    handleSizeChange(val) {
      this.pageSize = val
      this.fetchData()
    },
    handleCurrentChange(val) {
      this.currentPage = val
      this.fetchData()
    },
    handleAdd() {
      this.dialogTitle = '新增用户'
      this.form = {
        username: '',
        password: '',
        email: '',
        full_name: '',
        phone: '',
        role_ids: [],
        is_active: true
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        if (this.$refs.form) {
          this.$refs.form.clearValidate()
        }
      })
    },
    handleEdit(row) {
      this.dialogTitle = '编辑用户'
      this.form = {
        ...row,
        role_ids: row.roles.map(role => role.id)
      }
      delete this.form.password
      this.dialogVisible = true
      this.$nextTick(() => {
        if (this.$refs.form) {
          this.$refs.form.clearValidate()
        }
      })
    },
    async handleToggleStatus(row) {
      try {
        await axios.put(`/users/${row.id}`, {
          is_active: !row.is_active
        })
        this.$message.success(`${row.is_active ? '禁用' : '启用'}成功`)
        this.fetchData()
      } catch (error) {
        this.$message.error('操作失败')
      }
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该用户吗？')
        await axios.delete(`/users/${row.id}`)
        this.$message.success('删除成功')
        this.fetchData()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    async handleSubmit() {
      this.$refs.form.validate(async valid => {
        if (valid) {
          try {
            if (this.form.id) {
              await axios.put(`/users/${this.form.id}`, this.form)
              this.$message.success('更新成功')
            } else {
              await axios.post('/users', this.form)
              this.$message.success('创建成功')
            }
            this.dialogVisible = false
            this.fetchData()
          } catch (error) {
            this.$message.error(error.response?.data?.message || '操作失败')
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.user-container {
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
}

.operation-bar {
  margin-bottom: 20px;
}

.search-form {
  float: right;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>