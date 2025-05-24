<template>
  <div class="logs-container">
    <el-card>
      <div slot="header" class="header">
        <span>日志管理</span>
        <div class="filter-container">
          <el-form :inline="true" :model="queryParams">
            <el-form-item label="级别">
              <el-select v-model="queryParams.level" placeholder="选择日志级别" clearable>
                <el-option label="INFO" value="INFO"></el-option>
                <el-option label="ERROR" value="ERROR"></el-option>
                <el-option label="DB" value="DB"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="yyyy-MM-dd"
                @change="handleDateRangeChange">
              </el-date-picker>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchLogs">查询</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
      <el-table :data="logs" border style="width: 100%">
        <el-table-column prop="created_at" label="时间" width="180"></el-table-column>
        <el-table-column prop="level" label="级别" width="100"></el-table-column>
        <el-table-column prop="component" label="组件" width="120"></el-table-column>
        <el-table-column prop="message" label="内容"></el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-size="pageSize"
          layout="total, prev, pager, next"
          :total="total">
        </el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'LogManagement',
  data() {
    return {
      logs: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dateRange: [],
      queryParams: {
        level: '',
        startDate: '',
        endDate: ''
      }
    }
  },
  created() {
    this.fetchLogs()
  },
  methods: {
    handleDateRangeChange(dates) {
      if (dates) {
        this.queryParams.startDate = dates[0]
        this.queryParams.endDate = dates[1]
      } else {
        this.queryParams.startDate = ''
        this.queryParams.endDate = ''
      }
    },
    async fetchLogs() {
      try {
        const response = await this.$axios.get('/api/logs', {
          params: {
            ...this.queryParams,
            page: this.currentPage,
            pageSize: this.pageSize
          }
        })
        this.logs = response.data.logs || []
        this.total = response.data.total || 0
      } catch (error) {
        this.$message.error('获取日志失败')
      }
    },
    handleCurrentChange(page) {
      this.currentPage = page
      this.fetchLogs()
    }
  }
}
</script>

<style scoped>
.logs-container {
  padding: 20px;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter-container {
  margin-top: 10px;
}
.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>