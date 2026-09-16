export type RagProduct = {
	id: string;
	title: string;
	category: string;
	series: string;
	priceCny: number;
	stock: number;
	tags: string[];
	summary: string;
	description: string;
};

export type RagHit = {
	id: string;
	kind: 'product' | 'faq' | 'knowledge' | string;
	title: string;
	section?: string;
	sourceId?: string;
	category?: string;
	series?: string;
	priceCny?: number;
	stock?: number;
	tags: string[];
	summary: string;
	score: number;
	highlights?: string[];
	product?: RagProduct;
};

export type AskResponse = {
	question: string;
	answer: string;
	citations: RagHit[];
	mode: string;
};
