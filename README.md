# Go 项目

一个综合性的 Go Web 应用程序，实现了用户认证、CRUD 操作以及集成阿里云 OSS 的文件管理功能。

## 功能特性

- **用户管理**

  - 用户注册和登录
  - JWT 令牌认证
  - 安全的密码处理

- **文件管理**
  - 集成阿里云 OSS 的文件上传
  - 文件类型和大小验证
  - 文件列表分页支持
  - 文件的增删改查操作
  - OSS 对象的自动清理

## 技术栈

- Go (后端)
- 阿里云 OSS (云存储)
- MySQL (数据库)

## 环境要求

- Go 1.16 或更高版本
- MySQL
- 阿里云 OSS 账号
- Git

## 安装步骤

1. **克隆仓库**

   ```bash
   git clone [仓库地址]
   cd go-project
   ```

2. **配置环境变量**
   在项目根目录创建 `.env` 文件，添加以下配置：

   ```env
   # 数据库配置
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database_name

   # 阿里云 OSS 配置
   OSS_ENDPOINT=your_oss_endpoint
   OSS_ACCESS_KEY_ID=your_access_key_id
   OSS_ACCESS_KEY_SECRET=your_access_key_secret
   OSS_BUCKET=your_bucket_name
   ```

3. **安装依赖**

   ```bash
   go mod tidy
   ```

4. **运行应用**
   ```bash
   go run main.go
   ```

## API 文档

### 认证接口

#### 用户注册

- **POST** `/api/auth/register`
- **请求体:**
  ```json
  {
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```

#### 用户登录

- **POST** `/api/auth/login`
- **请求体:**
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```

### 文件管理接口

#### 上传文件

- **POST** `/api/files/upload`
- **请求头:** `Authorization: Bearer {token}`
- **请求体:** `multipart/form-data`
  - `file`: 要上传的文件

#### 文件列表

- **GET** `/api/files`
- **请求头:** `Authorization: Bearer {token}`
- **查询参数:**
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10)

#### 获取文件详情

- **GET** `/api/files/{id}`
- **请求头:** `Authorization: Bearer {token}`

#### 删除文件

- **DELETE** `/api/files/{id}`
- **请求头:** `Authorization: Bearer {token}`

## 响应格式

所有 API 响应都遵循以下标准格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

## 错误处理

API 使用标准的 HTTP 状态码，错误响应格式如下：

```json
{
  "code": 400,
  "message": "错误描述",
  "data": null
}
```

## 参与贡献

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m '添加一些很棒的功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 开源协议

本项目基于 MIT 协议开源 - 查看 LICENSE 文件了解详情。
