import { v4 as uuidv4 } from 'uuid';
import db from '$lib/server/db';
import type { Inquiry, InquiryKind } from '$lib/inquiry-types';

type InquiryRow = {
	id: string;
	user_id: string | null;
	username: string | null;
	kind: InquiryKind;
	name: string;
	contact: string;
	country: string | null;
	address: string | null;
	product_id: string | null;
	addon_id: string | null;
	quantity: number;
	message: string | null;
	status: string;
	created_at: string;
};

export type NewInquiry = {
	userId?: string | null;
	username?: string | null;
	kind: InquiryKind;
	name: string;
	contact: string;
	country?: string;
	address?: string;
	productId?: string;
	addonId?: string;
	quantity?: number;
	message?: string;
};

function mapInquiry(row: InquiryRow): Inquiry {
	return {
		id: row.id,
		userId: row.user_id,
		username: row.username,
		kind: row.kind,
		name: row.name,
		contact: row.contact,
		country: row.country,
		address: row.address,
		productId: row.product_id,
		addonId: row.addon_id,
		quantity: row.quantity,
		message: row.message,
		status: row.status,
		createdAt: row.created_at
	};
}

export function createInquiry(input: NewInquiry): Inquiry {
	const id = uuidv4();
	const quantity = Math.max(1, Math.min(99, Number(input.quantity) || 1));
	db.prepare(
		`INSERT INTO inquiries (
      id, user_id, username, kind, name, contact, country, address,
      product_id, addon_id, quantity, message, status
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'received')`
	).run(
		id,
		input.userId ?? null,
		input.username ?? null,
		input.kind,
		input.name.trim(),
		input.contact.trim(),
		input.country?.trim() || null,
		input.address?.trim() || null,
		input.productId || null,
		input.addonId || null,
		quantity,
		input.message?.trim() || null
	);
	return getInquiry(id)!;
}

export function getInquiry(id: string): Inquiry | undefined {
	const row = db.prepare('SELECT * FROM inquiries WHERE id = ?').get(id) as InquiryRow | undefined;
	return row ? mapInquiry(row) : undefined;
}

export function listInquiries(userId: string): Inquiry[] {
	const rows = db
		.prepare('SELECT * FROM inquiries WHERE user_id = ? ORDER BY created_at DESC')
		.all(userId) as InquiryRow[];
	return rows.map(mapInquiry);
}

export function saveRagQuestion(input: { userId?: string | null; question: string; answer?: string }) {
	db.prepare(
		'INSERT INTO rag_questions (id, user_id, question, answer) VALUES (?, ?, ?, ?)'
	).run(uuidv4(), input.userId ?? null, input.question, input.answer ?? null);
}

export function listRagQuestions(userId: string, limit = 8) {
	return db
		.prepare(
			'SELECT id, question, answer, created_at as createdAt FROM rag_questions WHERE user_id = ? ORDER BY created_at DESC LIMIT ?'
		)
		.all(userId, limit) as { id: string; question: string; answer: string | null; createdAt: string }[];
}
