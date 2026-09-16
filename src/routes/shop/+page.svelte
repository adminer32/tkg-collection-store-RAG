<script lang="ts">
	import Navbar from '$lib/Navbar.svelte';
	import Footer from '$lib/Footer.svelte';
	import { hitAction } from '$lib/hit-links';
	import { auth } from '$lib/stores/auth';
	import type { AskResponse, RagProduct } from '$lib/rag-types';

	auth.initialize();

	let { data }: { data: { items: RagProduct[]; ragReady: boolean } } = $props();

	let question = $state('日本明信片怎么选？');
	let asking = $state(false);
	let askError = $state('');
	let result = $state<AskResponse | null>(null);

	const suggestions = ['代寄多少钱', '有没有欧洲邮票', '怎么看齿孔品相', '邮路测试怎么用'];

	async function ask(e?: Event) {
		e?.preventDefault();
		if (asking) return;
		const q = question.trim();
		if (!q) {
			askError = '请先输入你想问的商品或服务';
			return;
		}
		asking = true;
		askError = '';
		try {
			const headers: Record<string, string> = { 'content-type': 'application/json' };
			if ($auth.token) headers.authorization = `Bearer ${$auth.token}`;
			const res = await fetch('/api/rag/ask', {
				method: 'POST',
				headers,
				body: JSON.stringify({ question: q })
			});
			const payload = await res.json();
			if (!res.ok) {
				askError = payload.error || '提问失败';
				result = null;
				return;
			}
			result = payload as AskResponse;
		} catch {
			askError = '网络异常，店员助手暂时离线';
			result = null;
		} finally {
			asking = false;
		}
	}
</script>

<svelte:head>
	<title>精选商店-高木桑收藏品店</title>
	<meta name="description" content="浏览明信片、邮票与代寄服务，并向店员助手提问库存、价格与邮路问题。" />
</svelte:head>

<Navbar />

