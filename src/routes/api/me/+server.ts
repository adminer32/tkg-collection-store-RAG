import { json } from '@sveltejs/kit';
import { getUserFromRequest } from '$lib/server/auth';
import { listInquiries, listRagQuestions } from '$lib/server/inquiries';

export async function GET({ request }) {
	const user = getUserFromRequest(request);
	if (!user) {
		return json({ error: '请先登录' }, { status: 401 });
	}
	return json({
		user,
		inquiries: listInquiries(user.id),
		questions: listRagQuestions(user.id)
	});
}
