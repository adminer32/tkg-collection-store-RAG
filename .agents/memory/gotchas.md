---
id: gotchas
updated: 2026-09-16
---

# Gotchas

- **客户端 / 服务端拆分**：`$lib/server/*` 进 `.svelte` 会在构建期炸掉。类型要拆到 `$lib/rag-types.ts` 等。
- **`fetchCatalog` 形状**：现返回 `{ items, ragReady }`，不再是数组。`src/routes/api/rag/catalog/+server.ts` 若仍 `json({ items })` 且把整个对象当 items，需要按新形状改。
- **JWT 默认值**：`src/lib/server/auth.ts` 仍有硬编码 fallback。上线必须用 `JWT_SECRET`。`.env.example` 写的是 `change-me`。
- **Vite on Windows**：文件监视可能 `spawn EPERM` / `EBUSY`。`vite.config.ts` 已 `server.watch.usePolling`。沙箱里启动 dev 往往需要能 spawn 子进程。
- **Go 测试缓存**：用户目录可能 Access is denied。把 `GOCACHE`/`GOMODCACHE` 指到 `services/rag/.gocache` 等（已 gitignore）。
- **本机未必有 docker/podman**：不要把「能 compose 起来」当成已验证。
- **首页不是关于页**：`src/routes/+page.svelte` 是店首页；关于在 `/about`。
- **登出双通道**：页面会话以 localStorage 为准；logout API 只清 cookie，调用它是为了对称，不能单靠 cookie。
- **RAG 离线**：商店仍应能列出 `catalog.json`；助手提问会 502，页面要能显示 `ragReady`。
