# Agent memory

给编码代理用的分层记忆，不是给顾客看的文档。

```
AGENTS.md                 常驻：边界、禁令、如何加载记忆
src/AGENTS.md             仅 src/
services/rag/AGENTS.md    仅 Go RAG
.agents/memory/INDEX.md   记忆目录（先读这个）
.agents/memory/*.md       按域拆开的稳定事实
.agents/scratch/          会话草稿（git 忽略）
```

写入规则：一条记忆一个主题；更新覆盖原段；不追加聊天记录。
