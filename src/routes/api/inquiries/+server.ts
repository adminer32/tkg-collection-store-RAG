import { json } from '@sveltejs/kit';
import { getUserFromRequest } from '$lib/server/auth';
import { isAddonId, isServiceId } from '$lib/server/catalog';
import { createInquiry, listInquiries } from '$lib/server/inquiries';
import type { InquiryKind } from '$lib/inquiry-types';

const KINDS = new Set<InquiryKind>(['contact', 'forward', 'route-test']);

export async function GET({ request }) {
	const user = getUserFromRequest(request);
	if (!user) {
		return json({ error: '请先登录' }, { status: 401 });
	}
	return json({ items: listInquiries(user.id) });
}

export async function POST({ request }) {
	let body: Record<string, unknown>;
	try {
		body = await request.json();
	} catch {
		return json({ error: '请求格式无效' }, { status: 400 });
	}

	const kind = String(body.kind ?? '') as InquiryKind;
	const name = String(body.name ?? '').trim();
	const contact = String(body.contact ?? '').trim();
	if (!KINDS.has(kind)) {
		return json({ error: '未知的咨询类型' }, { status: 400 });
	}
	if (!name || !contact) {
		return json({ error: '请填写称呼和联系方式' }, { status: 400 });
	}

	const productId = String(body.productId ?? '').trim();
	const addonId = String(body.addonId ?? '').trim();
	const address = String(body.address ?? '').trim();

	if (kind === 'forward' || kind === 'route-test') {
		if (!address) {
			return json({ error: '请填写完整收件地址' }, { status: 400 });
		}
		if (!isServiceId(productId)) {
			return json({ error: '请选择代寄或邮路测试服务' }, { status: 400 });
		}
		if (addonId && !isAddonId(addonId)) {
			return json({ error: '加购商品不存在' }, { status: 400 });
		}
	}

	const user = getUserFromRequest(request);
	const inquiry = createInquiry({
		userId: user?.id,
		username: user?.username,
		kind,
		name,
		contact,
		country: String(body.country ?? ''),
		address,
		productId,
		addonId,
		quantity: Number(body.quantity ?? 1),
		message: String(body.message ?? '')
	});

	return json({ message: '已收到，我们会按中国邮政渠道跟进。', inquiry }, { status: 201 });
}
