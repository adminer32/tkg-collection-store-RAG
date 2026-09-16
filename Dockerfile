# --- 第一阶段：构建 (Builder Stage) ---
FROM node:20-slim AS builder
WORKDIR /app

# 在构建阶段安装编译工具（g++、make、python3）
# 这些工具是为了让 npm 能够编译像 better-sqlite3 这样的原生 C++ 模块
RUN apt-get update && apt-get install -y \
    python3 \
    make \
    g++ \
    && rm -rf /var/lib/apt/lists/*

COPY package*.json ./

# 安装所有依赖（包括 devDependencies）
RUN npm install

COPY . .

# 执行 SvelteKit 构建
RUN npm run build

# --- 第二阶段：生产 (Production Stage) ---
FROM node:20-slim
WORKDIR /app

# 生产环境只需安装轻量级的 sqlite3 运行时库
# 此时不再需要 python/make/g++，镜像体积大幅缩小
RUN apt-get update && apt-get install -y \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# 设置环境变量
ENV HOST=0.0.0.0
ENV PORT=3000
ENV NODE_ENV=production
# 预设数据库路径环境变量（供代码调用）
ENV DATABASE_PATH="/app/data/database.sqlite"
# Go RAG sidecar（docker compose 中指向 rag 服务）
ENV RAG_URL="http://127.0.0.1:8081"

# 1. 关键：创建持久化数据目录并赋予权限
# 这样即便你还没挂载 Volume，程序启动时也不会因为找不到目录而报错
RUN mkdir -p /app/data && chmod 777 /app/data

# 2. 从 builder 阶段拷贝构建产物
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/services/rag/data/catalog.json ./services/rag/data/catalog.json

# 3. 仅安装生产环境依赖
# 注意：如果你的 SQLite 驱动在这一步仍然触发编译，说明上面的拷贝需要调整
# 但对于大多数 SvelteKit 项目，直接拷贝 node_modules 或在这里安装生产版即可
RUN npm install --production

EXPOSE 3000

# 使用 JSON 数组格式启动，对信号处理更友好（方便 docker stop 快速响应）
CMD ["node", "build"]