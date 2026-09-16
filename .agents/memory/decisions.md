---
id: decisions
updated: 2026-09-16
---

# Decisions

1. **RAG 用 Go，不用 Python**  
   用户明确选型。生态默认 FastAPI 被否决。后续生成层也在 Go 里加。

2. **认证留在 Node/SQLite**  
   账号与问询和页面在一起。Go 只做检索。

3. **先切块抽取，后向量 / LLM**  
   语料小、商品边界清晰。知识库用 Markdown 切块 + 倒排。未批准前不上 embedding。

4. **知识来自整理稿，不爬网**  
   用户停止网页搜索方案。`data/docs` 按集邮常识撰写，需标明非鉴定。

5. **申请不是订单**  
   代寄/咨询写入 `inquiries`，状态 `received`。支付、扣库存、后台审单未做。

6. **目录单一来源**  
   `catalog.json` 同时喂 Go 与 SvelteKit 回落，避免商店价和 RAG 价两套。

7. **本地双进程优先于容器**  
   用户曾要求不用 Docker 启动。Compose 文件保留给以后部署；当前验证用 `go run` + `npm run dev`。
