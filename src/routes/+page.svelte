<script lang="ts">
	import Navbar from '$lib/Navbar.svelte';
	import Footer from '$lib/Footer.svelte';
	import type { RagProduct } from '$lib/rag-types';

	let { data }: { data: { featured: RagProduct[]; ragReady: boolean } } = $props();

	const services = [
		{
			href: '/shop',
			title: '精选商店',
			text: '明信片、邮票、礼盒与信件定制，目录与店员助手共用同一份商品数据。'
		},
		{
			href: '/shop/sent?service=postcard-forward',
			title: '明信片代寄',
			text: '贴票、销戳、交寄，单张 18 元起。实寄一律走中国邮政。'
		},
		{
			href: '/shop/sent?service=mail-route-test',
			title: '邮路测试',
			text: '向指定地址寄出测试片，确认信箱能否收到。'
		},
		{
			href: '/shop#assistant',
			title: '店员助手',
			text: '问价格、库存、齿孔品相或代寄流程，回答带来源切块。'
		}
	];
</script>

<svelte:head>
	<title>高木桑收藏品店</title>
	<meta
		name="description"
		content="广州市白云区三元里高木桑收藏品店：明信片、邮票与中国邮政代寄。浏览商店、办理邮路测试，或向店员助手提问。"
	/>
</svelte:head>

<Navbar />

<div class="bg-slate-50 min-h-screen">
	<section class="relative overflow-hidden border-b border-rose-100 bg-white">
		<div class="mx-auto grid max-w-6xl items-center gap-10 px-6 py-16 md:grid-cols-2">
			<div>
				<p class="text-sm font-semibold tracking-wide text-rose-500">广州市白云区三元里</p>
				<h1 class="mt-3 text-4xl font-bold tracking-tight text-slate-900 sm:text-5xl">
					高木桑收藏品店
				</h1>
				<p class="mt-5 max-w-xl text-lg leading-8 text-slate-600">
					明信片、邮票与手写信件。收藏留下纸面温度，实寄交给中国邮政。这里是店铺首页，不是关于我们。
				</p>
				<div class="mt-8 flex flex-wrap gap-3">
					<a
						class="rounded-full bg-gradient-to-r from-rose-500 to-purple-600 px-6 py-3 font-semibold text-white shadow-sm"
						href="/shop"
					>
						进入商店
					</a>
					<a
						class="rounded-full border border-slate-200 bg-white px-6 py-3 font-semibold text-slate-700"
						href="/shop/sent"
					>
						代寄 / 邮路测试
					</a>
					<a class="rounded-full px-6 py-3 font-semibold text-purple-600" href="/about">了解店铺</a>
				</div>
			</div>
			<div class="rounded-3xl border border-slate-100 bg-slate-50 p-6 shadow-sm">
				<p class="text-xs uppercase tracking-wider text-slate-400">今日可问</p>
				<ul class="mt-4 space-y-3 text-sm text-slate-600">
					<li>日本明信片怎么选？</li>
					<li>代寄多少钱？</li>
					<li>欧洲散票是什么年代？</li>
					<li>怎么看齿孔品相？</li>
				</ul>
				<a class="mt-6 inline-block text-sm font-medium text-rose-500" href="/shop#assistant">去问店员 →</a>
			</div>
		</div>
	</section>

	<section class="mx-auto max-w-6xl px-6 py-14">
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
			{#each services as item}
				<a class="rounded-3xl border border-slate-100 bg-white p-5 shadow-sm hover:border-rose-200" href={item.href}>
					<h2 class="font-semibold text-slate-900">{item.title}</h2>
					<p class="mt-2 text-sm leading-6 text-slate-600">{item.text}</p>
				</a>
			{/each}
		</div>
	</section>

	<section class="mx-auto max-w-6xl px-6 pb-16">
		<div class="mb-6 flex items-end justify-between gap-4">
			<div>
				<h2 class="text-2xl font-bold text-slate-900">在售精选</h2>
				<p class="mt-1 text-sm text-slate-500">与商店、RAG 共用同一份目录。</p>
			</div>
			<a class="text-sm font-medium text-rose-500" href="/shop">查看全部</a>
		</div>
		{#if data.featured.length === 0}
			<div class="rounded-3xl border border-dashed border-slate-200 bg-white p-10 text-center text-slate-500">
				暂时没有读到商品目录。
			</div>
		{:else}
			<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each data.featured as item}
					<article class="flex flex-col rounded-3xl border border-slate-100 bg-white p-5 shadow-sm">
						<div class="flex items-center justify-between gap-3">
							<span class="rounded-full bg-rose-50 px-3 py-1 text-xs font-medium text-rose-500">{item.series}</span>
							<span class="text-sm text-slate-400">库存 {item.stock}</span>
						</div>
						<h3 class="mt-3 text-lg font-bold text-slate-900">{item.title}</h3>
						<p class="mt-3 flex-1 text-sm leading-6 text-slate-600">{item.summary}</p>
						<div class="mt-5 flex items-center justify-between">
							<p class="text-xl font-semibold text-rose-500">¥{item.priceCny}</p>
							<a class="text-sm font-medium text-purple-600" href="/shop">去商店</a>
						</div>
					</article>
				{/each}
			</div>
		{/if}
	</section>
</div>

<Footer />
