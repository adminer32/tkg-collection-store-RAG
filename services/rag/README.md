# tkg-rag

高木桑收藏品店的 Go RAG 后端。SvelteKit 继续负责页面和账号；本服务检索商品、FAQ 与 `data/docs/` 里的集邮知识，并生成带引用的回答。

当前模式是 **chunked extractive RAG**：

1. 商品拆成摘要块 + 说明块，FAQ 整条一块，Markdown 按 `##` 标题和约 420 字切块（带一句 overlap）。
2. 对切块建中文 bigram 倒排索引（TF-IDF）。
3. 问句检索 top-k 切块，用模板抽取拼接，不调用外部 LLM。

## 本地运行

```bash
cd services/rag
go test ./...
go run .
```

默认监听 `:8081`。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/healthz` | 健康检查（含 docs/chunks 数量） |
| GET | `/api/v1/catalog` | 商品目录 |
| GET | `/api/v1/search?q=` | 切块检索 |
| POST | `/api/v1/ask` | 切块抽取式问答 |

环境变量：

- `RAG_ADDR` 默认 `:8081`
- `RAG_DATA_DIR` 默认 `./data`

知识文档放在 `data/docs/*.md`。SvelteKit 通过 `RAG_URL`（默认 `http://127.0.0.1:8081`）反代到 `/api/rag/*`。
