<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { activeSession } from '$lib/sessionStore';
	import {
		generateKeyPair,
		exportPublicKey,
		importPublicKey,
		deriveSharedSecret,
		deriveAESKey,
		exportSigningPublicKey
	} from '$lib/crypto';

	import { env } from '$env/dynamic/public';

	// API and WS Configuration (dynamic with fallback to localhost)
	const BACKEND_HOST = env.PUBLIC_BACKEND_URL || 'http://localhost:8080';
	const BASE_URL = `${BACKEND_HOST}/v1/api`;
	const WS_BASE_URL = `${BACKEND_HOST.replace(/^http/, 'ws')}/ws`;

	let myNickname = 'Guest_User';
	let myEcdhKeyPair: CryptoKeyPair | null = null;
	let myEcdsaKeyPair: CryptoKeyPair | null = null;
	let myEcdsaPkBase64 = '';
	let peerEcdhPublicKey: CryptoKey | null = null;
	let aesKey: CryptoKey | null = null;
	let memberId = '';
	let roomId = '';
	let matchQueueId = '';
	let socket: WebSocket | null = null;

	function createClientMemberID(): string {
		if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
			return crypto.randomUUID();
		}
		return `${Date.now()}-${Math.random().toString(36).slice(2)}`;
	}

	async function generateKeys() {
		myEcdhKeyPair = await generateKeyPair();

		// Convert ECDH private key to ECDSA private key for signing
		const jwkPrivate = await crypto.subtle.exportKey('jwk', myEcdhKeyPair.privateKey);
		jwkPrivate.key_ops = ['sign'];
		const ecdsaPrivateKey = await crypto.subtle.importKey(
			'jwk',
			jwkPrivate,
			{ name: 'ECDSA', namedCurve: 'P-256' },
			true,
			['sign']
		);

		// Convert ECDH public key to ECDSA public key for verification
		const jwkPublic = await crypto.subtle.exportKey('jwk', myEcdhKeyPair.publicKey);
		jwkPublic.key_ops = ['verify'];
		const ecdsaPublicKey = await crypto.subtle.importKey(
			'jwk',
			jwkPublic,
			{ name: 'ECDSA', namedCurve: 'P-256' },
			true,
			['verify']
		);

		myEcdsaKeyPair = {
			privateKey: ecdsaPrivateKey,
			publicKey: ecdsaPublicKey
		};
		myEcdsaPkBase64 = await exportSigningPublicKey(myEcdsaKeyPair.publicKey);
	}

	function connectWebSocket(myEcdhPkBase64: string) {
		const wsUrl = `${WS_BASE_URL}?roomId=${roomId}&publicKey=${encodeURIComponent(myEcdhPkBase64)}`;
		socket = new WebSocket(wsUrl);

		socket.onmessage = async (event) => {
			try {
				const data = JSON.parse(event.data);
				if (data.roomId && data.roomId !== roomId) return;

				if (data.event === 'peer_joined') {
					peerEcdhPublicKey = await importPublicKey(data.payload.publicKey);
					const sharedSecret = await deriveSharedSecret(
						myEcdhKeyPair!.privateKey,
						peerEcdhPublicKey
					);
					aesKey = await deriveAESKey(sharedSecret);

					// Save active session to shared store and route to chatroom
					activeSession.set({
						roomId,
						myEcdhKeyPair,
						myEcdsaKeyPair,
						myEcdsaPkBase64,
						peerEcdhPublicKey,
						aesKey,
						myNickname,
						partnerNickname: '',
						memberId,
						socket: socket
					});

					goto(resolve('/chatroom'));
				}
			} catch (err) {
				console.error('Error handling WebSocket message in match:', err);
			}
		};
	}

	async function handleStartMatchmaking() {
		myNickname = `Guest_${Math.floor(1000 + Math.random() * 9000)}`;

		try {
			await generateKeys();
			const myEcdhPkBase64 = await exportPublicKey(myEcdhKeyPair!.publicKey);
			memberId = createClientMemberID();

			const res = await fetch(`${BASE_URL}/match/queue`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ publicKey: myEcdhPkBase64, nickname: myNickname })
			});
			const result = await res.json();

			if (result.responseCode === '17' && result.data?.status === 'waiting') {
				matchQueueId = result.data.queueId;
				roomId = result.data.queueId;
				connectWebSocket(myEcdhPkBase64);
				return;
			}

			if (result.responseCode === '00' && result.data?.status === 'matched') {
				matchQueueId = '';
				roomId = result.data.roomId;
				peerEcdhPublicKey = await importPublicKey(result.data.peerPublicKey);
				const sharedSecret = await deriveSharedSecret(myEcdhKeyPair!.privateKey, peerEcdhPublicKey);
				aesKey = await deriveAESKey(sharedSecret);

				activeSession.set({
					roomId,
					myEcdhKeyPair,
					myEcdsaKeyPair,
					myEcdsaPkBase64,
					peerEcdhPublicKey,
					aesKey,
					myNickname,
					partnerNickname: '',
					memberId,
					socket: null
				});

				goto(resolve('/chatroom'));
				return;
			}

			throw new Error(result.responseMessage || 'Gagal memulai random match.');
		} catch (err: unknown) {
			const errMsg =
				err instanceof Error ? err.message : 'Error memulai random match, silakan coba lagi.';
			const query = new URLSearchParams({ error: errMsg }).toString();
			goto(resolve(`/?${query}`));
		}
	}

	async function handleCancelMatchmaking() {
		if (socket) {
			socket.close();
		}
		if (matchQueueId) {
			try {
				await fetch(`${BASE_URL}/match/queue/${matchQueueId}`, { method: 'DELETE' });
			} catch (err) {
				console.error('Error cancelling match queue:', err);
			}
		}
		goto(resolve('/'));
	}

	onMount(() => {
		handleStartMatchmaking();

		return () => {
			// Clean up socket if we leave the page without matching
			if (socket && !aesKey) {
				socket.close();
			}
		};
	});
