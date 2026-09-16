# src/AGENTS.md

作用域：`src/`。SvelteKit 应用约定。跨进程 / 商品规则看根 `AGENTS.md` 与 `.agents/memory/`。

## Svelte 5

- 页面用 `$props()` / `$state()` / `$derived()`，不要把新页写成 Svelte 4 `export let`
- 旧文件（`Navbar.svelte`、`auth/+page.svelte`）可暂留 store 语法，改到再迁
- 客户端禁止 `import` `$lib/server/*`（会打进浏览器包）
- 跨端类型放 `$lib/rag-types.ts`、`$lib/inquiry-types.ts`、`$lib/catalog-labels.ts`

## Routes

| 路径 | 职责 |
| --- | --- |
| `/` | 店铺首页 + 精选目录 |
| `/about` | 关于我们（不是首页） |
| `/shop` | 目录 + 店员助手 |
| `/shop/sent` | 代寄 / 邮路测试表单 |
| `/about/contact` | 咨询表单 |
| `/userspace` | 登录后的申请与提问 |
| `/auth` | 登录注册 |

查询参数：`?service=postcard-forward|mail-route-test`，`?addon=<productId>`。

## Server modules

- 目录读取：`$lib/server/catalog.ts`（读同一份 `catalog.json`）
- RAG 反代：`$lib/server/rag.ts`（`RAG_URL`，失败时目录回落到本地 JSON）
- 认证：`$lib/server/auth.ts`（Bearer JWT）
- 问询：`$lib/server/inquiries.ts` + `db.ts` 表 `inquiries` / `rag_questions`
- 助手跳转：`$lib/hit-links.ts`

API：`/api/auth/*`、`/api/inquiries`、`/api/me`、`/api/rag/ask`、`/api/rag/catalog`。

问询和提问在浏览器带 `Authorization: Bearer <token>`；`auth` store 的 token 在 `localStorage`。

## UI

- Tailwind v4：`src/app.css` 只有 `@import "tailwindcss"`，不要加 v3 `tailwind.config`
- 新页沿用现有玫瑰 / 紫色圆角卡片，不要另起设计系统
- 文案：代寄、邮路、中国邮政、不承诺升值
