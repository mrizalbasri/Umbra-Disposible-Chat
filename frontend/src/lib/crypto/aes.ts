// AES-256-GCM Encryption — WebCrypto API (no external deps)

const enc = new TextEncoder();
const dec = new TextDecoder();

// ponytail: inline helpers instead of a base64 util module — two callers, no shared module needed yet
const toBase64 = (buf: ArrayBuffer): string =>
	btoa(Array.from(new Uint8Array(buf), (b) => String.fromCharCode(b)).join(''));
const fromBase64 = (s: string): Uint8Array =>
	Uint8Array.from(atob(s), (c) => c.charCodeAt(0));

/** Wrap a shared secret (ArrayBuffer) into an AES-256-GCM CryptoKey. */
export function deriveAESKey(sharedSecret: ArrayBuffer): Promise<CryptoKey> {
	return crypto.subtle.importKey('raw', sharedSecret, { name: 'AES-GCM' }, false, [
		'encrypt',
		'decrypt'
	]);
}

/** Encrypt plaintext → { ciphertext, iv } (both base64). IV is random per call. */
export async function encrypt(
	plaintext: string,
	aesKey: CryptoKey
): Promise<{ ciphertext: string; iv: string }> {
	const iv = crypto.getRandomValues(new Uint8Array(12)); // 96-bit IV — GCM standard
	const cipherBuf = await crypto.subtle.encrypt(
		{ name: 'AES-GCM', iv },
		aesKey,
		enc.encode(plaintext)
	);
	return { ciphertext: toBase64(cipherBuf), iv: toBase64(iv) };
}

/** Decrypt ciphertext (base64) + iv (base64) → plaintext string. */
export async function decrypt(
	ciphertext: string,
	iv: string,
	aesKey: CryptoKey
): Promise<string> {
	const plainBuf = await crypto.subtle.decrypt(
		{ name: 'AES-GCM', iv: fromBase64(iv) },
		aesKey,
		fromBase64(ciphertext)
	);
	return dec.decode(plainBuf);
}
