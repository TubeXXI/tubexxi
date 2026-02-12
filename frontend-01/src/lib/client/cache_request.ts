class RequestCache {
	private cache = new Map<string, { data: any; timestamp: number }>();
	private ttl = 5 * 60 * 1000; // 5 minutes default

	set(key: string, data: any): void {
		this.cache.set(key, { data, timestamp: Date.now() });
	}

	get(key: string): any | undefined {
		const cached = this.cache.get(key);
		if (!cached) return undefined;

		if (Date.now() - cached.timestamp > this.ttl) {
			this.cache.delete(key);
			return undefined;
		}

		return cached.data;
	}

	clear(key?: string): void {
		if (key) {
			this.cache.delete(key);
		} else {
			this.cache.clear();
		}
	}
}

export const requestCache = new RequestCache();
