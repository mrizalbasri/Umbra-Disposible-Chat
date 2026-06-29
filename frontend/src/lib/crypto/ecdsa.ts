// ECDSA P-256 Digital Signature — WebCrypto API (no external deps)

const ECDSA_PARAMS = { name: 'ECDSA', namedCurve: 'P-256' } as const;
const SIGN_ALGO = { name: 'ECDSA', hash: 'SHA-256' } as const;

const toBase64 = (buf: ArrayBuffer): string =>
	btoa(Array.from(new Uint8Array(buf), (b) => String.fromCharCode(b)).join(''));
const fromBase64 = (s: string): Uint8Array =>
	Uint8Array.from(atob(s), (c) => c.charCodeAt(0));

/** Generate ECDSA keypair for signing. Separate from ECDH keypair. */
export function generateSigningKeyPair(): Promise<CryptoKeyPair> {
	return crypto.subtle.generateKey(ECDSA_PARAMS, true, ['sign', 'verify']);
}

/** Export signing public key → base64 to share with peer. */
export async function exportSigningPublicKey(publicKey: CryptoKey): Promise<string> {
	const raw = await crypto.subtle.exportKey('raw', publicKey);
	return toBase64(raw);
}

/** Import peer's base64 signing public key → CryptoKey for verify(). */
export function importSigningPublicKey(base64Key: string): Promise<CryptoKey> {
	return crypto.subtle.importKey('raw', fromBase64(base64Key), ECDSA_PARAMS, true, ['verify']);
}

/** Sign data with our private key → base64 signature. */
export async function sign(data: ArrayBuffer, privateKey: CryptoKey): Promise<string> {
	const sigBuf = await crypto.subtle.sign(SIGN_ALGO, privateKey, data);
	return toBase64(sigBuf);
}

/**
 * Verify a signature against data + peer's public key.
 * Returns false on any failure (tampered, wrong key, malformed) — never throws.
 */
export async function verify(
	data: ArrayBuffer,
	signature: string,
	publicKey: CryptoKey
): Promise<boolean> {
	try {
		return await crypto.subtle.verify(SIGN_ALGO, publicKey, fromBase64(signature), data);
	} catch {
		return false; // malformed signature — treat as invalid, not crash
	}
}
