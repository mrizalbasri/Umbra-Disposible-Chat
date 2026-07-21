import { env } from '$env/dynamic/public';

export function getBackendHost(): string {
	if (env.PUBLIC_BACKEND_URL && env.PUBLIC_BACKEND_URL.trim() !== '') {
		return env.PUBLIC_BACKEND_URL.replace(/\/$/, '');
	}
	if (env.PUBLIC_API_URL && env.PUBLIC_API_URL.trim() !== '') {
		return env.PUBLIC_API_URL.replace(/\/$/, '');
	}
	const legacyViteApiUrl = (import.meta.env as Record<string, string | undefined>).VITE_API_URL;
	if (legacyViteApiUrl && legacyViteApiUrl.trim() !== '') {
		return legacyViteApiUrl.replace(/\/$/, '');
	}
	if (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')) {
		return 'http://localhost:8080';
	}
	// Default to live Azure backend
	return 'https://backend-umbra.ashydune-5f4fef09.southeastasia.azurecontainerapps.io';
}

export function getApiBaseUrl(): string {
	return `${getBackendHost()}/v1/api`;
}

export function getWsBaseUrl(): string {
	const host = getBackendHost();
	const wsProtocol = host.startsWith('https') ? 'wss' : 'ws';
	const cleanHost = host.replace(/^https?:\/\//, '');
	return `${wsProtocol}://${cleanHost}/ws`;
}
