<script lang="ts">
	import Navbar from '$lib/Navbar.svelte';
	import Footer from '$lib/Footer.svelte';
	import { auth } from '$lib/stores/auth';

	auth.initialize();

	let name = $state('');
	let contact = $state('');
	let message = $state('');
	let submitting = $state(false);
	let error = $state('');
	let success = $state('');

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
				body: JSON.stringify({ kind: 'contact', name, contact, message })
			});
			const payload = await res.json();
			if (!res.ok) {
				error = payload.error || '提交失败';
				return;
			}
			success = '已收到。库存、鉴定和合作问题我们会通过你留下的联系方式回复。';
			message = '';
		} catch {
			error = '网络异常，请稍后重试';
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>联系我们-高木桑收藏品店</title>
	<meta name="description" content="联系广州市白云区三元里高木桑收藏品店，咨询库存、代寄、邮路测试与合作。" />
</svelte:head>

<Navbar />

<div class="bg-gray-50 min-h-screen">
	<main class="mx-auto max-w-3xl px-6 py-16">
		<p class="text-sm font-semibold text-rose-500">广州市白云区三元里</p>
		<h1 class="mt-2 text-4xl font-bold text-gray-900">联系我们</h1>
		<p class="mt-4 text-gray-600">
			适合问库存、限量组是否还在、信件定制文案，或商业合作。需要实寄请走
			<a class="text-rose-500" href="/shop/sent">代寄 / 邮路测试</a>。
		</p>

		<form class="mt-10 space-y-4 rounded-3xl bg-white p-6 shadow-sm" onsubmit={submit}>
			<label class="block text-sm text-gray-600">称呼
				<input class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2" bind:value={name} required />
			</label>
			<label class="block text-sm text-gray-600">联系方式
				<input class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2" bind:value={contact} required />
			</label>
			<label class="block text-sm text-gray-600">想咨询什么
				<textarea class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2" rows="6" bind:value={message} required></textarea>
			</label>
			<button class="rounded-full bg-gradient-to-r from-rose-500 to-purple-600 px-6 py-3 font-semibold text-white disabled:opacity-60" disabled={submitting}>
				{submitting ? '提交中…' : '发送咨询'}
			</button>
			{#if error}<p class="text-sm text-rose-500">{error}</p>{/if}
			{#if success}<p class="text-sm text-emerald-600">{success}</p>{/if}
		</form>
	</main>
</div>

<Footer />
