<script lang="ts">
	import Navbar from '$lib/Navbar.svelte';
	import Footer from '$lib/Footer.svelte';
	import { auth } from '$lib/stores/auth';
	import type { RagProduct } from '$lib/rag-types';

	auth.initialize();

	let {
		data
	}: {
		data: {
			services: RagProduct[];
			addons: RagProduct[];
			selectedService: string;
			selectedAddon: string;
		};
	} = $props();

	let name = $state('');
	let contact = $state('');
	let country = $state('中国');
	let address = $state('');
	let productId = $state(data.selectedService);
	let addonId = $state(data.selectedAddon);
	let quantity = $state(1);
	let message = $state('');
	let submitting = $state(false);
	let error = $state('');
	let success = $state('');

	const selected = $derived(data.services.find((item) => item.id === productId));
	const addon = $derived(data.addons.find((item) => item.id === addonId));
	const estimate = $derived((selected?.priceCny ?? 0) * quantity + (addon?.priceCny ?? 0));

	async function submit(e: Event) {
		e.preventDefault();
		if (submitting) return;
		submitting = true;
		error = '';
		success = '';
		try {
			const headers: Record<string, string> = { 'content-type': 'application/json' };
			if ($auth.token) headers.authorization = `Bearer ${$auth.token}`;
			const res = await fetch('/api/inquiries', {
				method: 'POST',
				headers,
				body: JSON.stringify({
					kind: productId === 'mail-route-test' ? 'route-test' : 'forward',
					name,
					contact,
					country,
					address,
					productId,
					addonId,
					quantity,
					message
				})
			});
			const payload = await res.json();
			if (!res.ok) {
				error = payload.error || '提交失败';
				return;
			}
			success = payload.message + (payload.inquiry?.id ? ` 编号 ${payload.inquiry.id.slice(0, 8)}` : '');
			address = '';
			message = '';
		} catch {
			error = '网络异常，请稍后重试';
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>代寄 / 邮路测试-高木桑收藏品店</title>
	<meta name="description" content="填写收件地址，委托高木桑收藏品店通过中国邮政代寄明信片或进行邮路测试。" />
</svelte:head>

<Navbar />

<div class="bg-slate-50 min-h-screen">
	<section class="mx-auto max-w-4xl px-6 py-14">
		<p class="text-sm font-semibold tracking-wide text-rose-500">中国邮政渠道</p>
		<h1 class="mt-2 text-4xl font-bold text-slate-900">代寄与邮路测试</h1>
		<p class="mt-4 text-slate-600">
			本店不提供商业快递。填写地址后，我们按现行邮资贴票、销戳并交寄。
			{#if $auth.isAuthenticated}
				提交后可在 <a class="text-rose-500" href="/userspace">个人中心</a> 查看记录。
			{:else}
				<a class="text-rose-500" href="/auth?from=/shop/sent">登录</a> 后申请会记到你的账号。
			{/if}
		</p>

		<form class="mt-10 space-y-6 rounded-3xl border border-slate-100 bg-white p-6 shadow-sm" onsubmit={submit}>
			<div class="grid gap-4 md:grid-cols-2">
				{#each data.services as service}
					<label class="flex cursor-pointer flex-col rounded-2xl border p-4" class:border-rose-400={productId === service.id} class:border-slate-200={productId !== service.id}>
						<input class="sr-only" type="radio" name="service" value={service.id} bind:group={productId} />
						<span class="text-sm text-rose-500">{service.series}</span>
						<span class="mt-1 font-semibold text-slate-900">{service.title}</span>
						<span class="mt-2 text-sm text-slate-500">{service.summary}</span>
						<span class="mt-3 font-medium text-rose-500">¥{service.priceCny} 起</span>
					</label>
				{/each}
			</div>

			<div class="grid gap-4 md:grid-cols-2">
				<label class="text-sm text-slate-600">称呼
					<input class="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2" bind:value={name} required />
				</label>
				<label class="text-sm text-slate-600">联系方式（微信 / 邮箱 / 电话）
					<input class="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2" bind:value={contact} required />
				</label>
				<label class="text-sm text-slate-600">国家或地区
					<input class="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2" bind:value={country} />
				</label>
				<label class="text-sm text-slate-600">数量
					<input class="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2" type="number" min="1" max="99" bind:value={quantity} />
				</label>
			</div>

			<label class="block text-sm text-slate-600">完整收件地址
				<textarea class="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2" rows="4" bind:value={address} required placeholder="姓名、邮编、省市区或海外地址、电话"></textarea>
			</label>

			<label class="block text-sm text-slate-600">加购明信片 / 礼盒（可选）
				<select class="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2" bind:value={addonId}>
					<option value="">不加购，仅代寄或测试</option>
					{#each data.addons as item}
						<option value={item.id}>{item.title} · ¥{item.priceCny}</option>
					{/each}
				</select>
			</label>

			<label class="block text-sm text-slate-600">备注
				<textarea class="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2" rows="3" bind:value={message} placeholder="纪念戳偏好、是否需要手写信件等"></textarea>
			</label>

			<div class="flex flex-wrap items-center justify-between gap-3">
				<p class="text-slate-700">预估服务+加购 <span class="font-semibold text-rose-500">¥{estimate}</span>（不含可能调整的邮政资费）</p>
				<button class="rounded-full bg-gradient-to-r from-rose-500 to-purple-600 px-6 py-3 font-semibold text-white disabled:opacity-60" disabled={submitting}>
					{submitting ? '提交中…' : '提交代寄申请'}
				</button>
			</div>
			{#if error}<p class="text-sm text-rose-500">{error}</p>{/if}
			{#if success}<p class="text-sm text-emerald-600">{success}</p>{/if}
		</form>
	</section>
</div>

<Footer />
