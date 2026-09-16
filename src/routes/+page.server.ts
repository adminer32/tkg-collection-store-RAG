import { fetchCatalog } from '$lib/server/rag';

export async function load() {
	const { items, ragReady } = await fetchCatalog();
	const featured = items.slice(0, 6);
	return { featured, ragReady };
}
