// AES-256-GCM Encryption — WebCrypto API (no external deps)

const enc = new TextEncoder();
const dec = new TextDecoder();

// ponytail: inline helpers instead of a base64 util module — two callers, no shared module needed yet
const toBase64 = (buf: ArrayBuffer | ArrayBufferView): string => {
	const u8 = buf instanceof ArrayBuffer ? new Uint8Array(buf) : new Uint8Array(buf.buffer, buf.byteOffset, buf.byteLength);
	return btoa(Array.from(u8, (b) => String.fromCharCode(b)).join(''));
};
const fromBase64 = (s: string): ArrayBuffer => {
	const bin = atob(s);
	const u8 = new Uint8Array(bin.length);
	for (let i = 0; i < bin.length; i++) {
		u8[i] = bin.charCodeAt(i);
	}
	return u8.buffer;
};

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
