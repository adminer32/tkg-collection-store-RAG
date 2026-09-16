import { env } from '$env/dynamic/private';
import { loadLocalCatalog } from '$lib/server/catalog';
import type { AskResponse, RagHit, RagProduct } from '$lib/rag-types';

export type { AskResponse, RagHit, RagProduct };

const DEFAULT_RAG_URL = 'http://127.0.0.1:8081';

function ragBase(): string {
	return (env.RAG_URL || DEFAULT_RAG_URL).replace(/\/$/, '');
}

async function ragFetch(path: string, init?: RequestInit): Promise<Response> {
	const url = `${ragBase()}${path}`;
	return fetch(url, {
		...init,
		headers: {
			accept: 'application/json',
			...(init?.headers ?? {})
		}
	});
}

export async function fetchCatalog(): Promise<{ items: RagProduct[]; ragReady: boolean }> {
	try {
		const res = await ragFetch('/api/v1/catalog');
		if (!res.ok) {
			throw new Error(`RAG catalog failed: ${res.status}`);
		}
		const data = (await res.json()) as { items?: RagProduct[] };
		const items = data.items?.length ? data.items : loadLocalCatalog();
		return { items, ragReady: true };
	} catch {
		return { items: loadLocalCatalog(), ragReady: false };
	}
}

export async function askRag(question: string, limit = 4): Promise<AskResponse> {
	const res = await ragFetch('/api/v1/ask', {
		method: 'POST',
		headers: { 'content-type': 'application/json' },
		body: JSON.stringify({ question, limit })
	});
	if (!res.ok) {
		const text = await res.text();
		throw new Error(text || `RAG ask failed: ${res.status}`);
	}
	return (await res.json()) as AskResponse;
}

export async function searchRag(query: string, limit = 6): Promise<RagHit[]> {
	const params = new URLSearchParams({ q: query, limit: String(limit) });
	const res = await ragFetch(`/api/v1/search?${params.toString()}`);
	if (!res.ok) {
		throw new Error(`RAG search failed: ${res.status}`);
	}
	const data = (await res.json()) as { hits?: RagHit[] };
	return data.hits ?? [];
}