</script>

<svelte:head>
	<title>UMBRA — Random Match Waiting</title>
	<link
		href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Space+Grotesk:wght@500;700&family=Inter:wght@400;500;600&display=swap"
		rel="stylesheet"
	/>
	<link
		rel="stylesheet"
		href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200"
	/>
</svelte:head>

<div
	class="h-screen flex flex-col bg-[#F4F6F9] text-slate-800 font-sans select-none overflow-hidden"
>
	<!-- Header -->
	<header
		class="h-16 flex-shrink-0 flex items-center justify-between px-8 border-b border-slate-200/50 bg-white shadow-sm"
	>
		<!-- Title -->
		<div class="flex items-center space-x-2.5">
			<img src="/logo.webp" alt="UMBRA Logo" class="w-8 h-8 object-contain rounded-lg" />
			<span class="font-bold text-lg text-[#00658D] tracking-tight font-['Space_Grotesk']"
				>UMBRA</span
			>
		</div>

		<!-- Status Badge -->
		<div
			class="flex items-center space-x-2 bg-[#E6F4EA] text-[#137333] px-3.5 py-1.5 rounded-full text-xs font-semibold"
		>
			<span class="w-1.5 h-1.5 rounded-full bg-[#137333]"></span>
			<span class="font-mono text-[11px] uppercase tracking-wider font-bold">E2EE Ready</span>
		</div>
	</header>

	<!-- Main Content Area -->
	<div
		class="flex-grow flex flex-col items-center justify-center p-6 max-w-lg mx-auto w-full text-center space-y-12"
	>
		<!-- Avatars & Line Connection -->
		<div class="flex items-center justify-between w-full max-w-sm relative">
			<!-- Connection Line -->
			<div class="absolute left-16 right-16 top-10 -translate-y-1/2 h-[2px] bg-slate-200">
				<!-- Glowing Dot Flow Animation -->
				<div class="absolute top-0 bottom-0 w-4 bg-[#00aeef] rounded-full animate-flow-dot"></div>
			</div>

			<!-- Left Avatar (ANDA) -->
			<div class="flex flex-col items-center space-y-3 z-10">
				<div
					class="w-20 h-20 rounded-full border-4 border-[#00aeef] bg-white flex items-center justify-center shadow-lg relative animate-pulse-glow"
				>
					<!-- Stylized Profile Avatar SVG -->
					<svg
						class="w-10 h-10 text-[#00658D]"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
					>
						<path
							d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"
							stroke-linecap="round"
							stroke-linejoin="round"
						></path>
						<circle cx="12" cy="7" r="4" stroke-linecap="round" stroke-linejoin="round"></circle>
					</svg>
					<!-- Blue connector dot on the right edge -->
					<div
						class="absolute right-[-6px] top-1/2 -translate-y-1/2 w-2.5 h-2.5 rounded-full bg-[#00aeef]"
					></div>
				</div>
				<span class="text-xs font-bold text-[#00aeef] tracking-widest uppercase font-mono"
					>ANDA</span
				>
			</div>

			<!-- Right Avatar (MENCARI) -->
			<div class="flex flex-col items-center space-y-3 z-10">
				<div
					class="w-20 h-20 rounded-full bg-[#F1F5F9] border-4 border-white flex items-center justify-center text-slate-400 font-bold text-3xl shadow-sm animate-pulse"
				>
					?
				</div>
				<span class="text-xs font-bold text-slate-400 tracking-widest uppercase font-mono"
					>MENCARI...</span
				>
			</div>
		</div>

		<!-- Message -->
		<div class="space-y-3">
			<h3 class="text-xl font-bold text-slate-900">Mencari lawan bicara...</h3>
			<p class="text-sm text-slate-500 leading-relaxed">
				Sistem kami sedang mengenkripsi jalur komunikasi<br />untuk sesi anonim Anda.
			</p>
		</div>

		<!-- Infinite Loader Line -->
		<div class="w-full max-w-sm h-1.5 bg-slate-200/60 rounded-full overflow-hidden relative">
			<div
				class="absolute top-0 bottom-0 bg-gradient-to-r from-[#00aeef] to-[#00658d] w-1/3 rounded-full animate-infinite-slide"
			></div>
		</div>

		<!-- Stats Banner (Figma Capsule style) -->
		<div
			class="w-full bg-[#0F1C2C] text-[#CBD5E1] py-3.5 px-6 rounded-full flex items-center justify-between text-[11px] font-bold font-mono shadow-md"
		>
			<div class="flex items-center space-x-2">
				<span class="w-2 h-2 rounded-full bg-[#10B981]"></span>
				<span>1.240 pengguna online</span>
			</div>
			<span class="opacity-30">||</span>
			<div class="flex items-center space-x-2">
				<span class="w-2.5 h-2.5 rounded-full bg-[#00aeef]"></span>
				<span>Enkripsi aktif</span>
			</div>
			<span class="opacity-30">||</span>
			<div class="flex items-center space-x-1.5 flex-row">
				<span class="material-symbols-outlined text-[16px] text-[#EF4444] leading-none">timer</span>
				<span>&lt; 30 detik</span>
			</div>
		</div>

		<!-- Cancel Button -->
		<div class="space-y-4 pt-4">
			<button
				onclick={handleCancelMatchmaking}
				class="px-10 py-4 bg-[#0A6C8B] text-white font-bold rounded-xl shadow-md hover:bg-[#085A75] hover:shadow-lg transition-all active:scale-[0.98] text-sm uppercase tracking-wider cursor-pointer"
			>
				Batalkan pencarian
			</button>
			<p class="text-[10px] font-bold text-slate-400 tracking-widest font-mono uppercase">
				EST. WAIT: &lt; 30S
			</p>
		</div>
	</div>
</div>

<style>
	/* Animations */
	@keyframes flow-dot {
		0% {
			left: 0%;
			opacity: 0;
		}
		10% {
			opacity: 1;
		}
		90% {
			opacity: 1;
		}
		100% {
			left: 100%;
			opacity: 0;
		}
	}
	@keyframes infinite-slide {
		0% {
			left: -33%;
		}
		100% {
			left: 100%;
		}
	}
	:global(.animate-flow-dot) {
		animation: flow-dot 2.5s infinite linear;
	}
	:global(.animate-infinite-slide) {
		animation: infinite-slide 1.5s infinite linear;
	}
</style>
