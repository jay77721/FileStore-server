# 文件存储服务器 (FileStore Server)

基于 Go 语言开发的轻量级文件存储服务器，支持文件上传、下载、用户管理和分片上传等功能。

## 🚀 功能特性

- 📁 **文件管理**
  - 文件上传（支持元数据）
  - 文件下载
  - 文件元数据更新
  - 文件删除
  - 文件信息查询

- 🔐 **用户认证**
  - 用户注册
  - 用户登录
  - 用户信息获取
  - 基于 JWT 的身份验证

- ⚡ **分片上传**
  - 支持大文件分片上传
  - 分片上传状态检查
  - 自动合并分片

- 🗄️ **存储后端**
  - MySQL 数据库存储元数据
  - Redis 缓存和会话管理
  - 本地文件系统存储

## 🏗️ 项目架构

```
filestore-server/
├── main.go              # 程序入口和 HTTP 路由
├── db/                  # 数据库操作
│   ├── mysql/conn.go    # MySQL 连接
│   ├── file.go          # 文件相关数据库操作
│   └── user.go          # 用户相关数据库操作
├── handler/             # HTTP 请求处理器
│   ├── auth.go          # 认证中间件
│   ├── handler.go       # 文件上传下载处理器
│   └── user.go          # 用户管理处理器
├── meta/                # 文件元数据管理
│   └── filemeta.go      # 文件元数据结构
├── rd/                  # Redis 操作
│   └── redis.go         # Redis 连接和操作
├── util/                # 工具函数
│   └── chunk.go         # 分片上传工具
├── static/              # 静态文件（前端资源）
├── uploads/             # 上传文件存储
└── go.mod               # Go 模块定义
```

## 🛠️ 技术栈

- **编程语言**: Go 1.24.0
- **数据库**: MySQL
- **缓存**: Redis
- **Web 框架**: net/http (标准库)
- **认证方式**: JWT

## 📋 API 接口

### 文件操作

| 方法 | 接口 | 描述 |
|------|------|------|
| POST | `/file/upload` | 上传文件 |
| GET | `/file/meta` | 获取文件元数据 |
| GET | `/file/query` | 查询文件 |
| GET | `/file/download` | 下载文件 |
| POST | `/file/update` | 更新文件元数据 |
| POST | `/file/delete` | 删除文件 |

### 分片上传

| 方法 | 接口 | 描述 |
|------|------|------|
| POST | `/file/upload/chunk` | 上传文件分片 |
| GET | `/file/upload/status` | 检查分片上传状态 |
| POST | `/file/upload/merge` | 合并上传的分片 |

### 用户操作

| 方法 | 接口 | 描述 |
|------|------|------|
| POST | `/user/signup` | 用户注册 |
| POST | `/user/signin` | 用户登录 |
| GET | `/user/info` | 获取用户信息 |

## 🚦 快速开始

### 环境要求

- Go 1.24.0 或更高版本
- MySQL 数据库
- Redis 服务器

### 安装部署

1. **克隆项目**
   ```bash
   git clone <仓库地址>
   cd filestore-server
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **配置数据库连接**
   - 在 `db/mysql/conn.go` 中更新 MySQL 连接配置
   - 在 `rd/redis.go` 中更新 Redis 连接配置

4. **创建数据库表**
   ```sql
   -- 创建用户表
   CREATE TABLE users (
       id INT AUTO_INCREMENT PRIMARY KEY,
       username VARCHAR(50) UNIQUE NOT NULL,
       password VARCHAR(100) NOT NULL,
       email VARCHAR(100),
       create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- 创建文件元数据表
   CREATE TABLE file_meta (
       id INT AUTO_INCREMENT PRIMARY KEY,
       file_hash VARCHAR(100) NOT NULL,
       file_name VARCHAR(255) NOT NULL,
       file_size BIGINT DEFAULT 0,
       file_path VARCHAR(255) NOT NULL,
       create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       status INT DEFAULT 0
   );
   ```

5. **启动服务器**
   ```bash
   go run main.go
   ```

6. **访问服务**
   - 服务器将在 `http://localhost:8080` 启动
   - 静态文件通过 `/static/` 访问

## 📝 使用示例

### 上传文件
```bash
curl -X POST -F "file=@/path/to/your/file.txt" http://localhost:8080/file/upload
```

### 下载文件
```bash
curl -X GET "http://localhost:8080/file/download?filehash=abc123" --output file.txt
```

### 用户注册
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  http://localhost:8080/user/signup
```

### 用户登录
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  http://localhost:8080/user/signin
```

## ⚙️ 配置说明

### 数据库配置
在相应文件中更新连接字符串：

**MySQL** (`db/mysql/conn.go`):
```go
db, err := sql.Open("mysql", "用户名:密码@tcp(地址:端口)/数据库名")
```

**Redis** (`rd/redis.go`):
```go
client := redis.NewClient(&redis.Options{
    Addr:     "地址:端口",
    Password: "密码",
    DB:       0,
})
```

### 服务器配置
默认在 8080 端口运行。要修改端口，请编辑 `main.go`:
```go
err := http.ListenAndServe(":8080", nil)  // 将 8080 改为所需端口
```

## 🔒 安全注意事项

- 密码在存储前应该进行哈希处理（需要在 user.go 中实现）
- 生产环境请使用 HTTPS
- 验证文件类型和大小限制
- 为上传功能实现速率限制
- 清理文件名和路径，防止路径遍历攻击

## 🚧 开发状态

该项目正在积极开发中，功能可能会发生变化，API 接口可能会修改。

## 🤝 贡献指南

1. Fork 项目仓库
2. 创建功能分支
3. 进行代码修改
4. 如果适用，添加测试
5. 提交 Pull Request

## 📄 许可证

本项目用于教育和学习目的。

## 🆘 技术支持

如有问题或建议，请在仓库中提交 Issue。

---

**注意**: 本项目仅供学习和研究使用，请勿用于商业用途。