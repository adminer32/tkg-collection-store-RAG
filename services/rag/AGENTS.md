# services/rag/AGENTS.md

作用域：`services/rag/`。Go 检索服务，不负责登录或 SQLite。

## Runtime

- 入口 `main.go`，引擎 `engine.go`，切块 `chunk.go`
- 默认 `RAG_ADDR=:8081`，`RAG_DATA_DIR=./data`
- 接口：`GET /healthz`、`GET /api/v1/catalog`、`GET /api/v1/search`、`POST /api/v1/ask`
- 回答 `mode`：`chunked-extractive-rag`

## Data

改商品只动 `data/catalog.json`。FAQ 在 `data/faq.json`。长知识在 `data/docs/*.md`。

切块规则（`chunk.go`）：

- 商品：`#summary` + `#description`
- FAQ：整条 `#body`
- Markdown：按 `##`，再按约 420 字，块间一句 overlap

不要在 Go 里再维护一份价格表。

## Retrieval

- 中文单字 + bigram，标题/小节权重高于正文
- 无 embedding、无 LLM；要换生成层时保留切块与 citations
- 改分词或切块后：`go test ./...`
- Windows 沙箱可设 `GOCACHE` / `GOMODCACHE` 到本目录，`GOTELEMETRY=off`

## Don't

- 不要引入 Python 依赖或爬虫入库流程（用户已否决网页抓取）
- 不要在回答里编造库存或证书级鉴定
- 不要把知识库原文当成商品报价