<div class="bg-slate-50 min-h-screen">
	<section id="assistant" class="relative overflow-hidden border-b border-rose-100 bg-white">
		<div class="mx-auto max-w-6xl px-6 py-14">
			<p class="text-sm font-semibold tracking-wide text-rose-500">精选商店 · 店员助手</p>
			<h1 class="mt-2 text-4xl font-bold tracking-tight text-slate-900">问一问店里还有什么</h1>
			<p class="mt-4 max-w-2xl text-slate-600">
				商品目录、FAQ 和邮票知识库由 Go RAG 按段落切块后检索。你可以问价格、库存、齿孔品相或代寄；回答会附上切块来源。
			</p>

			{#if !data.ragReady}
				<div class="mt-6 rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
					店员助手暂未连上 Go 后端。本地可先运行
					<code class="mx-1 rounded bg-white px-1.5 py-0.5">go run .</code>
					于 <code class="rounded bg-white px-1.5 py-0.5">services/rag</code>。
				</div>
			{/if}

			<form class="mt-8 flex flex-col gap-3 sm:flex-row" onsubmit={ask}>
				<input
					bind:value={question}
					class="flex-1 rounded-full border border-slate-200 bg-white px-5 py-3 text-slate-800 shadow-sm outline-none ring-rose-200 focus:ring-2"
					placeholder="例如：日本明信片、代寄怎么收费、还有礼盒吗"
				/>
				<button
					type="submit"
					disabled={asking}
					class="rounded-full bg-gradient-to-r from-rose-500 to-purple-600 px-6 py-3 font-semibold text-white shadow-sm disabled:opacity-60"
				>
					{asking ? '检索中…' : '询问店员'}
				</button>
			</form>
			<div class="mt-4 flex flex-wrap gap-2">
				{#each suggestions as s}
					<button
						type="button"
						class="rounded-full border border-slate-200 bg-white px-3 py-1 text-sm text-slate-600 hover:border-rose-300 hover:text-rose-500"
						onclick={() => {
							question = s;
							ask();
						}}
					>
						{s}
					</button>
				{/each}
			</div>
			{#if askError}
				<p class="mt-4 text-sm text-rose-500">{askError}</p>
			{/if}
		</div>
	</section>

	{#if result}
		<section class="mx-auto max-w-6xl px-6 py-10">
			<div class="rounded-3xl border border-slate-100 bg-white p-6 shadow-sm">
				<p class="text-xs uppercase tracking-wider text-slate-400">RAG 回答 · {result.mode}</p>
				<h2 class="mt-2 text-lg font-semibold text-slate-900">「{result.question}」</h2>
				<pre class="mt-4 whitespace-pre-wrap font-sans text-sm leading-7 text-slate-700">{result.answer}</pre>
				{#if result.citations?.length}
					<div class="mt-6 grid gap-3 md:grid-cols-2">
						{#each result.citations as hit}
							{@const action = hitAction(hit)}
							<article class="rounded-2xl border border-slate-100 bg-slate-50 p-4">
								<p class="text-xs text-rose-500">
									{hit.kind === 'product' ? '商品' : hit.kind === 'knowledge' ? '知识切块' : 'FAQ'}
									· score {hit.score}
								</p>
								<h3 class="mt-1 font-semibold text-slate-800">{hit.title}</h3>
								{#if hit.section && hit.section !== hit.title}
									<p class="mt-0.5 text-xs text-slate-400">{hit.section}</p>
								{/if}
								<p class="mt-2 text-sm text-slate-600">{hit.summary}</p>
								{#if action}
									<a class="mt-3 inline-block text-sm font-medium text-purple-600 hover:text-rose-500" href={action.href}>
										{action.label}
									</a>
								{/if}
							</article>
						{/each}
					</div>
				{/if}
			</div>
		</section>
	{/if}

	<section class="mx-auto max-w-6xl px-6 pb-16">
		<div class="mb-6 flex items-end justify-between gap-4">
			<div>
				<h2 class="text-2xl font-bold text-slate-900">在售目录</h2>
				<p class="mt-1 text-sm text-slate-500">目录与 RAG 共用 `services/rag/data/catalog.json`；助手离线时商店仍可读这份文件。</p>
			</div>
		</div>

		{#if data.items.length === 0}
			<div class="rounded-3xl border border-dashed border-slate-200 bg-white p-10 text-center text-slate-500">
				暂时没有读到商品目录。确认 `services/rag/data/catalog.json` 存在后刷新本页。
			</div>
		{:else}
			<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each data.items as item}
					<article class="flex flex-col rounded-3xl border border-slate-100 bg-white p-5 shadow-sm">
						<div class="flex items-center justify-between gap-3">
							<span class="rounded-full bg-rose-50 px-3 py-1 text-xs font-medium text-rose-500">{item.series}</span>
							<span class="text-sm text-slate-400">库存 {item.stock}</span>
						</div>
						<h3 class="mt-3 text-lg font-bold text-slate-900">{item.title}</h3>
						<p class="mt-1 text-sm text-slate-500">{item.category}</p>
						<p class="mt-3 flex-1 text-sm leading-6 text-slate-600">{item.summary}</p>
						<div class="mt-5 flex items-center justify-between gap-3">
							<p class="text-xl font-semibold text-rose-500">¥{item.priceCny}</p>
							<div class="flex items-center gap-3">
								{#if item.id === 'postcard-forward' || item.id === 'mail-route-test'}
									<a class="text-sm font-medium text-rose-500" href={`/shop/sent?service=${item.id}`}>去办理</a>
								{:else}
									<a class="text-sm font-medium text-rose-500" href={`/shop/sent?addon=${item.id}`}>加购代寄</a>
								{/if}
								<button
									type="button"
									class="text-sm font-medium text-purple-600 hover:text-rose-500"
									onclick={() => {
										question = `${item.title}还剩多少？怎么买？`;
										ask();
									}}
								>
									问这件
								</button>
							</div>
						</div>
					</article>
				{/each}
			</div>
		{/if}
	</section>
</div>

<Footer />
