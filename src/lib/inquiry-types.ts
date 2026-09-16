export type InquiryKind = 'contact' | 'forward' | 'route-test';

export type Inquiry = {
	id: string;
	userId: string | null;
	username: string | null;
	kind: InquiryKind;
	name: string;
	contact: string;
	country: string | null;
	address: string | null;
	productId: string | null;
	addonId: string | null;
	quantity: number;
	message: string | null;
	status: string;
	createdAt: string;
};
