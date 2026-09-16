import type { RagHit } from '$lib/rag-types';

export function hitAction(hit: RagHit): { href: string; label: string } | null {
	const source = hit.sourceId || hit.product?.id || '';
	if (hit.kind === 'product' || hit.product) {
		if (source === 'postcard-forward' || source === 'mail-route-test') {
			return { href: `/shop/sent?service=${encodeURIComponent(source)}`, label: '去填代寄/邮路' };
		}
		return { href: `/shop/sent?addon=${encodeURIComponent(source)}`, label: '加购并代寄' };
	}
	if (hit.kind === 'faq' || hit.kind === 'knowledge') {
		if ((hit.title + hit.summary).includes('代寄') || (hit.title + hit.summary).includes('邮路')) {
			return { href: '/shop/sent', label: '去代寄页' };
		}
		return { href: '/about/contact', label: '联系店铺' };
	}
	return null;
}
