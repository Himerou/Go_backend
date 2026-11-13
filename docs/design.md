# 项目设计文档

## 1. 项目概述

GO_web 是一个基于 Go 和 MySQL 实现的 Todo Web 应用，提供用户注册、登录、个人资料管理以及 Todo 任务管理等功能。

项目采用分层架构，将 HTTP 请求处理、业务逻辑和数据访问分离，并通过接口组织各层之间的依赖关系。

## 2. 项目架构

项目主要分为以下几个部分：

| 模块       | 职责                                  |
| ---------- | ------------------------------------- |
| Router     | 注册路由、组织中间件                  |
| Handler    | 处理 HTTP 请求、解析参数、构造响应    |
| Service    | 实现业务逻辑和业务规则                |
| Repository | 负责数据库访问和 Session 存储         |
| Model      | 定义数据结构、上下文键和业务错误      |
| Web        | 提供静态 HTML、CSS 和 JavaScript 页面 |

整体请求处理流程：

```
Web 前端
   |
   v
Router / Middleware
   |
   v
Handler
   |
   v
Service
   |
   v
Repository
   |
   v
MySQL
```

## 3. 分层设计

### 3.1 Handler 层

Handler 层负责 HTTP 请求和响应处理。

主要职责：

- 接收并解析 HTTP 请求。
- 读取路径参数、请求体和请求上下文。
- 调用 Service 层完成业务操作。
- 根据执行结果返回相应的 HTTP 状态码和 JSON 数据。

Handler 不直接访问数据库，而是依赖 Service 接口。

### 3.2 Service 层

Service 层负责业务逻辑。

主要职责：

- 实现用户注册、登录、登出等认证逻辑。
- 实现用户资料管理逻辑。
- 实现 Todo 的增删改查。
- 处理业务规则和业务错误。
- 协调 Repository 层完成数据操作。

Service 通过 Repository 接口访问数据，不直接依赖具体的数据存储实现。

### 3.3 Repository 层

Repository 层负责数据访问。

主要职责：

- 执行 MySQL 数据库操作。
- 完成用户和 Todo 数据的增删改查。
- 管理 Session 数据的存储和删除。

Repository 层隐藏具体的数据访问细节，为 Service 层提供统一的数据操作接口。

### 3.4 接口与依赖关系

项目采用面向接口编程的方式组织层间依赖。

- Handler 所需的接口定义在 Handler 包中，由 Service 层实现。
- Service 所需的接口定义在 Service 包中，由 Repository 层实现。
- 具体实现通过程序入口进行初始化和依赖注入。

这种设计使各层不必依赖其他层的具体实现，有利于降低耦合度，也便于使用 Mock 对象进行单元测试。

## 4. 用户认证设计

项目采用 Session + Cookie 机制维护用户登录状态。

### 4.1 登录流程

1. 用户提交用户名和密码。
2. Handler 接收请求并调用认证 Service。
3. Service 查询用户信息，并使用 bcrypt 校验密码。
4. 验证成功后创建 Session。
5. 服务端将 Session ID 通过 `session_id` Cookie 返回给客户端。

### 4.2 认证中间件

访问需要登录的接口时，中间件执行以下操作：

1. 从请求 Cookie 中读取 `session_id`。
2. 根据 Session ID 查询 Session 存储。
3. 验证 Session 是否存在且有效。
4. 验证通过后，将对应的用户 ID 写入请求上下文。
5. 将请求交给后续 Handler 处理。

如果认证失败：

- API 请求返回 `401 Unauthorized`。
- 页面请求重定向到首页，并附带未登录提示参数。

### 4.3 登出流程

1. 客户端发送登出请求。
2. 服务端根据 Cookie 获取 Session ID。
3. 删除服务端对应的 Session。
4. 返回过期 Cookie，使客户端清除登录凭证。

## 5. 数据库设计

项目使用 MySQL 存储用户和 Todo 数据。

主要包含两张数据表：

### 5.1 users 表

| 字段          | 说明                 |
| ------------- | -------------------- |
| id            | 用户 ID，主键        |
| username      | 用户名，具有唯一约束 |
| password_hash | bcrypt 密码哈希      |
| profile       | 个人资料             |
| created_at    | 创建时间             |

### 5.2 todos 表

| 字段       | 说明          |
| ---------- | ------------- |
| id         | Todo ID，主键 |
| title      | Todo 标题     |
| content    | Todo 内容     |
| completed  | 是否完成      |
| user_id    | 所属用户 ID   |
| created_at | 创建时间      |
| updated_at | 更新时间      |

### 5.3 数据关联

`todos.user_id` 是外键，关联 `users.id`。

当用户被删除时，通过 `ON DELETE CASCADE` 自动删除该用户关联的 Todo 数据。

## 6. 用户数据隔离

项目通过用户身份和数据库查询条件实现 Todo 数据隔离。

认证中间件验证用户身份后，将用户 ID 写入请求上下文。

Service 层从请求上下文中获取当前用户 ID，并在相关数据操作中传递给 Repository 层。

Repository 层根据 Todo ID 和用户 ID 限制查询、更新和删除范围，避免用户访问或修改其他用户的数据。

## 7. 单元测试设计

项目使用 Go 标准库 `testing` 和 `net/http/httptest` 编写单元测试。

测试主要分为两个层次：

### 7.1 Handler 层测试

通过 Mock Service 隔离业务逻辑，重点测试：

- HTTP 请求解析。
- 参数和 JSON 格式校验。
- HTTP 状态码。
- 响应数据。
- Cookie 设置与处理。

### 7.2 Service 层测试

通过 Mock Repository 隔离数据库访问，重点测试：

- 业务逻辑是否正确。
- 不同业务条件下的处理结果。
- 错误传递和业务错误判断。

运行全部测试：

```shell
go test ./...
```