<script lang="ts">
	import { onMount } from 'svelte';
	import Navbar from '$lib/Navbar.svelte';
	import Footer from '$lib/Footer.svelte';
	import { auth } from '$lib/stores/auth';
	import { productLabel } from '$lib/catalog-labels';
	import type { Inquiry } from '$lib/inquiry-types';

	let loading = $state(true);
	let error = $state('');
	let inquiries = $state<Inquiry[]>([]);
	let questions = $state<{ id: string; question: string; answer: string | null; createdAt: string }[]>([]);

	const kindLabel: Record<string, string> = {
		contact: '店铺咨询',
		forward: '明信片代寄',
		'route-test': '邮路测试'
	};

	onMount(() => {
		auth.initialize();
		const unsub = auth.subscribe(async (state) => {
			if (!state.isAuthenticated) {
				loading = false;
				return;
			}
			loading = true;
			try {
				const res = await fetch('/api/me', {
					headers: { authorization: `Bearer ${state.token}` }
				});
				const payload = await res.json();
				if (!res.ok) {
					error = payload.error || '无法读取个人中心';
					return;
				}
				inquiries = payload.inquiries ?? [];
				questions = payload.questions ?? [];
			} catch {
				error = '网络异常';
			} finally {
				loading = false;
			}
		});
		return unsub;
	});
</script>

<svelte:head>
	<title>个人中心-高木桑收藏品店</title>
</svelte:head>

<Navbar />

<div class="bg-slate-50 min-h-screen">
	<main class="mx-auto max-w-4xl px-6 py-14">
		<h1 class="text-4xl font-bold text-slate-900">个人中心</h1>
		<p class="mt-3 text-slate-600">查看你提交的代寄申请、咨询，以及问过店员助手的问题。</p>

		{#if !$auth.isAuthenticated}
			<div class="mt-10 rounded-3xl border border-dashed border-slate-200 bg-white p-10 text-center">
				<p class="text-slate-600">登录后才能把申请记到你的账号下。</p>
				<a class="mt-4 inline-block rounded-full bg-gradient-to-r from-rose-500 to-purple-600 px-6 py-3 font-semibold text-white" href="/auth?from=/userspace">去登录</a>
			</div>
		{:else if loading}
			<p class="mt-10 text-slate-500">读取中…</p>
		{:else if error}
			<p class="mt-10 text-rose-500">{error}</p>
		{:else}
			<section class="mt-10">
				<h2 class="text-xl font-semibold text-slate-900">我的申请</h2>
				{#if inquiries.length === 0}
					<p class="mt-3 text-sm text-slate-500">还没有记录。可以从 <a class="text-rose-500" href="/shop/sent">代寄页</a> 或 <a class="text-rose-500" href="/about/contact">联系我们</a> 提交。</p>
				{:else}
					<div class="mt-4 space-y-3">
						{#each inquiries as item}
							<article class="rounded-2xl border border-slate-100 bg-white p-4">
								<div class="flex flex-wrap items-center justify-between gap-2">
									<p class="font-medium text-slate-800">{kindLabel[item.kind] ?? item.kind}</p>
									<span class="text-xs text-slate-400">{item.status} · {item.createdAt}</span>
								</div>
								<p class="mt-2 text-sm text-slate-600">
									{#if item.productId}{productLabel(item.productId)}{/if}
									{#if item.addonId} · 加购 {productLabel(item.addonId)}{/if}
									{#if !item.productId && !item.addonId}{item.message || '店铺咨询'}{/if}
								</p>
								{#if item.address}
									<p class="mt-1 text-sm text-slate-500">{item.address}</p>
								{/if}
								<p class="mt-1 text-xs text-slate-400">联系 {item.contact} · 数量 {item.quantity}</p>
							</article>
						{/each}
					</div>
				{/if}
			</section>

			<section class="mt-12">
				<h2 class="text-xl font-semibold text-slate-900">问过店员的问题</h2>
				{#if questions.length === 0}
					<p class="mt-3 text-sm text-slate-500">还没有提问。去 <a class="text-rose-500" href="/shop#assistant">商店助手</a> 问库存或品相。</p>
				{:else}
					<div class="mt-4 space-y-3">
						{#each questions as q}
							<article class="rounded-2xl border border-slate-100 bg-white p-4">
								<p class="font-medium text-slate-800">{q.question}</p>
								<p class="mt-2 whitespace-pre-wrap text-sm text-slate-600">{q.answer}</p>
							</article>
						{/each}
					</div>
				{/if}
			</section>
		{/if}
	</main>
</div>

<Footer />
