import jwt from 'jsonwebtoken';
import { env } from '$env/dynamic/private';

export const JWT_SECRET = env.JWT_SECRET || 'TAKAGISANWAKAWAIIDESU520';

export type AuthUser = {
	id: string;
	username: string;
};

export function getUserFromRequest(request: Request): AuthUser | null {
	const header = request.headers.get('authorization');
	if (!header?.startsWith('Bearer ')) return null;
	try {
		const payload = jwt.verify(header.slice(7), JWT_SECRET) as {
			userId: string;
			username: string;
		};
		if (!payload?.userId || !payload?.username) return null;
		return { id: payload.userId, username: payload.username };
	} catch {
		return null;
	}
}
