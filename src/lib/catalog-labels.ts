export const PRODUCT_LABELS: Record<string, string> = {
	'jp-postcard-set': '日本风景明信片套装',
	'eu-vintage-stamps': '欧洲复古邮票散票',
	'letter-custom': '手写信件定制',
	'gift-box': '收藏礼盒套装',
	'mail-route-test': '全球邮路测试',
	'postcard-forward': '明信片代寄服务',
	'limited-stamps': '限量珍藏邮票',
	'vintage-postcard': '复古系列明信片'
};

export function productLabel(id: string | null | undefined): string {
	if (!id) return '';
	return PRODUCT_LABELS[id] ?? id;
}
