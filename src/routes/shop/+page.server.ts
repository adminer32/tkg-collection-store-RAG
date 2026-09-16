import { fetchCatalog } from '$lib/server/rag';

export async function load() {
	return fetchCatalog();
}
