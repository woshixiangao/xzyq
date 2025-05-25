<template>
  <div class="role-container">
    <div class="operation-bar">
      <el-button type="primary" @click="handleAdd">新增角色</el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="roles"
      border
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="name" label="角色名称"></el-table-column>
      <el-table-column prop="description" label="描述"></el-table-column>
      <el-table-column label="权限">
        <template slot-scope="scope">
          <div v-if="scope.row.permissions" class="permission-tags">
            <template v-for="(perms, module) in JSON.parse(scope.row.permissions)">
              <el-tag
                :key="module"
                size="small"
                style="margin: 2px"
              >
                {{ getModuleText(module) }}
                <template v-for="perm in perms">
                  <el-tag
                    :key="`${module}-${perm}`"
                    size="mini"
                    type="info"
                    style="margin: 0 2px"
                  >{{ getPermissionText(perm) }}</el-tag>
                </template>
              </el-tag>
            </template>
          </div>
        </template>
      </el-table-column>
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
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="form.description"></el-input>
        </el-form-item>
        <el-form-item label="权限配置" prop="permissions">
          <div class="permissions-config">
            <div v-for="module in modules" :key="module.key" class="permission-module">
              <div class="module-header">
                <el-checkbox
                  :indeterminate="isModuleIndeterminate(module.key)"
                  v-model="moduleCheckAll[module.key]"
                  @change="handleModuleCheckAllChange($event, module.key)"
                >{{ module.label }}</el-checkbox>
              </div>
              <div class="module-permissions">
                <el-checkbox-group
                  v-model="selectedPermissions[module.key]"
                  @change="handlePermissionChange(module.key)"
                >
                  <el-checkbox
                    v-for="perm in permissions"
                    :key="`${module.key}-${perm.key}`"
                    :label="perm.key"
                  >{{ perm.label }}</el-checkbox>
                </el-checkbox-group>
              </div>
            </div>
          </div>
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
  name: 'RoleList',
  data() {
    return {
      loading: false,
      roles: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      modules: [
        { key: 'organizations', label: '组织管理' },
        { key: 'projects', label: '项目管理' },
        { key: 'products', label: '产品管理' },
        { key: 'users', label: '用户管理' },
        { key: 'roles', label: '角色管理' },
        { key: 'logs', label: '日志管理' }
      ],
      permissions: [
        { key: 'create', label: '创建' },
        { key: 'read', label: '查看' },
        { key: 'update', label: '更新' },
        { key: 'delete', label: '删除' }
      ],
      moduleCheckAll: {},
      selectedPermissions: {},
      form: {
        name: '',
        description: '',
        permissions: '{}'
      },
      rules: {
        name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.initPermissionData()
    this.fetchData()
  },
  methods: {
    initPermissionData() {
      this.modules.forEach(module => {
        this.moduleCheckAll[module.key] = false
        this.selectedPermissions[module.key] = []
      })
    },
    getModuleText(module) {
      const found = this.modules.find(m => m.key === module)
      return found ? found.label : module
    },
    getPermissionText(permission) {
      const found = this.permissions.find(p => p.key === permission)
      return found ? found.label : permission
    },
    isModuleIndeterminate(moduleKey) {
      const selected = this.selectedPermissions[moduleKey]
      return selected.length > 0 && selected.length < this.permissions.length
    },
    handleModuleCheckAllChange(checked, moduleKey) {
      this.selectedPermissions[moduleKey] = checked
        ? this.permissions.map(p => p.key)
        : []
    },
    handlePermissionChange(moduleKey) {
      const itemCount = this.permissions.length
      const checkedCount = this.selectedPermissions[moduleKey].length
      this.moduleCheckAll[moduleKey] = checkedCount === itemCount
    },
    async fetchData() {
      this.loading = true
      try {
        const response = await axios.get('/roles', {
          params: {
            page: this.currentPage,
            per_page: this.pageSize
          }
        })
        this.roles = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取角色列表失败')
      } finally {
        this.loading = false
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
      this.dialogTitle = '新增角色'
      this.form = {
        name: '',
        description: '',
        permissions: '{}'
      }
      this.initPermissionData()
      this.dialogVisible = true
    },
    handleEdit(row) {
      this.dialogTitle = '编辑角色'
      this.form = { ...row }
      
      // 解析已有权限
      const permissions = JSON.parse(row.permissions)
      this.initPermissionData()
      Object.keys(permissions).forEach(moduleKey => {
        this.selectedPermissions[moduleKey] = permissions[moduleKey]
        this.handlePermissionChange(moduleKey)
      })
      
      this.dialogVisible = true
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该角色吗？')
        await axios.delete(`/roles/${row.id}`)
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
            // 构建权限对象
            const permissions = {}
            Object.keys(this.selectedPermissions).forEach(moduleKey => {
              if (this.selectedPermissions[moduleKey].length > 0) {
                permissions[moduleKey] = this.selectedPermissions[moduleKey]
              }
            })
            this.form.permissions = JSON.stringify(permissions)

            if (this.form.id) {
              await axios.put(`/roles/${this.form.id}`, this.form)
              this.$message.success('更新成功')
            } else {
              await axios.post('/roles', this.form)
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
.role-container {
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

.permission-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.permissions-config {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
}

.permission-module {
  margin-bottom: 15px;
}

.module-header {
  margin-bottom: 10px;
}

.module-permissions {
  padding-left: 20px;
}
</style>