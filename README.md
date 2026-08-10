# AI 岗位分析助手

一个基于 Go 和 Vue 的前后端项目，用于根据岗位描述调用大模型生成求职分析，并将分析记录保存到 MySQL。项目适合作为 Go 后端实习简历项目，重点展示接口开发、数据库持久化、前后端联调和第三方大模型 API 接入能力。

## 功能介绍

- 岗位智能分析：输入岗位描述和问题，调用大模型生成中文分析结果。
- 历史记录保存：每次分析结果会保存到 MySQL 数据库。
- 分页查询历史：支持按页查询历史分析记录。
- 关键词搜索：支持按岗位描述、问题、回答内容进行模糊搜索。
- 删除历史记录：支持根据记录 ID 删除指定历史数据。
- 前端页面展示：使用 Vue3 实现岗位分析、历史查询、搜索和删除操作。

## 技术栈

### 后端

- Go
- go-zero / goctl
- GORM
- MySQL
- OpenAI 兼容格式的大模型 API

### 前端

- Vue 3
- Vite
- Axios

## 项目结构

```text
job-assistant
├─ etc
│  ├─ job-api.example.yaml    # 示例配置文件，可提交到 GitHub
│  └─ job-api.yaml            # 本地真实配置文件，不提交
├─ internal
│  ├─ config                  # 配置结构
│  ├─ handler                 # HTTP handler
│  ├─ llm                     # 大模型 API 调用封装
│  ├─ logic                   # 业务逻辑
│  ├─ model                   # GORM 数据模型
│  ├─ svc                     # 服务依赖初始化
│  └─ types                   # goctl 生成的请求/响应结构体
├─ web                        # Vue 前端项目
├─ job.api                    # go-zero API 描述文件
├─ job.go                     # 后端启动入口
├─ go.mod
└─ README.md
```

## 数据库准备

先创建数据库和数据表：

```sql
CREATE DATABASE job_analyzer
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE job_analyzer;

CREATE TABLE analyze_records (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    job_description LONGTEXT NOT NULL,
    question TEXT NOT NULL,
    answer LONGTEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 配置说明

复制示例配置文件：

```powershell
Copy-Item .\etc\job-api.example.yaml .\etc\job-api.yaml
```

然后修改 `etc/job-api.yaml`：

```yaml
Name: job-api
Host: 0.0.0.0
Port: 8888
Timeout: 180000

MySQL:
  DataSource: "root:你的MySQL密码@tcp(127.0.0.1:3306)/job_analyzer?charset=utf8mb4&parseTime=True&loc=Local"

LLM:
  BaseURL: "https://api.gemai.cc/v1"
  Model: "grok-3"
