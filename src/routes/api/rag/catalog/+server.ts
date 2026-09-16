import { json } from '@sveltejs/kit';
import { fetchCatalog } from '$lib/server/rag';

export async function GET() {
	try {
		const { items, ragReady } = await fetchCatalog();
		return json({ items, ragReady });
	} catch (error) {
		console.error('RAG catalog proxy error:', error);
		return json({ error: '商店目录暂时不可用，请确认 Go RAG 服务已启动。' }, { status: 502 });
	}
}
