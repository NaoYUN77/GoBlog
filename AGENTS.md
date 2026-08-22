# AGENTS.md

个人博客项目：Go + Gin 后端（模块名 `Blog`）+ Vue 3 管理后台前端（`admin-web/`）。后端提供 `/admin` 系列 API，前端是文章管理后台（登录、文章 CRUD、Markdown 编辑器）。

## 技术栈

- **后端**：Go 1.25.6、Gin、sqlx + MySQL、go-redis/v9、golang-jwt/v5（HS256）、viper、zap、bwmarrin/snowflake
- **前端**：Vue 3（`<script setup lang="ts">`）+ TypeScript + Vite + Tailwind CSS v4 + Pinia + vue-router + axios + md-editor-v3 + json-bigint

## 目录结构

| 路径 | 说明 |
|---|---|
| `main.go` | 后端入口：依次初始化 settings(viper) → snowflake → MySQL → Redis → zap，启动 HTTP 服务，支持 5s 优雅停机 |
| `settings/` | viper 读取 `config.yaml`（app/mysql/redis 配置），暴露全局 `Cnf` |
| `router/` | Gin 路由，全部挂在 `/admin` 前缀下；注册 zap 请求日志中间件 |
| `handler/` | HTTP 层：参数绑定、调用 service、返回 `gin.H` |
| `service/` | 业务层：组装 db 结构体、生成雪花 ID |
| `repository/` | 数据层：`database/`（sqlx）、`redis/`（go-redis，仅初始化，尚未投入使用） |
| `models/` | 结构体定义 + SQL 建表脚本（原目录名 `modles`，已重命名为 `models`） |
| `middleware/` | `ParseTokenMiddleware`（Bearer JWT）、`LoggerMiddleWare`（zap） |
| `pkg/` | `jwt`（生成/解析）、`snowflake`（文章 ID）、`logger`（zap，development/production 模式） |
| `admin-web/` | Vue 3 管理后台（见下方前端说明） |
| `web/` | Vue 3 用户端前台：首页文章列表 + slug 详情页，端口 5174，`/post` 前缀代理到 8080（见 `web/README.md`） |

## 常用命令

### 后端（仓库根目录）

```bash
go run .          # 运行（需 MySQL/Redis 可连，配置见 settings/config.yaml）
go build -o blog.exe .   # 构建（Windows）
```

- 默认监听端口 `8080`；配置：MySQL `root/123@127.0.0.1:3306/blog`，Redis `127.0.0.1:6379`（密码 `123456`）
- 数据库建表脚本在 `modles/*.sql`：`tb_post`（文章表）、`tb_user`（管理员表，默认 `admin/admin123`）
- 依赖安装：`go mod tidy`

### 前端（admin-web/）

```bash
npm install
npm run dev      # 开发服务器 http://localhost:5173，/admin 代理到 http://127.0.0.1:8080（vite.config.ts）
npm run build    # vue-tsc 类型检查 + vite 构建，产物 dist/
npm run preview
```

## API 约定（后端）

所有接口前缀 `/admin`，除 `login` 外均需请求头 `Authorization: Bearer <token>`（10 分钟有效期）。

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/admin/login` | 登录，返回 `{code: 200, token}` |
| POST | `/admin/post` | 创建文章 |
| PATCH | `/admin/patch/:id` | 部分更新文章（仅更新非 nil 字段，动态拼 SQL） |
| DELETE | `/admin/delete/:id` | 删除文章 |
| GET | `/admin/get/:id` | 文章详情 |
| GET | `/admin/list` | 分页列表，query：`page`、`page_size`（默认 1/10，上限 100，按 `created_at desc`） |

- JWT：HS256，盐 `lolololo`，有效期 10 分钟（`models/jwt.go`）
- 默认管理员：admin / admin123（`models/admin.sql`）
- 文章 ID 由雪花算法生成（int64），超出 JS Number 安全整数范围
- 响应统一 `gin.H`：`{code?, msg?, data?, page?, page_size?}`；错误时 HTTP 状态码多为 400
- 分层调用约定：handler → service → repository，禁止跨层；repository 返回中文 `errors.New`，service 透传，handler 拼装响应

## 用户端接口（无鉴权）

面向访客公开，不需要 token。只返回 **已发布（status=1）** 的文章：

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/post/:slug` | 通过 slug 打开文章详情，字段 snake_case（含 `category_id`、`view_count`），每次访问 view_count +1 |
| GET | `/posts` | 已发布文章分页列表，query：`page`、`page_size`（默认 1/10，上限 100，按 `created_at desc`），字段风格同 `/admin/list`（大驼峰） |

- 草稿/不存在的 slug 返回 404
- 实现：`handler/post.go` 的 `UserGetPostHandler`/`UserGetPostListHandler` → `service/get.go` 的 `GetPostBySlug`/`GetPublishedPostList` → `repository/database/get.go` 对应查询 + `IncrementViewCount`

## 前后端字段兼容（重要）

后端两个读取接口的字段风格不一致，前端统一在 `admin-web/src/utils/normalizePost.ts` 归一化，业务代码只使用 `types/post.ts` 的 `Post` 类型：

- `GET /admin/list` 直接序列化 `Postdb`（结构体无 json tag）→ 字段为大驼峰：`ID`、`Title`、`CoverURL`、`CreatedAt`…
- `GET /admin/get/:id` 手动拼装 `gin.H` → 字段为 snake_case，且时间字段是 `create_at`/`update_at`（**少一个 d**），并缺少 `category_id`、`view_count`

改动后端结构体字段/命名时，必须同步检查 `normalizePost.ts` 与 `types/post.ts`。

## 已知问题（改动前先了解，勿当作新鲜事顺手"修复"）

- `service/post.go` 的 `CreatePost` **未透传 `CategoryID`**，请求体带 `category_id` 也不会写入；`UpdatePost` 同样把 `upr.CategoryID` 丢弃了（`database.UpdatePost` 本身支持该字段）
- `repository/database/DeletePost.go` 执行 `delete from post`，但表名实际是 `tb_post`（会报错）
- `handler/post.go` 的 `DeletePostHandler` 在 `ParseInt` 失败分支缺少 `return`，会带着 `id=0` 继续执行
- `handler/post.go` 的 `GetPostSigleHandler` 响应缺 `category_id`、`view_count`（前端已在 `normalizePost.ts` 兼容）
- `repository/redis/redis.go` 的 `Init` 中 `rdb := redis.NewClient(...)` 遮蔽了包级变量 `rdb`，全局变量仍为 nil，后续 `Close()` 也会 panic
- `modles/post.go` 注释说明"每 get 一次 view_count 加 1"，但当前实现未增加阅读量
- 后端未配置 CORS、也未托管前端静态文件（生产部署方案见 `admin-web/README.md`）

## 前端约定

- 组合式 API + `<script setup lang="ts">`，组件使用 `.vue` 单文件
- 状态管理：Pinia（`stores/`），auth 状态存 localStorage，key 为 `blog-admin-auth`
- 请求统一走 `src/api/http.ts`：baseURL `/admin`，响应用 `json-bigint` 解析（防雪花 ID 精度丢失），401/鉴权类错误自动登出并跳转登录页
- 路由守卫在 `src/router/index.ts`：非公开路由需登录，本地按 10 分钟 TTL 提前判过期
- API 接口封装在 `src/api/`（`auth.ts`、`posts.ts`），页面只调用 api 层，不直接发请求
- 构建前必须通过 `vue-tsc` 类型检查（`npm run build` 内置）
