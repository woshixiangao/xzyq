<template>
  <div class="project-container">
    <div class="operation-bar">
      <el-button type="primary" @click="handleAdd">新增项目</el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="projects"
      border
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="name" label="项目名称"></el-table-column>
      <el-table-column prop="code" label="项目代码"></el-table-column>
      <el-table-column prop="organization_name" label="所属组织"></el-table-column>
      <el-table-column prop="status" label="状态">
        <template slot-scope="scope">
          <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
        </template>
      </el-table-column>
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
        <el-form-item label="项目名称" prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="项目代码" prop="code">
          <el-input v-model="form.code"></el-input>
        </el-form-item>
        <el-form-item label="所属组织" prop="organization_id">
          <el-select v-model="form.organization_id" placeholder="请选择">
            <el-option
              v-for="item in organizations"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择">
            <el-option label="未开始" value="not_started"></el-option>
            <el-option label="进行中" value="in_progress"></el-option>
            <el-option label="已完成" value="completed"></el-option>
            <el-option label="已暂停" value="suspended"></el-option>
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
  name: 'ProjectList',
  data() {
    return {
      loading: false,
      projects: [],
      organizations: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      form: {
        name: '',
        code: '',
        organization_id: '',
        status: 'not_started',
        description: ''
      },
      rules: {
        name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
        code: [{ required: true, message: '请输入项目代码', trigger: 'blur' }],
        organization_id: [{ required: true, message: '请选择所属组织', trigger: 'change' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }]
      }
    }
  },
  created() {
    this.fetchData()
    this.fetchOrganizations()
  },
  methods: {
    getStatusType(status) {
      const types = {
        not_started: 'info',
        in_progress: 'primary',
        completed: 'success',
        suspended: 'warning'
      }
      return types[status] || 'info'
    },
    getStatusText(status) {
      const texts = {
        not_started: '未开始',
        in_progress: '进行中',
        completed: '已完成',
        suspended: '已暂停'
      }
      return texts[status] || status
    },
    async fetchData() {
      this.loading = true
      try {
        const response = await axios.get('/projects', {
          params: {
            page: this.currentPage,
            per_page: this.pageSize
          }
        })
        this.projects = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取项目列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchOrganizations() {
      try {
        const response = await axios.get('/organizations', {
          params: { per_page: 1000 }
        })
        this.organizations = response.data.items
      } catch (error) {
        this.$message.error('获取组织列表失败')
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
      this.dialogTitle = '新增项目'
      this.form = {
        name: '',
        code: '',
        organization_id: '',
        status: 'not_started',
        description: ''
      }
      this.dialogVisible = true
    },
    handleEdit(row) {
      this.dialogTitle = '编辑项目'
      this.form = { ...row }
      this.dialogVisible = true
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该项目吗？')
        await axios.delete(`/projects/${row.id}`)
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
              await axios.put(`/projects/${this.form.id}`, this.form)
              this.$message.success('更新成功')
            } else {
              await axios.post('/projects', this.form)
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
.project-container {
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