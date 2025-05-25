<template>
  <div class="log-container">
    <div class="operation-bar">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-select v-model="searchForm.user_id" placeholder="操作用户" clearable>
            <el-option
              v-for="item in users"
              :key="item.id"
              :label="item.username"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.operation_type" placeholder="操作类型" clearable>
            <el-option label="创建" value="create"></el-option>
            <el-option label="更新" value="update"></el-option>
            <el-option label="删除" value="delete"></el-option>
            <el-option label="查询" value="read"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.resource_type" placeholder="资源类型" clearable>
            <el-option label="组织" value="organization"></el-option>
            <el-option label="项目" value="project"></el-option>
            <el-option label="产品" value="product"></el-option>
            <el-option label="用户" value="user"></el-option>
            <el-option label="角色" value="role"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="searchForm.time_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="yyyy-MM-dd"
          ></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button type="success" @click="handleExport">导出</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="loading"
      :data="logs"
      border
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="user.username" label="操作用户"></el-table-column>
      <el-table-column prop="operation_type" label="操作类型">
        <template slot-scope="scope">
          <el-tag :type="getOperationTypeTag(scope.row.operation_type)">
            {{ getOperationTypeText(scope.row.operation_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="resource_type" label="资源类型">
        <template slot-scope="scope">
          {{ getResourceTypeText(scope.row.resource_type) }}
        </template>
      </el-table-column>
      <el-table-column prop="resource_id" label="资源ID"></el-table-column>
      <el-table-column prop="description" label="操作描述"></el-table-column>
      <el-table-column prop="ip_address" label="IP地址"></el-table-column>
      <el-table-column prop="created_at" label="操作时间" width="180"></el-table-column>
      <el-table-column label="操作" width="100">
        <template slot-scope="scope">
          <el-button
            size="mini"
            @click="handleViewDetails(scope.row)"
          >详情</el-button>
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

    <!-- 日志详情对话框 -->
    <el-dialog title="日志详情" :visible.sync="detailsVisible" width="60%">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="操作用户">{{ currentLog.user?.username }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ currentLog.created_at }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{ getOperationTypeText(currentLog.operation_type) }}</el-descriptions-item>
        <el-descriptions-item label="资源类型">{{ getResourceTypeText(currentLog.resource_type) }}</el-descriptions-item>
        <el-descriptions-item label="资源ID">{{ currentLog.resource_id }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip_address }}</el-descriptions-item>
        <el-descriptions-item label="操作描述" :span="2">{{ currentLog.description }}</el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <pre>{{ formatJson(currentLog.request_data) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应结果" :span="2">
          <pre>{{ formatJson(currentLog.response_data) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'LogList',
  data() {
    return {
      loading: false,
      logs: [],
      users: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      detailsVisible: false,
      currentLog: {},
      searchForm: {
        user_id: '',
        operation_type: '',
        resource_type: '',
        time_range: []
      }
    }
  },
  created() {
    this.fetchData()
    this.fetchUsers()
  },
  methods: {
    getOperationTypeTag(type) {
      const tags = {
        create: 'success',
        update: 'warning',
        delete: 'danger',
        read: 'info'
      }
      return tags[type] || 'info'
    },
    getOperationTypeText(type) {
      const texts = {
        create: '创建',
        update: '更新',
        delete: '删除',
        read: '查询'
      }
      return texts[type] || type
    },
    getResourceTypeText(type) {
      const texts = {
        organization: '组织',
        project: '项目',
        product: '产品',
        user: '用户',
        role: '角色'
      }
      return texts[type] || type
    },
    formatJson(json) {
      try {
        return JSON.stringify(JSON.parse(json), null, 2)
      } catch (e) {
        return json
      }
    },
    async fetchData() {
      this.loading = true
      try {
        const params = {
          page: this.currentPage,
          per_page: this.pageSize,
          user_id: this.searchForm.user_id,
          operation_type: this.searchForm.operation_type,
          resource_type: this.searchForm.resource_type
        }
        
        if (this.searchForm.time_range?.length === 2) {
          params.start_date = this.searchForm.time_range[0]
          params.end_date = this.searchForm.time_range[1]
        }

        const response = await axios.get('/logs', { params })
        this.logs = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取日志列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchUsers() {
      try {
        const response = await axios.get('/users', {
          params: { per_page: 1000 }
        })
        this.users = response.data.items
      } catch (error) {
        this.$message.error('获取用户列表失败')
      }
    },
    handleSearch() {
      this.currentPage = 1
      this.fetchData()
    },
    handleReset() {
      this.searchForm = {
        user_id: '',
        operation_type: '',
        resource_type: '',
        time_range: []
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
    handleViewDetails(log) {
      this.currentLog = log
      this.detailsVisible = true
    },
    async handleExport() {
      try {
        const params = {
          user_id: this.searchForm.user_id,
          operation_type: this.searchForm.operation_type,
          resource_type: this.searchForm.resource_type
        }
        
        if (this.searchForm.time_range?.length === 2) {
          params.start_date = this.searchForm.time_range[0]
          params.end_date = this.searchForm.time_range[1]
        }

        const response = await axios.get('/logs/export', {
          params,
          responseType: 'blob'
        })

        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'system_logs.xlsx')
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      } catch (error) {
        this.$message.error('导出日志失败')
      }
    }
  }
}
</script>

<style scoped>
.log-container {
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

pre {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>