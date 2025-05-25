<template>
  <div class="organization-container">
    <div class="operation-bar">
      <el-button type="primary" @click="handleAdd">新增组织</el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="organizations"
      border
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="name" label="组织名称"></el-table-column>
      <el-table-column prop="code" label="组织代码"></el-table-column>
      <el-table-column prop="parent_name" label="上级组织"></el-table-column>
      <el-table-column prop="description" label="描述"></el-table-column>
      <el-table-column label="操作" width="200">
        <template slot-scope="scope">
          <el-button
            size="mini"
            @click="handleEdit(scope.row)"
          >编辑</el-button>
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
        <el-form-item label="组织名称" prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="组织代码" prop="code">
          <el-input v-model="form.code"></el-input>
        </el-form-item>
        <el-form-item label="上级组织">
          <el-select v-model="form.parent_id" clearable placeholder="请选择">
            <el-option
              v-for="item in organizationOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="form.description"></el-input>
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
  name: 'OrganizationList',
  data() {
    return {
      loading: false,
      organizations: [],
      organizationOptions: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      form: {
        name: '',
        code: '',
        parent_id: null,
        description: ''
      },
      rules: {
        name: [{ required: true, message: '请输入组织名称', trigger: 'blur' }],
        code: [{ required: true, message: '请输入组织代码', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.fetchData()
    this.fetchOrganizationOptions()
  },
  methods: {
    async fetchData() {
      this.loading = true
      try {
        const response = await axios.get('/organizations', {
          params: {
            page: this.currentPage,
            per_page: this.pageSize
          }
        })
        this.organizations = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取组织列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchOrganizationOptions() {
      try {
        const response = await axios.get('/organizations', {
          params: { per_page: 1000 }
        })
        this.organizationOptions = response.data.items
      } catch (error) {
        this.$message.error('获取组织选项失败')
      }
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
      this.dialogTitle = '新增组织'
      this.form = {
        name: '',
        code: '',
        parent_id: null,
        description: ''
      }
      this.dialogVisible = true
    },
    handleEdit(row) {
      this.dialogTitle = '编辑组织'
      this.form = { ...row }
      this.dialogVisible = true
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该组织吗？')
        await axios.delete(`/organizations/${row.id}`)
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
              await axios.put(`/organizations/${this.form.id}`, this.form)
              this.$message.success('更新成功')
            } else {
              await axios.post('/organizations', this.form)
              this.$message.success('创建成功')
            }
            this.dialogVisible = false
            this.fetchData()
            this.fetchOrganizationOptions()
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
.organization-container {
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
}

.operation-bar {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>