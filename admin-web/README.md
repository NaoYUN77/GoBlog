# 笔耕 · 博客管理后台

基于 `Blog` 项目 `/admin` 系列接口实现的博客文章管理前端（Vue 3 + TypeScript + Vite + Tailwind CSS v4）。

## 启动

```bash
npm install
npm run dev
```

开发服务器默认运行在 `http://localhost:5173`，并将 `/admin/*` 请求代理到 `http://127.0.0.1:8080`（见 `vite.config.ts`），启动前请确保后端服务已运行、MySQL/Redis 可连接。

```bash
npm run build    # 类型检查 + 生产构建，产物在 dist/
npm run preview  # 预览生产构建
```

## 生产部署提示

后端当前未启用 CORS，也没有静态文件托管。上线时二选一：

1. 在 Gin 中用 `r.Static` / `r.StaticFile` 直接托管 `dist/` 产物，前后端同源，无需 CORS。
2. 或者给后端加 CORS 中间件，前端单独部署（Nginx/CDN），并把 `axios` 的 `baseURL` 改成后端完整地址。

## 功能范围

- 登录（`POST /admin/login`），JWT 存储在浏览器本地，顶部有实时倒计时提示（后端 token 有效期 10 分钟）。
- 文章列表（`GET /admin/list`）：因接口不返回总数，采用「按 100 条/页分批拉取 + 前端搜索/筛选/排序/加载更多」的方式。
- 新建文章（`POST /admin/post`）、编辑文章（`PATCH /admin/patch/:id`，仅提交实际改动字段）、删除文章（`DELETE /admin/delete/:id`）、文章详情（`GET /admin/get/:id`）。
- Markdown 编辑器（md-editor-v3），实时字数与阅读时长估算。

## 关于后端返回数据的兼容处理

- `GET /admin/list` 直接序列化 `Postdb`（无 json tag），字段是大驼峰（`ID`/`CoverURL`/`CreatedAt`…）；`GET /admin/get/:id` 是手动拼装的 `gin.H`，字段是 `snake_case`，且时间字段是 `create_at`/`update_at`（少一个 `d`）。前端在 `src/utils/normalizePost.ts` 里统一做了兼容归一化。
- 文章 id 由雪花算法生成，超过 JS `Number` 安全整数范围，响应体用 `json-bigint` 解析，避免大整数精度丢失（`src/api/http.ts`）。
- `GET /admin/get/:id` 的响应里没有 `category_id` 字段。编辑页据此只在用户主动修改「分类 ID」时才会把该字段放进 `PATCH` 请求体，避免把后端已有的分类误清空。

## 建议关注的后端问题（未在本次改动中修复）

在对接接口时发现以下几处后端行为可能不是预期效果，供参考：

1. `repository/database/CreatePost.go` 里 `INSERT` 语句用的列名是 `coverurl`，而表结构里实际列名是 `cover_url`，创建文章时如果不做修正大概率会插入失败。
2. `service/post.go` 的 `CreatePost` 没有把 `CategoryID` 透传给 `Postdb`，即便请求体带了 `category_id` 创建时也不会写入。
3. `handler/post.go` 的 `GetPostSigleHandler` 返回体里缺少 `category_id`、`view_count`，且时间字段名是 `create_at`/`update_at`（其余接口是 `created_at`/`updated_at`），建议统一。
