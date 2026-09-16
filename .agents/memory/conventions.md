---
id: conventions
updated: 2026-09-16
---

# Conventions

- 新代码用 TypeScript / Go；不要为 RAG 引入 Python
- 文件改动保持小而贴合现有命名：`+page.svelte`、`+page.server.ts`、`+server.ts`
- 服务端 IO 集中在 `$lib/server/`；纯函数可放 `$lib/`
- 不要为了类型去 import `./$types`（生成文件可能不在仓库）
- 不要加版权头、不要无关重构、不要新加格式化工具
- 组件文案用中文；id、kind、环境变量用英文 kebab / snake
- 测试：Go 用 `engine_test.go`；前端目前无单测，不要凭空搭测试框架

## Catalog edits

改一件商品：只改 `catalog.json` 对应对象。商店卡片、代寄加购、RAG 商品块都会跟。

## Inquiry kinds

`contact` | `forward` | `route-test`。代寄必须有地址和服务 id。
