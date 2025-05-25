# 新质元启(XZYQ) MES系统

## 项目简介

新质元启MES系统是一个基于Go语言开发的小型制造执行系统，采用前后端分离架构，后端使用Gin框架，前端使用Vue2框架。系统提供组织管理、项目管理、产品管理、用户管理、角色管理和日志管理等核心功能。

## 技术栈

### 后端
- Go 1.21
- Gin Web框架
- GORM ORM框架
- PostgreSQL 数据库
- JWT认证

### 前端
- Vue2
- Element UI
- Axios
- Vue Router
- Vuex

## 功能模块

1. 组织管理
   - 组织的创建、查询、修改、删除
   - 组织层级关系管理

2. 项目管理
   - 项目的创建、查询、修改、删除
   - 项目状态跟踪
   - 项目与组织关联

3. 产品管理
   - 产品的创建、查询、修改、删除
   - 产品分类管理
   - 产品与项目关联

4. 用户管理
   - 用户的创建、查询、修改、删除
   - 用户角色分配
   - 用户状态管理

5. 角色管理
   - 角色的创建、查询、修改、删除
   - 角色权限配置

6. 日志管理
   - 系统操作日志记录
   - 日志查询和导出

## 快速开始

### 环境要求
- Go 1.21或更高版本
- PostgreSQL 14或更高版本
- Node.js 16或更高版本

### 后端启动

1. 克隆项目
```bash
git clone [项目地址]
cd golang
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境变量
复制`.env.example`文件为`.env`，并根据实际情况修改配置：
```bash
cp .env.example .env
```

4. 启动服务
```bash
go run main.go
```

服务将在`http://localhost:8080`启动

### 前端启动

1. 进入前端目录
```bash
cd ../frontend
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run serve
```

前端将在`http://localhost:8081`启动

## API文档

### 认证相关
- POST /api/auth/register - 用户注册
- POST /api/auth/login - 用户登录

### 组织管理
- GET /api/organizations - 获取组织列表
- POST /api/organizations - 创建组织
- GET /api/organizations/:id - 获取组织详情
- PUT /api/organizations/:id - 更新组织
- DELETE /api/organizations/:id - 删除组织

### 项目管理
- GET /api/projects - 获取项目列表
- POST /api/projects - 创建项目
- GET /api/projects/:id - 获取项目详情
- PUT /api/projects/:id - 更新项目
- DELETE /api/projects/:id - 删除项目

### 产品管理
- GET /api/products - 获取产品列表
- POST /api/products - 创建产品
- GET /api/products/:id - 获取产品详情
- PUT /api/products/:id - 更新产品
- DELETE /api/products/:id - 删除产品

### 用户管理
- GET /api/users - 获取用户列表
- POST /api/users - 创建用户
- GET /api/users/:id - 获取用户详情
- PUT /api/users/:id - 更新用户
- DELETE /api/users/:id - 删除用户

### 角色管理
- GET /api/roles - 获取角色列表
- POST /api/roles - 创建角色
- GET /api/roles/:id - 获取角色详情
- PUT /api/roles/:id - 更新角色
- DELETE /api/roles/:id - 删除角色

### 日志管理
- GET /api/logs - 获取日志列表
- GET /api/logs/:id - 获取日志详情

## 默认账户

系统初始化时会创建一个默认的管理员账户：
- 用户名：admin
- 密码：admin123

首次登录后请及时修改密码。

## 注意事项

1. 生产环境部署时请修改JWT密钥
2. 确保数据库配置正确且可访问
3. 建议定期备份数据库
4. 请及时更新系统依赖包以修复潜在的安全问题

## 许可证

[待定]