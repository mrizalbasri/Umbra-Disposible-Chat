// ECDH P-256 Key Exchange — WebCrypto API (no external deps)

const ECDH_PARAMS = { name: 'ECDH', namedCurve: 'P-256' } as const;

/** Generate ECDH keypair for this session. Call once on room join. */
export function generateKeyPair(): Promise<CryptoKeyPair> {
	return crypto.subtle.generateKey(ECDH_PARAMS, true, ['deriveKey', 'deriveBits']);
}

/** Export public key → base64 for sending to server/peer. */
export async function exportPublicKey(publicKey: CryptoKey): Promise<string> {
	const raw = await crypto.subtle.exportKey('raw', publicKey);
	return btoa(Array.from(new Uint8Array(raw), (b) => String.fromCharCode(b)).join(''));
}

/** Import a peer's base64 public key → CryptoKey. */
export function importPublicKey(base64Key: string): Promise<CryptoKey> {
	const bin = atob(base64Key);
	const raw = new Uint8Array(bin.length);
	for (let i = 0; i < bin.length; i++) {
		raw[i] = bin.charCodeAt(i);
	}
	return crypto.subtle.importKey('raw', raw.buffer, ECDH_PARAMS, true, []);
}

/** Derive 256-bit shared secret from our private key + peer's public key. */
export function deriveSharedSecret(
	privateKey: CryptoKey,
	peerPublicKey: CryptoKey
): Promise<ArrayBuffer> {
	return crypto.subtle.deriveBits({ name: 'ECDH', public: peerPublicKey }, privateKey, 256);
}
