import { json } from '@sveltejs/kit';
import { getUserFromRequest } from '$lib/server/auth';
import { saveRagQuestion } from '$lib/server/inquiries';
import { askRag } from '$lib/server/rag';

export async function POST({ request }) {
	let question = '';
	try {
		const body = await request.json();
		question = String(body?.question ?? '').trim();
	} catch {
		return json({ error: '请求格式无效' }, { status: 400 });
	}

	if (!question) {
		return json({ error: '请输入问题' }, { status: 400 });
	}

	try {
		const result = await askRag(question);
		const user = getUserFromRequest(request);
		saveRagQuestion({
			userId: user?.id,
			question,
			answer: result.answer
		});
		return json(result);
	} catch (error) {
		console.error('RAG ask proxy error:', error);
		return json({ error: '店员助手暂时无法回答，请稍后重试。' }, { status: 502 });
	}
}
