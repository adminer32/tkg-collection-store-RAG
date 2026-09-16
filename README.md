# sv

Everything you need to build a Svelte project, powered by [`sv`](https://github.com/sveltejs/cli).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```sh
# create a new project
npx sv create my-app
```

To recreate this project with the same configuration:

```sh
# recreate this project
npx sv@0.17.0 create --template minimal --types ts --install npm ./
```

## RAG 店员助手（Go）

商品目录检索与问答跑在独立的 Go 服务里，不占用现有 SvelteKit/SQLite 认证链路。

```sh
# 终端 1：Go RAG
cd services/rag
go test ./...
go run .

# 终端 2：SvelteKit
npm run dev
```

打开 `/shop` 即可浏览目录并向店员助手提问。SvelteKit 通过 `RAG_URL`（默认 `http://127.0.0.1:8081`）反代到 Go。商店、代寄页和 RAG 共用 `services/rag/data/catalog.json`。

相关页面：

- `/shop` 目录 + 店员助手
- `/shop/sent` 代寄 / 邮路测试申请
- `/about/contact` 店铺咨询
- `/userspace` 我的申请与提问记录

也可用 Compose 同时启动：

```sh
docker compose up --build
```

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```sh
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```sh
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.
