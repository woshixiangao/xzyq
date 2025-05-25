<template>
  <div class="product-container">
    <div class="operation-bar">
      <el-button type="primary" @click="handleAdd">新增产品</el-button>
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-select v-model="searchForm.project_id" placeholder="选择项目" clearable>
            <el-option
              v-for="item in projects"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.category" placeholder="产品类别" clearable>
            <el-option label="电子产品" value="electronics"></el-option>
            <el-option label="机械设备" value="machinery"></el-option>
            <el-option label="原材料" value="raw_material"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchForm.keyword" placeholder="产品名称/代码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="loading"
      :data="products"
      border
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="name" label="产品名称"></el-table-column>
      <el-table-column prop="code" label="产品代码"></el-table-column>
      <el-table-column prop="category" label="类别">
        <template slot-scope="scope">
          <el-tag>{{ getCategoryText(scope.row.category) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="project_name" label="所属项目"></el-table-column>
      <el-table-column prop="specifications" label="规格"></el-table-column>
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
        <el-form-item label="产品名称" prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="产品代码" prop="code">
          <el-input v-model="form.code"></el-input>
        </el-form-item>
        <el-form-item label="所属项目" prop="project_id">
          <el-select v-model="form.project_id" placeholder="请选择">
            <el-option
              v-for="item in projects"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="产品类别" prop="category">
          <el-select v-model="form.category" placeholder="请选择">
            <el-option label="电子产品" value="electronics"></el-option>
            <el-option label="机械设备" value="machinery"></el-option>
            <el-option label="原材料" value="raw_material"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="规格">
          <el-input v-model="form.specifications"></el-input>
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
  name: 'ProductList',
  data() {
    return {
      loading: false,
      products: [],
      projects: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      searchForm: {
        project_id: '',
        category: '',
        keyword: ''
      },
      form: {
        name: '',
        code: '',
        project_id: '',
        category: '',
        specifications: '',
        description: ''
      },
      rules: {
        name: [{ required: true, message: '请输入产品名称', trigger: 'blur' }],
        code: [{ required: true, message: '请输入产品代码', trigger: 'blur' }],
        project_id: [{ required: true, message: '请选择所属项目', trigger: 'change' }],
        category: [{ required: true, message: '请选择产品类别', trigger: 'change' }]
      }
    }
  },
  created() {
    this.fetchData()
    this.fetchProjects()
  },
  methods: {
    getCategoryText(category) {
      const texts = {
        electronics: '电子产品',
        machinery: '机械设备',
        raw_material: '原材料',
        other: '其他'
      }
      return texts[category] || category
    },
    async fetchData() {
      this.loading = true
      try {
        const response = await axios.get('/products', {
          params: {
            page: this.currentPage,
            per_page: this.pageSize,
            project_id: this.searchForm.project_id,
            category: this.searchForm.category,
            keyword: this.searchForm.keyword
          }
        })
        this.products = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取产品列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchProjects() {
      try {
        const response = await axios.get('/projects', {
          params: { per_page: 1000 }
        })
        this.projects = response.data.items
      } catch (error) {
        this.$message.error('获取项目列表失败')
      }
    },
    handleSearch() {
      this.currentPage = 1
      this.fetchData()
    },
    handleReset() {
      this.searchForm = {
        project_id: '',
        category: '',
        keyword: ''
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
      this.dialogTitle = '新增产品'
      this.form = {
        name: '',
        code: '',
        project_id: '',
        category: '',
        specifications: '',
        description: ''
      }
      this.dialogVisible = true
    },
    handleEdit(row) {
      this.dialogTitle = '编辑产品'
      this.form = { ...row }
      this.dialogVisible = true
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该产品吗？')
        await axios.delete(`/products/${row.id}`)
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
              await axios.put(`/products/${this.form.id}`, this.form)
              this.$message.success('更新成功')
            } else {
              await axios.post('/products', this.form)
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
.product-container {
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