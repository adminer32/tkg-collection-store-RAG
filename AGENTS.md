# AGENTS.md

高木桑收藏品店（`tkg-collection-store`）。先读本文件，再按任务打开 `.agents/memory/` 里对应条目，不要一次读完所有记忆。

## Role

你是本仓库的实现代理：改代码、补页面、改 RAG、跑本地验证。默认中文沟通，标识符与路径保持英文。

## Stack (always on)

- 前端 / 账号：SvelteKit 2 + Svelte 5 + Vite 8 + Tailwind v4
- 认证库：SQLite（`better-sqlite3`）+ JWT + bcrypt
- 目录 / RAG：Go sidecar `services/rag`（extractive，无 LLM）
- 实寄承诺：中国邮政；不做商业快递替代

## Boundary

| 放 Node | 放 Go | 不要做 |
| --- | --- | --- |
| 页面、登录注册、问询/申请、RAG 反代 | 目录检索、切块、问答生成 | 把认证迁出 Node |
| `inquiries` / `rag_questions` | `catalog.json` / `faq.json` / `data/docs` | 新增 Python RAG |
| JWT 校验（`Authorization: Bearer`） | 倒排索引 + 抽取回答 | 未要求就上 embedding / 云模型 |

## Load memory by task

先打开 `.agents/memory/INDEX.md`，只读命中的文件：

- 改文案、品类、代寄规则 → `product.md`
- 改路由、API、进程边界 → `architecture.md`
- 改检索 / 知识库 / 切块 → `rag.md`
- 改 Svelte / Go 代码风格 → `conventions.md`
- 排障、已知坑 → `gotchas.md`
- 为什么这样实现 → `decisions.md`

目录级规则（比本文件更具体）：

- `src/AGENTS.md` — SvelteKit / 客户端边界
- `services/rag/AGENTS.md` — Go RAG 服务

## Do

- 商品价格、库存、标题只改 `services/rag/data/catalog.json`（商店、代寄、RAG 共用）
- 新集邮知识写成 `services/rag/data/docs/*.md`，按 `##` 分节，便于切块
- 服务端模块只放 `$lib/server/*`；客户端可共用类型放 `$lib/*.ts`
- 代寄/咨询走 `/api/inquiries`，不要先做支付和库存扣减
- 本地默认：`services/rag` 里 `go run .`，仓库根 `npm run dev`（`RAG_URL=http://127.0.0.1:8081`）
- 改 Go 检索后跑 `go test ./...`（在 `services/rag`，可设 `GOCACHE` 到该目录）
- 用户说用 Podman / 不要 Docker 时，不要改去调用 `docker`

## Don't

- 不要从 `$lib/server` 向 `.svelte` 页面 import
- 不要把 RAG 回答写成无来源的 LLM 闲聊；保持 citations
- 不要把 `/` 做成关于页；`/` 是店铺首页，`/about` 才是关于我们
- 不要提交 `.env`、`local.db`、`node_modules`、Go 缓存
- 不要新开替代 Web GUI / 第二套 Vite 入口来「修复当前页」
- 不要扩大范围：没要求就不上向量库、Ollama、支付、后台审单

## Verify

改动触及的最小验证即可：

1. 相关页面能打开（`/` `/shop` `/shop/sent` `/about/contact` `/userspace`）
2. 若改 RAG：`GET /healthz`，必要时 `POST /api/v1/ask`
3. 若改问询：未登录可提交，登录后带 `Authorization` 并出现在 `/userspace`

## Memory hygiene

- 新的稳定事实写入 `.agents/memory/` 对应文件，并在 `INDEX.md` 加一行
- 记忆只记结论、约束、路径；不贴日志、不贴对话、不贴大段代码
- 过时条目直接改文件，不要另开 `memory-v2`
- 会话临时笔记可放 `.agents/scratch/`（git 忽略），不要写进 memory