```

`etc/job-api.yaml` 是本地真实配置文件，可能包含数据库密码，所以不会提交到 GitHub。GitHub 中只保留 `etc/job-api.example.yaml` 示例配置。

## 启动后端

进入项目根目录：

```powershell
cd C:\Users\35102\Desktop\job-assistant
```

设置大模型 API Key：

```powershell
$env:LLM_API_KEY="your-api-key"
```

启动后端服务：

```powershell
go run . -f etc/job-api.yaml
```

启动成功后会看到：

```text
Starting server at 0.0.0.0:8888...
```

## 启动前端

打开另一个 PowerShell 窗口：

```powershell
cd C:\Users\35102\Desktop\job-assistant\web
npm install
npm run dev
```

浏览器访问：

```text
http://localhost:5173/
```

前端通过 Vite 代理访问后端接口，代理配置在 `web/vite.config.js`：

```js
server: {
  proxy: {
    '/api': 'http://127.0.0.1:8888',
  },
}
```

## 接口说明

### 1. 岗位分析

```http
POST /api/analyze
```

请求示例：

```json
{
  "jobDescription": "Go 后端开发，要求 MySQL 和 RESTful API",
  "question": "这个岗位需要哪些技能？"
}
```

返回示例：

```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "data": {
    "id": 1,
    "answer": "根据岗位描述，该岗位需要掌握 Go 基础、MySQL、RESTful API..."
  }
}
```

PowerShell 测试：

```powershell
$body = @{
    jobDescription = "Go 后端开发，要求 MySQL 和 RESTful API"
    question = "这个岗位需要哪些技能？"
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri "http://127.0.0.1:8888/api/analyze" `
    -Method Post `
    -ContentType "application/json; charset=utf-8" `
    -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
```

### 2. 查询历史记录

```http
GET /api/analyze/history?page=1&pageSize=10
```

参数说明：

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| page | int | 否 | 当前页，默认 1 |
| pageSize | int | 否 | 每页条数，默认 10，最大 50 |
| keyword | string | 否 | 搜索关键词 |

关键词搜索示例：

```http
GET /api/analyze/history?page=1&pageSize=10&keyword=Go
```

返回示例：

```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "pageSize": 10
  }
}
```

PowerShell 测试：

```powershell
Invoke-RestMethod `
  -Uri "http://127.0.0.1:8888/api/analyze/history?page=1&pageSize=10&keyword=Go" `
  -Method Get
```

### 3. 删除历史记录

```http
DELETE /api/analyze/history/:id
```

示例：

```http
DELETE /api/analyze/history/1
```

PowerShell 测试：

```powershell
Invoke-RestMethod `
  -Uri "http://127.0.0.1:8888/api/analyze/history/1" `
  -Method Delete
```

返回示例：

```json
{
  "code": 200,
  "message": "删除成功",
  "success": true
}
```

## 前端功能

前端页面主要包含：

- 岗位描述输入框
- 问题输入框
- 分析结果展示区
- 历史记录列表
- 关键词搜索框
- 删除历史记录按钮

前端请求后端接口时使用相对路径：

```js
axios.post('/api/analyze', data)
axios.get('/api/analyze/history', { params })
axios.delete(`/api/analyze/history/${id}`)
```

## 常见问题

### 1. npm 提示 EPERM 权限错误

如果 npm 报错路径在 `C:\Program Files\nodejs`，说明 npm 的缓存或全局目录配置到了系统目录。可以修改用户目录下的 `.npmrc`：

```ini
prefix=C:\Users\35102\AppData\Roaming\npm
cache=C:\Users\35102\AppData\Local\npm-cache
```

### 2. 页面一直显示“分析中”

优先检查后端是否启动：

```powershell
go run . -f etc/job-api.yaml
```

还要确认大模型 API Key 是否已设置：

```powershell
$env:LLM_API_KEY="your-api-key"
```

如果后端请求大模型响应较慢，前端会等待接口返回。

### 3. 浏览器打不开 localhost:5173

说明前端服务没有启动，重新执行：

```powershell
cd web
npm run dev
```

如果 Vite 显示的是 `5174` 或其他端口，要访问实际输出的地址。

### 4. 修改 job.api 后类型没有更新

保存 `job.api` 后重新执行：

```powershell
goctl api go --api .\job.api --dir .
```

goctl 默认不会覆盖已经存在的 logic 文件，所以业务逻辑文件通常需要自己手动修改。

## 项目亮点

- 使用 go-zero 根据 `job.api` 生成接口结构，代码组织清晰。
- 使用 GORM 操作 MySQL，实现分析记录持久化。
- 封装大模型 API 调用，支持 OpenAI 兼容接口。
- 实现分页、关键词搜索、删除等常见后端业务能力。
- 使用 Vue3 + Axios 完成前后端联调，具备完整项目展示效果。

## 后续优化方向

- 增加登录注册和用户维度的数据隔离。
- 增加 Markdown 渲染，让大模型回答展示更美观。
- 增加 Docker Compose，一键启动 MySQL 和后端服务。
- 增加 Swagger 或独立接口文档。
- 增加单元测试和接口测试。
