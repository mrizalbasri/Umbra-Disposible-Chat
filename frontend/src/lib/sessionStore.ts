import { writable } from 'svelte/store';

export interface ActiveSession {
	roomId: string;
	myEcdhKeyPair: CryptoKeyPair | null;
	myEcdsaKeyPair: CryptoKeyPair | null;
	myEcdsaPkBase64: string;
	peerEcdhPublicKey: CryptoKey | null;
	aesKey: CryptoKey | null;
	myNickname: string;
	partnerNickname: string;
	memberId: string;
	socket: WebSocket | null;
}

export const activeSession = writable<ActiveSession | null>(null);
