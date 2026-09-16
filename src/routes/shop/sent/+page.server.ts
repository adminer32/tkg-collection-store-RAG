import { isServiceId, loadLocalCatalog } from '$lib/server/catalog';

export async function load({ url }) {
	const catalog = loadLocalCatalog();
	const services = catalog.filter((item) => isServiceId(item.id));
	const addons = catalog.filter((item) => !isServiceId(item.id));
	const service = url.searchParams.get('service') || 'postcard-forward';
	const addon = url.searchParams.get('addon') || '';
	return {
		services,
		addons,
		selectedService: isServiceId(service) ? service : 'postcard-forward',
		selectedAddon: addons.some((item) => item.id === addon) ? addon : ''
	};
}
