# Zero-Life

Zero-Life 是 [Firefly III](https://www.firefly-iii.org/) 的 Go + Vue 重写版本，一个开源的个人财务管理系统。

## 技术栈

| 层 | 技术 |
|---|------|
| 后端 | Go 1.22+、Gin、GORM、SQLite、Redis、JWT |
| 前端 | Vue 3、TypeScript 6.0、Element Plus、Pinia、vue-i18n、ECharts、Vite |
| 部署 | Docker、Docker Compose、Nginx |

## 功能模块

- **账户管理** — 资产、负债、支出、收入账户，虚拟账户支持
- **交易管理** — 收入/支出/转账，拆分交易，批量操作，高级搜索
- **预算管理** — 月度/年度预算，分类关联，使用率追踪与预警
- **分类与标签** — 树形分类，彩色标签，交易归类
- **账单管理** — 定期账单追踪，到期提醒
- **存钱罐** — 目标储蓄，存取记录，进度追踪
- **规则引擎** — 触发条件 + 自动动作，规则组管理
- **定期交易** — 循环交易自动执行，到期处理
- **报表** — 收支报表、分类报表、预算报表、净值趋势、标签报表、审计报表
- **数据洞察** — 收入/支出洞察，按分类/账户/日期分析
- **导入导出** — CSV 导入（字段映射），CSV/JSON 导出
- **对账** — 账户对账，余额差异计算
- **货币与汇率** — 多货币支持，汇率管理
- **MFA 两步验证** — TOTP 验证，备用恢复码
- **Webhook** — 交易事件通知，HMAC 签名验证
- **用户管理** — 注册/登录，角色权限，账户锁定
- **系统管理** — 系统配置，定时任务，邮件测试

## 项目结构

```
zero-life/
├── server/                  # Go 后端
│   ├── cmd/server/          # 入口 main.go
│   ├── internal/
│   │   ├── config/          # 配置加载
│   │   ├── controller/      # HTTP 控制器（34 个）
│   │   ├── dto/             # 请求/响应 DTO
│   │   ├── middleware/       # 中间件（认证、CORS、日志等）
│   │   ├── model/           # 数据模型（22 个）
│   │   ├── pkg/             # 工具包（errcode、jwt、pagination）
│   │   ├── repository/      # 数据访问层（23 个）
│   │   ├── router/          # 路由注册
│   │   └── service/         # 业务逻辑层（31 个）
│   ├── config.yaml          # 配置文件
│   └── Dockerfile
├── web/                     # Vue 前端
│   ├── src/
│   │   ├── api/             # API 请求封装（25 个）
│   │   ├── components/      # 组件（图表、表单、布局）
│   │   ├── i18n/            # 国际化（中/英）
│   │   ├── pages/           # 页面（35 个）
│   │   ├── router/          # 路由配置
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── types/           # TypeScript 类型定义
│   │   └── utils/           # 工具函数
│   ├── nginx.conf           # Nginx 配置
│   └── Dockerfile
├── docker-compose.yml       # Docker Compose 编排
└── firefly-iii/             # PHP 参考实现（仅供参考）
```

## 快速开始

### 环境要求

- **Go** 1.22+
- **Node.js** 20+
- **Redis** 7.0+
- **SQLite**（Go 内置驱动，无需单独安装）

---

### 方式一：本地开发启动

#### 1. 启动 Redis

```bash
# 方式 A：本地安装 Redis 后启动
redis-server

# 方式 B：用 Docker 快速启动 Redis
docker run -d --name redis -p 6379:6379 redis:7.0-alpine
```

#### 2. 启动后端

```bash
cd server

# 安装 Go 依赖
go mod download

# 启动后端服务（默认监听 8080 端口）
go run ./cmd/server/
```

后端启动后会：
- 读取 `server/config.yaml` 配置文件
- 自动创建 SQLite 数据库文件 `server/data/zero-life.db`
- 自动执行数据库迁移（建表）
- 初始化默认货币（CNY、USD、EUR、JPY、GBP）
- 在 `http://localhost:8080` 提供服务

**配置说明**（`server/config.yaml`）：

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `app.port` | 后端监听端口 | `8080` |
| `db.path` | SQLite 数据库路径 | `./data/zero-life.db` |
| `redis.host` | Redis 地址 | `localhost` |
| `redis.port` | Redis 端口 | `6379` |
| `jwt.secret` | JWT 签名密钥（**生产环境务必修改**） | `change-me-in-production` |
| `cors.allow_origins` | 允许的前端跨域来源 | `http://localhost:5173` |
| `attach.path` | 附件存储路径 | `./data/attachments` |

#### 3. 启动前端

```bash
cd web

# 安装前端依赖
npm install

# 启动开发服务器（默认监听 5173 端口，自动代理 API 到后端）
npm run dev
```

前端开发服务器配置了 Vite 代理，所有 `/api/*` 请求会自动转发到 `http://localhost:8080`。

访问 `http://localhost:5173` 即可使用。

#### 4. 构建生产版本

```bash
# 构建后端
cd server
go build -o server ./cmd/server/

# 构建前端
cd web
npm run build
# 产物在 web/dist/ 目录
```

---

### 方式二：Docker Compose 一键启动（推荐）

这是最简单的方式，一条命令启动所有服务（Redis + 后端 + 前端）。

#### 1. 修改 JWT 密钥（重要）

```bash
# 在项目根目录创建 .env 文件
echo 'JWT_SECRET=your-secure-random-secret-here' > .env
```

#### 2. 一键启动

```bash
docker compose up -d
```

启动后包含三个服务：

| 服务 | 端口 | 说明 |
|------|------|------|
| `redis` | 6379 | Redis 缓存 |
| `api` | 8080 | Go 后端 API |
| `web` | 80 | Nginx 前端（反向代理 API） |

访问 `http://localhost` 即可使用。

#### 3. 常用命令

```bash
# 查看日志
docker compose logs -f

# 只看后端日志
docker compose logs -f api

# 停止所有服务
docker compose down

# 停止并删除数据卷（清空数据库）
docker compose down -v

# 重新构建并启动（代码更新后）
docker compose up -d --build
```

#### 4. 数据持久化

Docker Compose 使用了两个命名卷来持久化数据：

- `sqlite_data` — SQLite 数据库文件
- `redis_data` — Redis 数据

数据不会因容器重启而丢失，除非执行 `docker compose down -v`。

---

### 方式三：单独用 Docker 构建镜像

#### 构建后端镜像

```bash
cd server
docker build -t zero-life-api .

# 运行（需要先启动 Redis）
docker run -d \
  --name zero-life-api \
  -p 8080:8080 \
  -e REDIS_HOST=host.docker.internal \
  -e JWT_SECRET=your-secret \
  -v zero-life-data:/app/data \
  zero-life-api
```

#### 构建前端镜像

```bash
cd web
docker build -t zero-life-web .

# 运行
docker run -d \
  --name zero-life-web \
  -p 80:80 \
  zero-life-web
```

> 注意：前端 Nginx 配置中 API 代理地址为 `http://api:8080`，单独运行时需修改 `web/nginx.conf` 中的 `proxy_pass` 为实际后端地址。

---

## API 概览

后端提供 RESTful API，基础路径为 `/api/v1/`，主要端点：

| 模块 | 端点前缀 | 说明 |
|------|----------|------|
| 认证 | `/api/v1/auth/` | 注册、登录、刷新令牌、忘记密码 |
| 账户 | `/api/v1/accounts/` | 账户 CRUD |
| 交易 | `/api/v1/transactions/` | 交易 CRUD、搜索、拆分、批量操作 |
| 预算 | `/api/v1/budgets/` | 预算 CRUD、使用率、历史 |
| 分类 | `/api/v1/categories/` | 分类 CRUD（树形） |
| 标签 | `/api/v1/tags/` | 标签 CRUD |
| 账单 | `/api/v1/bills/` | 账单 CRUD、到期提醒 |
| 存钱罐 | `/api/v1/piggy-banks/` | 存钱罐 CRUD、存取款 |
| 规则 | `/api/v1/rules/` | 规则 CRUD、执行 |
| 定期交易 | `/api/v1/recurring-transactions/` | 循环交易 CRUD、到期处理 |
| 报表 | `/api/v1/reports/` | 收支、分类、预算、净值、趋势、标签报表 |
| 导入 | `/api/v1/import/` | CSV 上传、解析、执行 |
| 导出 | `/api/v1/export/` | CSV/JSON 导出 |
| 货币 | `/api/v1/currencies/` | 货币管理、汇率 |
| 仪表盘 | `/api/v1/dashboard/` | 月度概览、预算预警、账单提醒 |
| 用户管理 | `/api/v1/admin/users/` | 用户列表、锁定、角色变更 |

所有需要认证的接口需在请求头中携带 `Authorization: Bearer <token>`。

## 开发

```bash
# 后端代码检查
cd server
go vet ./...

# 前端类型检查 + 构建
cd web
npm run build
```

## License

MIT
