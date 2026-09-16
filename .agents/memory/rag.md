---
id: rag
updated: 2026-09-16
---

# RAG

模式：chunked extractive RAG。不是向量检索，不是 LLM 生成。

## Pipeline

1. 启动时读 `catalog.json` + `faq.json` + `data/docs/*.md`
2. 切块（商品摘要/说明、FAQ 整条、Markdown 按标题与约 420 字）
3. 中文单字 + bigram 倒排，TF-IDF 变体，标题/小节加权
4. `POST /api/v1/ask` 取 top-k 块，`composeAnswer` 按意图模板拼接
5. 返回 `citations[]`（`kind`: product | faq | knowledge），`mode=chunked-extractive-rag`

SvelteKit：`$lib/server/rag.ts` 转发到 Go。`fetchCatalog()` 失败时回落本地 JSON，并设 `ragReady: false`。

前端跳转：`$lib/hit-links.ts` 把商品/代寄命中链到 `/shop/sent`。

## Knowledge docs

`services/rag/data/docs/` 现有 10 篇（入门、品相、欧洲散票、日本明信片、专题、新票旧票、保存、代寄、赝品、目录编号）。来自整理稿，不是网页抓取。对外按科普说明，非证书。

## When to change layer

- 同义词对不上、文档明显变长 → 再考虑 embedding
- 需要对话口吻 → 只换生成层，检索仍喂切块
- 用户明确要 Go，不要加 Python 服务
