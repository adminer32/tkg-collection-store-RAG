---
id: architecture
updated: 2026-09-16
---

# Architecture

两进程：

```
browser → SvelteKit (:5173 dev / :3000 prod)
                ├─ SQLite  users / inquiries / rag_questions
                └─ HTTP    RAG_URL → Go :8081
```

Node 不实现检索；Go 不实现登录。

## Env

见 `.env.example`：`DATABASE_PATH`、`JWT_SECRET`、`RAG_URL`。
Go：`RAG_ADDR`、`RAG_DATA_DIR`。
Compose 里 web 的 `RAG_URL=http://rag:8081`。本地默认 `http://127.0.0.1:8081`。

## Auth

- 注册/登录：`/api/auth/register`、`/api/auth/login`
- 客户端：`src/lib/stores/auth.ts` 把 `{id, username}` 与 JWT 放进 `localStorage`
- 受保护读取：`Authorization: Bearer`，`getUserFromRequest`
- 登出：清 localStorage，并 POST `/api/auth/logout`（cookie 路径历史上不一致，以 localStorage 为准）
- JWT 用户 id 是 **string**（uuid），不要当成 number

## Persistence

`src/lib/server/db.ts` 启动时 `CREATE TABLE IF NOT EXISTS`：

- `users`
- `inquiries`（contact / forward / route-test）
- `rag_questions`（登录后提问才绑得上 user）

无迁移框架；改表结构时同时改这份 `exec` SQL，并考虑已有 `local.db`。

## Pages vs APIs

| 用户动作 | 页面 | API |
| --- | --- | --- |
| 逛店 / 提问 | `/` `/shop` | `/api/rag/ask` |
| 代寄 | `/shop/sent` | `POST /api/inquiries` |
| 咨询 | `/about/contact` | 同上，`kind=contact` |
| 看记录 | `/userspace` | `GET /api/me` |

## Deploy notes

`Dockerfile` 生产镜像需带上 `services/rag/data/catalog.json`，因为 Node 会本地回落读它。用户若要求 Podman，用 compose 文件但命令走 `podman`，不要假设本机有 `docker`。
