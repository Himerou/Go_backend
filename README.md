# GO_web

一个基于 Go 和 MySQL 实现的轻量级 Todo Web 应用，提供用户注册、登录、个人资料管理和个人任务管理功能。

后端采用 Handler、Service、Repository 分层架构，使用 Session + Cookie 实现用户身份认证。前端使用原生 HTML、CSS 和 JavaScript 实现。

## 技术栈

| 分类       | 技术                         |
| ---------- | ---------------------------- |
| 编程语言   | Go                           |
| Web 服务   | `net/http`                   |
| 数据库     | MySQL                        |
| 数据库访问 | `database/sql`               |
| MySQL 驱动 | `go-sql-driver/mysql`        |
| 身份认证   | Session + Cookie             |
| 密码加密   | bcrypt                       |
| 前端       | HTML、CSS、JavaScript        |
| 自动化测试 | Go testing、httptest、Apifox |

## 项目结构

```
GO_web/
├── database/       # 数据库初始化脚本
├── docs/           # 项目相关文档
├── handler/        # HTTP 请求处理
├── model/          # 数据模型
├── repository/     # 数据库访问与 Session 存储
├── router/         # 路由配置
├── service/        # 业务逻辑
├── web/            # 静态前端页面及资源
├── main.go         # 程序入口
├── go.mod
└── go.sum
```

## 环境要求

运行项目前需要准备：

- Go 1.26.0 或更高版本
- MySQL 8.0 或更高版本
- Git

## 项目运行

### 1. 克隆项目

```shell
git clone https://github.com/Himerou/Go_backend
cd GO_web
```

### 2. 下载依赖

```shell
go mod download
```

### 3. 初始化数据库

在项目根目录下执行：

```shell
sudo mysql < database/init.sql
```

初始化脚本会创建 `todo_app` 数据库及 `users`、`todos` 两张表。

两张表通过外键关联，每个任务属于一个用户。删除用户时，会级联删除该用户的所有任务。

### 4. 配置环境变量

在项目根目录创建 `.env` 文件：

```shell
cp .env.example .env
```

根据本地 MySQL 配置修改 `.env` 中的 DSN：

```
DSN=用户名:密码@tcp(127.0.0.1:3306)/todo_app?parseTime=true&loc=Local
```

### 5. 启动项目

```shell
go run .
```

浏览器访问：

http://localhost:8080

## 项目测试

在项目根目录执行：

```shell
go test ./...	
```

## 项目文档

- `docs/api.md`：API 接口文档
- `docs/design.md`：项目架构与设计说明
- `docs/tasks.md`：项目开发任务记录