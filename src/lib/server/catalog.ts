import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import type { RagProduct } from '$lib/rag-types';

export function loadLocalCatalog(): RagProduct[] {
	const path = join(process.cwd(), 'services/rag/data/catalog.json');
	return JSON.parse(readFileSync(path, 'utf8')) as RagProduct[];
}

export function getCatalogItem(id: string | null | undefined): RagProduct | undefined {
	if (!id) return undefined;
	return loadLocalCatalog().find((item) => item.id === id);
}

export const SERVICE_IDS = ['postcard-forward', 'mail-route-test'] as const;

export function isServiceId(id: string): boolean {
	return (SERVICE_IDS as readonly string[]).includes(id);
}

export function isAddonId(id: string): boolean {
	return Boolean(id) && !isServiceId(id) && Boolean(getCatalogItem(id));
}
