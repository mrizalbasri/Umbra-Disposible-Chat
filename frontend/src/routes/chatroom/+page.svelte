<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import {
		generateKeyPair,
		exportPublicKey,
		importPublicKey,
		deriveSharedSecret,
		deriveAESKey,
		encrypt,
		decrypt,
		exportSigningPublicKey,
		importSigningPublicKey,
		sign,
		verify
	} from '$lib/crypto';

	import { env } from '$env/dynamic/public';
	import { activeSession, type ActiveSession } from '$lib/sessionStore';

	// API and WS Configuration (dynamic with fallback to localhost)
	const BACKEND_HOST = env.PUBLIC_BACKEND_URL || 'http://localhost:8080';
	const BASE_URL = `${BACKEND_HOST}/v1/api`;
	const WS_BASE_URL = `${BACKEND_HOST.replace(/^http/, 'ws')}/ws`;

	// roomState: set to 'setup' by default to render the real flow
	let roomState = $state<'setup' | 'waiting' | 'match_waiting' | 'chat'>('setup');
	let isSidebarOpen = $state(false);
	let nicknameInput = $state('');
	let roomCodeInput = $state('');
	let matchQueueId = $state('');

	// Interactive features state
	let showSearchInput = $state(false);
	let searchQuery = $state('');
	let showMoreMenu = $state(false);

	// Session variables
	let myNickname = $state('');
	let roomCode = $state('');
	let roomId = $state('');
	let memberId = $state('');
	let partnerNickname = $state('');

	// Deriving user initials dynamically
	let myInitials = $derived(myNickname ? myNickname.slice(0, 2).toUpperCase() : 'ME');
	let partnerInitials = $derived(
		partnerNickname ? partnerNickname.slice(0, 2).toUpperCase() : '..'
	);

	// Cryptographic keys
	let myEcdhKeyPair: CryptoKeyPair | null = null;
	let myEcdsaKeyPair: CryptoKeyPair | null = null;
	let myEcdsaPkBase64 = $state(''); // stores ECDSA public key in base64
	let peerEcdhPublicKey: CryptoKey | null = null;
	let peerEcdsaPublicKey: CryptoKey | null = null;
	let aesKey: CryptoKey | null = $state(null); // start with null for real E2EE

	// Chat room variables
	let chatContainer: HTMLDivElement | null = $state(null);
	let textareaElement: HTMLTextAreaElement | null = $state(null);
	let messageInput = $state('');
	let showToast = $state(false);
	let toastMessage = $state('');

	interface ChatMessage {
		id: string;
		type: 'system' | 'partner' | 'self';
		text?: string;
		boxed?: boolean;
		sender?: string;
		time?: string;
	}

	// Initial clean messages array
	let messages = $state<ChatMessage[]>([
		{
			id: 'm1',
			type: 'system',
			text: 'Session established. Messages will self-destruct on disconnect.',
			boxed: true
		}
	]);

	// Filter messages by search query in real-time
	let filteredMessages = $derived(
		searchQuery.trim()
			? messages.filter((m) => m.text && m.text.toLowerCase().includes(searchQuery.toLowerCase()))
			: messages
	);

	// WebSocket connection state
	let socket: WebSocket | null = null;
	let isConnected = $state(true); // default true for mock simulation

	// TTL Timer
	let timeLeft = $state(899); // 14:59 in seconds
	let timerInterval: ReturnType<typeof setInterval> | undefined = undefined;

	// Start TTL countdown
	function startTimer() {
		timeLeft = 899;
		if (timerInterval !== undefined) clearInterval(timerInterval);
		timerInterval = setInterval(() => {
			if (timeLeft > 0) {
				timeLeft -= 1;
			} else {
				if (timerInterval !== undefined) {
					clearInterval(timerInterval);
					timerInterval = undefined;
				}
				toastMessage = 'Sesi obrolan telah berakhir (Waktu TTL Habis)!';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 4000);
				leaveRoom();
			}
		}, 1000);
	}

	// Stop TTL countdown
	function stopTimer() {
		if (timerInterval !== undefined) {
			clearInterval(timerInterval);
			timerInterval = undefined;
		}
	}

	function formatTime(seconds: number): string {
		const m = Math.floor(seconds / 60);
		const s = seconds % 60;
		return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}s`;
	}

	function createClientMemberID(): string {
		if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
			return crypto.randomUUID();
		}
		return `${Date.now()}-${Math.random().toString(36).slice(2)}`;
	}

	async function transitionToMatchedRoom(peerPublicKey: string, myEcdhPkBase64: string) {
		peerEcdhPublicKey = await importPublicKey(peerPublicKey);
		const sharedSecret = await deriveSharedSecret(myEcdhKeyPair!.privateKey, peerEcdhPublicKey);
		aesKey = await deriveAESKey(sharedSecret);

		messages = [
			{
				id: 'init-sys',
				type: 'system',
				text: 'Session established. Messages will self-destruct on disconnect.',
				boxed: true
			}
		];

		roomState = 'chat';
		startTimer();
		connectWebSocket(myEcdhPkBase64);
	}

	// Auto-growing input box
	$effect(() => {
		if (textareaElement && messageInput !== undefined) {
			textareaElement.style.height = 'auto';
			textareaElement.style.height = `${textareaElement.scrollHeight}px`;
		}
	});

	// Scroll chat container to bottom
	function scrollToBottom() {
		if (chatContainer) {
			chatContainer.scrollTop = chatContainer.scrollHeight;
		}
	}

	// Generate ECDH & ECDSA key pairs sharing the same key material
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

	// Bind WebSocket handlers to the active socket
	function bindWebSocketHandlers() {
		if (!socket) return;

		socket.onopen = () => {
			isConnected = true;

			if (aesKey) {
				sendHandshake();
			}
		};

		socket.onmessage = async (event) => {
			try {
				const data = JSON.parse(event.data);
				if (data.roomId && data.roomId !== roomId) return;

				if (data.event === 'error') {
					if (data.code === 16 || data.payload?.code === 16) {
						toastMessage = 'Security Alert: Tanda tangan pesan tidak valid (ECDSA Verify Failed)!';
						showToast = true;
						setTimeout(() => {
							showToast = false;
						}, 4000);
					} else {
						console.error('Backend error:', data.message || data.payload?.message);
					}
					return;
				}

				if (data.event === 'peer_joined') {
					peerEcdhPublicKey = await importPublicKey(data.payload.publicKey);

					const sharedSecret = await deriveSharedSecret(
						myEcdhKeyPair!.privateKey,
						peerEcdhPublicKey
					);
					aesKey = await deriveAESKey(sharedSecret);

					roomState = 'chat';
					startTimer();
					await tick();
					scrollToBottom();

					await sendHandshake();
				} else if (data.event === 'message') {
					if (data.senderMemberId === memberId) return;
					const { ciphertext, iv, signature, timestamp } = data.payload;

					if (!aesKey) return;

					const decryptedText = await decrypt(ciphertext, iv, aesKey);
					const decryptedObj = JSON.parse(decryptedText);

					if (decryptedObj.type === 'handshake') {
						peerEcdsaPublicKey = await importSigningPublicKey(decryptedObj.signingPublicKey);
						partnerNickname = decryptedObj.nickname; // ponytail: store partner name dynamically
						messages = [
							...messages,
							{
								id: Date.now().toString(),
								type: 'system',
								text: `${decryptedObj.nickname} joined the session.`,
								boxed: true
							}
						];
						await tick();
						scrollToBottom();
					} else if (decryptedObj.type === 'chat') {
						if (peerEcdsaPublicKey) {
							const msgData = new TextEncoder().encode(ciphertext + iv + timestamp);
							const isValid = await verify(msgData.buffer, signature, peerEcdsaPublicKey);
							if (!isValid) {
								console.error('ECDSA signature verification failed! Message untrusted.');
								return;
							}
						}

						if (data.senderMemberId === memberId) return;

						const timeObj = new Date(timestamp);
						const timeString = timeObj.toLocaleTimeString([], {
							hour: '2-digit',
							minute: '2-digit'
						});

						messages = [
							...messages,
							{
								id: Date.now().toString(),
								type: 'partner',
								sender: decryptedObj.nickname.substring(0, 2).toUpperCase(),
								text: decryptedObj.text,
								time: timeString
							}
						];
						await tick();
						scrollToBottom();
					}
				} else if (data.event === 'peer_left' || data.event === 'room_destroyed') {
					messages = [
						...messages,
						{
							id: Date.now().toString(),
							type: 'system',
							text: 'Partner left. Session destroyed.',
							boxed: true
						}
					];
					isConnected = false;
					stopTimer();

					// Clear all private keys and session data from browser memory
					aesKey = null;
					myEcdhKeyPair = null;
					myEcdsaKeyPair = null;
					myEcdsaPkBase64 = '';
					peerEcdhPublicKey = null;
					peerEcdsaPublicKey = null;

					await tick();
					scrollToBottom();
				}
			} catch (err) {
				console.error('Error handling websocket message:', err);
			}
		};

		socket.onclose = () => {
			isConnected = false;
			stopTimer();
		};
	}

	// Connect to WebSocket and bind message handlers (deferred hook)
	function connectWebSocket(myEcdhPkBase64: string) {
		const wsUrl = `${WS_BASE_URL}?roomId=${roomId}&publicKey=${encodeURIComponent(myEcdhPkBase64)}`;
		socket = new WebSocket(wsUrl);
		bindWebSocketHandlers();
	}

	// Encrypted Handshake over WebSocket
	async function sendHandshake() {
		if (!aesKey || !myEcdsaKeyPair || !socket) return;

		const payload = {
			type: 'handshake',
			signingPublicKey: myEcdsaPkBase64,
			nickname: myNickname
		};

		const plaintext = JSON.stringify(payload);
		const { ciphertext, iv } = await encrypt(plaintext, aesKey);
		const timestamp = new Date().toISOString();

		const dataToSign = new TextEncoder().encode(ciphertext + iv + timestamp);
		const signature = await sign(dataToSign.buffer, myEcdsaKeyPair.privateKey);

		const msg = {
			event: 'send_message',
			roomId,
			senderMemberId: memberId,
			payload: {
				ciphertext,
				iv,
				signature,
				timestamp
			}
		};

		const s = socket;
		if (s.readyState === WebSocket.OPEN) {
			s.send(JSON.stringify(msg));
		} else if (s.readyState === WebSocket.CONNECTING) {
			s.addEventListener(
				'open',
				() => {
					s.send(JSON.stringify(msg));
				},
				{ once: true }
			);
		}
	}

	// Call POST /v1/api/room/create (deferred hook)
	async function handleCreateRoom() {
		if (!nicknameInput.trim()) {
			return;
		}
		try {
			await generateKeys();
			const myEcdhPkBase64 = await exportPublicKey(myEcdhKeyPair!.publicKey);

			const res = await fetch(`${BASE_URL}/room/create`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ publicKey: myEcdhPkBase64, nickname: nicknameInput.trim() })
			});

			const result = await res.json();
			if (result.responseCode !== '00') {
				throw new Error(result.responseMessage || 'Failed to create room');
			}

			myNickname = nicknameInput.trim();
			roomCode = result.data.roomCode;
			roomId = result.data.roomId;
			memberId = createClientMemberID();

			messages = [
				{
					id: 'init-sys',
					type: 'system',
					text: 'Session established. Messages will self-destruct on disconnect.',
					boxed: true
				}
			];

			roomState = 'waiting';
			window.history.replaceState({}, '', '/chatroom'); // ponytail: clear url params
			connectWebSocket(myEcdhPkBase64);
		} catch (err: unknown) {
			const errMsg =
				err instanceof Error
					? err.message
					: 'Network error, please start/check the backend server.';
			goto(resolve(`/create?error=${encodeURIComponent(errMsg)}`));
		}
	}

	// Call POST /v1/api/room/join (deferred hook)
	async function handleJoinRoom() {
		if (!roomCodeInput.trim() || roomCodeInput.length < 7) {
			return;
		}
		if (!nicknameInput.trim()) {
			return;
		}
		try {
			await generateKeys();
			const myEcdhPkBase64 = await exportPublicKey(myEcdhKeyPair!.publicKey);

			const res = await fetch(`${BASE_URL}/room/join`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					roomCode: roomCodeInput.trim(),
					publicKey: myEcdhPkBase64,
					nickname: nicknameInput.trim()
				})
			});

			const result = await res.json();
			if (result.responseCode !== '00') {
				throw new Error(result.responseMessage || 'Failed to join room');
			}

			myNickname = nicknameInput.trim();
			roomCode = roomCodeInput.trim();
			roomId = result.data.roomId;
			memberId = createClientMemberID();

			// Kunci ECDH peer untuk enkripsi
			peerEcdhPublicKey = await importPublicKey(result.data.peerPublicKey);

			const sharedSecret = await deriveSharedSecret(myEcdhKeyPair!.privateKey, peerEcdhPublicKey);
			aesKey = await deriveAESKey(sharedSecret);

			messages = [
				{
					id: 'init-sys',
					type: 'system',
					text: 'Session established. Messages will self-destruct on disconnect.',
					boxed: true
				}
			];

			roomState = 'chat';
			window.history.replaceState({}, '', '/chatroom'); // ponytail: clear url params
			startTimer();
			connectWebSocket(myEcdhPkBase64);
		} catch (err: unknown) {
			const errMsg =
				err instanceof Error
					? err.message
					: 'Network error, please start/check the backend server.';
			goto(resolve(`/join?error=${encodeURIComponent(errMsg)}`));
		}
	}

	// Start Random Matchmaking Queue
	async function handleStartMatchmaking() {
		let finalNickname = nicknameInput.trim();
		if (!finalNickname) {
			finalNickname = `Guest_${Math.floor(1000 + Math.random() * 9000)}`;
		}

		try {
			await generateKeys();
			const myEcdhPkBase64 = await exportPublicKey(myEcdhKeyPair!.publicKey);
			myNickname = finalNickname;
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
				roomState = 'match_waiting';
				window.history.replaceState({}, '', '/chatroom'); // ponytail: clear url params
				connectWebSocket(myEcdhPkBase64);
				return;
			}

			if (result.responseCode === '00' && result.data?.status === 'matched') {
				matchQueueId = '';
				roomId = result.data.roomId;
				window.history.replaceState({}, '', '/chatroom'); // ponytail: clear url params
				await transitionToMatchedRoom(result.data.peerPublicKey, myEcdhPkBase64);
				return;
			}

			throw new Error(result.responseMessage || 'Gagal memulai random match.');
		} catch (err: unknown) {
			const errMsg =
				err instanceof Error ? err.message : 'Error memulai random match, silakan coba lagi.';
			goto(resolve(`/?error=${encodeURIComponent(errMsg)}`));
		}
	}

	// Cancel matchmaking queue
	async function handleCancelMatchmaking() {
		if (matchQueueId) {
			try {
				await fetch(`${BASE_URL}/match/queue/${matchQueueId}`, { method: 'DELETE' });
			} catch (err) {
				console.error('Error cancelling random match queue:', err);
			}
		}
		matchQueueId = '';
		if (socket) {
			socket.close();
			socket = null;
		}

		const params = new URLSearchParams(window.location.search);
		if (params.get('tab') === 'match') {
			await goto(resolve('/'));
		} else {
			roomState = 'setup';
		}
	}

	// Send message handler (works locally in mock mode)
	async function sendMessage() {
		const text = messageInput.trim();
		if (!text) return;

		// Append self message locally
		const timeString = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
		messages = [
			...messages,
			{
				id: Date.now().toString(),
				type: 'self',
				sender: 'ME',
				text,
				time: timeString
			}
		];

		// If connected to WebSocket, encrypt and transmit
		if (socket && socket.readyState === WebSocket.OPEN && aesKey && myEcdsaKeyPair) {
			try {
				const payload = {
					type: 'chat',
					nickname: myNickname,
					text
				};
				const plaintext = JSON.stringify(payload);
				const { ciphertext, iv } = await encrypt(plaintext, aesKey);
				const timestamp = new Date().toISOString();

				const dataToSign = new TextEncoder().encode(ciphertext + iv + timestamp);
				const signature = await sign(dataToSign.buffer, myEcdsaKeyPair.privateKey);

				const msg = {
					event: 'send_message',
					roomId,
					senderMemberId: memberId,
					payload: {
						ciphertext,
						iv,
						signature,
						timestamp
					}
				};

				socket.send(JSON.stringify(msg));
			} catch (err) {
				console.error('Error sending encrypted socket message:', err);
			}
		}

		messageInput = '';
		await tick();
		scrollToBottom();
	}

	// Handle Enter keypress in textarea
	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			sendMessage();
		}
	}

	// Copy Room Code to clipboard
	function copyCode() {
		navigator.clipboard.writeText(roomCode).then(() => {
			toastMessage = 'Room code copied to clipboard';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 2000);
		});
	}

	// Leave the room and reset connection states
	function leaveRoom() {
		if (socket) {
			socket.close();
		}
		socket = null;
		isConnected = false;
		roomState = 'setup';
		messages = [];

		// Clear all private keys and session data from browser memory
		aesKey = null;
		myEcdhKeyPair = null;
		myEcdsaKeyPair = null;
		myEcdsaPkBase64 = '';
		peerEcdhPublicKey = null;
		peerEcdsaPublicKey = null;

		nicknameInput = '';
		roomCodeInput = '';
		matchQueueId = '';
		stopTimer();
		goto(resolve('/')); // ponytail: redirect to home to reset query params
	}

	// Reset messages (used for Clear Chat in dropdown)
	function clearChat() {
		messages = [
			{
				id: 'init-sys',
				type: 'system',
				text: 'Session established. Messages will self-destruct on disconnect.',
				boxed: true
			}
		];
	}

	// Initialize view based on query tab
	onMount(() => {
		let session: ActiveSession | null = null;
		activeSession.subscribe((val) => {
			session = val;
		})();

		const activeSessionData = session as ActiveSession | null;
		if (activeSessionData) {
			roomId = activeSessionData.roomId;
			myEcdhKeyPair = activeSessionData.myEcdhKeyPair;
			myEcdsaKeyPair = activeSessionData.myEcdsaKeyPair;
			myEcdsaPkBase64 = activeSessionData.myEcdsaPkBase64;
			peerEcdhPublicKey = activeSessionData.peerEcdhPublicKey;
			aesKey = activeSessionData.aesKey;
			myNickname = activeSessionData.myNickname;
			memberId = activeSessionData.memberId;
			roomState = 'chat';

			// Connect or reuse WebSocket connection
			if (activeSessionData.socket) {
				socket = activeSessionData.socket;
				isConnected = true;
				bindWebSocketHandlers();
				if (aesKey) {
					sendHandshake();
				}
			} else {
				exportPublicKey(myEcdhKeyPair!.publicKey).then((myEcdhPkBase64) => {
					connectWebSocket(myEcdhPkBase64);
				});
			}

			messages = [
				{
					id: 'init-sys',
					type: 'system',
					text: 'Session established. Messages will self-destruct on disconnect.',
					boxed: true
				}
			];
			startTimer();
			activeSession.set(null); // clear session store so fresh reload redirects out
		} else {
			const params = new URLSearchParams(window.location.search);
			const tab = params.get('tab');
			if (tab === 'match') {
				nicknameInput = `Guest_${Math.floor(1000 + Math.random() * 9000)}`;
				handleStartMatchmaking(); // ponytail: bypass nickname form and start match directly
			} else if (tab === 'create') {
				const nickname = params.get('nickname');
				if (nickname) {
					nicknameInput = nickname;
					handleCreateRoom(); // ponytail: auto-create room if redirected from /create
				}
			} else if (tab === 'join') {
				const roomCode = params.get('roomCode');
				const nickname = params.get('nickname');
				if (roomCode && nickname) {
					roomCodeInput = roomCode;
					nicknameInput = nickname;
					handleJoinRoom(); // ponytail: auto-join room if redirected from /join
				}
			}
		}

		scrollToBottom();

		return () => {
			stopTimer();
			if (socket) {
				socket.close();
			}
		};
	});
</script>

<div class="h-screen bg-background select-none">
	{#if roomState === 'setup'}
		<!-- INITIALIZING SECURITY / REDIRECTING SCREEN -->
		<div class="flex flex-col items-center justify-center h-full space-y-4 bg-white">
			<div
				class="w-10 h-10 border-4 border-[#00aeef] border-t-transparent rounded-full animate-spin"
			></div>
			<p class="text-sm font-label-mono text-outline uppercase tracking-wider">
				Inisialisasi Keamanan...
			</p>
		</div>
	{:else if roomState === 'match_waiting'}
		<!-- RANDOM MATCH WAITING SCREEN -->
		<div class="flex flex-col h-full bg-[#F4F6F9] text-on-background relative">
			<!-- Header -->
			<header
				class="h-16 flex-shrink-0 flex items-center justify-between px-8 border-b border-outline-variant/30 bg-white"
			>
				<!-- Title (No Back Button in Figma Header) -->
				<div class="flex items-center space-x-2.5">
					<img src="/logo.webp" alt="UMBRA Logo" class="w-8 h-8 object-contain rounded-lg" />
					<span class="font-bold text-lg text-[#00658D] tracking-tight font-['Space_Grotesk']"
						>UMBRA</span
					>
				</div>

				<!-- Status Badge -->
				<div
					class="flex items-center space-x-2 bg-[#E6F4EA] text-[#137333] px-3.5 py-1.5 rounded-full text-xs font-semibold select-none"
				>
					<span class="w-1.5 h-1.5 rounded-full bg-[#137333]"></span>
					<span class="font-label-mono text-[11px] uppercase tracking-wider font-bold"
						>E2EE Ready</span
					>
				</div>
			</header>

			<!-- Main Content Area -->
			<div
				class="flex-grow flex flex-col items-center justify-center p-6 max-w-lg mx-auto w-full text-center space-y-12"
			>
				<!-- Avatars & Line Connection -->
				<div class="flex items-center justify-between w-full max-w-sm relative">
					<!-- Connection Line -->
					<div class="absolute left-16 right-16 top-10 -translate-y-1/2 h-[2px] bg-gray-200">
						<!-- Glowing Dot Flow Animation -->
						<div
							class="absolute top-0 bottom-0 w-4 bg-[#00aeef] rounded-full animate-flow-dot"
						></div>
					</div>

					<!-- Left Avatar (ANDA) -->
					<div class="flex flex-col items-center space-y-3 z-10">
						<div
							class="w-20 h-20 rounded-full border-4 border-[#00aeef] bg-white flex items-center justify-center shadow-premium relative select-none"
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
								<circle cx="12" cy="7" r="4" stroke-linecap="round" stroke-linejoin="round"
								></circle>
							</svg>
							<!-- Blue connector dot on the right edge -->
							<div
								class="absolute right-[-6px] top-1/2 -translate-y-1/2 w-2.5 h-2.5 rounded-full bg-[#00aeef]"
							></div>
						</div>
						<span class="text-xs font-bold text-[#00aeef] tracking-widest uppercase font-label-mono"
							>ANDA</span
						>
					</div>

					<!-- Right Avatar (MENCARI) -->
					<div class="flex flex-col items-center space-y-3 z-10">
						<div
							class="w-20 h-20 rounded-full bg-[#F1F5F9] border-4 border-white flex items-center justify-center text-gray-400 font-bold text-3xl shadow-sm animate-pulse select-none"
						>
							?
						</div>
						<span class="text-xs font-bold text-gray-400 tracking-widest uppercase font-label-mono"
							>MENCARI...</span
						>
					</div>
				</div>

				<!-- Message -->
				<div class="space-y-3">
					<h3 class="text-xl font-bold text-on-surface">Mencari lawan bicara...</h3>
					<p class="text-sm text-outline leading-relaxed">
						Sistem kami sedang mengenkripsi jalur komunikasi<br />untuk sesi anonim Anda.
					</p>
				</div>

				<!-- Infinite Loader Line -->
				<div class="w-full max-w-sm h-1.5 bg-gray-200/60 rounded-full overflow-hidden relative">
					<div
						class="absolute top-0 bottom-0 bg-gradient-to-r from-[#00aeef] to-[#00658d] w-1/3 rounded-full animate-infinite-slide"
					></div>
				</div>

				<!-- Stats Banner (Figma Capsule style) -->
				<div
					class="w-full bg-[#0F1C2C] text-[#CBD5E1] py-3.5 px-6 rounded-full flex items-center justify-between text-[11px] font-bold shadow-premium font-label-mono select-none"
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
					<div class="flex items-center space-x-1.5">
						<span class="material-symbols-outlined text-[16px] text-[#EF4444] leading-none"
							>timer</span
						>
						<span>&lt; 30 detik</span>
					</div>
				</div>

				<!-- Cancel Button -->
				<div class="space-y-4 pt-4">
					<button
						onclick={handleCancelMatchmaking}
						class="px-10 py-4 bg-[#0A6C8B] text-white font-bold rounded-xl shadow-premium hover:bg-[#085A75] hover:shadow-lg transition-all active:scale-[0.98] text-sm uppercase tracking-wider cursor-pointer"
					>
						Batalkan pencarian
					</button>
					<p class="text-[10px] font-bold text-outline tracking-widest font-label-mono uppercase">
						EST. WAIT: &lt; 30S
					</p>
				</div>
			</div>
		</div>
	{:else}
		<!-- CHAT ROOM INTERFACE -->
		<div class="flex h-screen overflow-hidden w-full bg-surface-container-lowest relative">
			{#if isSidebarOpen}
				<!-- svelte-ignore a11y_click_events_have_key_events -->
				<!-- svelte-ignore a11y_no_static_element_interactions -->
				<div
					class="fixed inset-0 bg-black/40 z-20 md:hidden animate-fade-in"
					onclick={() => (isSidebarOpen = false)}
				></div>
			{/if}

			<!-- Sidebar (Left, 300px) -->
			<aside
				class="w-[300px] flex-shrink-0 bg-surface-container-low border-r border-outline-variant/20 flex flex-col justify-between sidebar-transition fixed md:relative z-30 h-full {isSidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}"
			>
				<div class="p-gutter space-y-8">
					<!-- Brand Anchor -->
					<div class="flex items-center space-x-3">
						<img src="/logo.webp" alt="UMBRA Logo" class="w-10 h-10 object-contain rounded-xl" />
						<span class="font-bold text-2xl text-primary tracking-tight font-['Space_Grotesk']"
							>UMBRA</span
						>
					</div>

					<!-- Session Identity -->
					<div class="space-y-stack-md">
						<div class="flex flex-col space-y-2">
							<span
								class="font-label-mono text-sm text-outline uppercase tracking-widest text-[12px]"
								>Active Alias</span
							>
							<div
								class="flex items-center space-x-3 bg-surface-container-highest p-3 rounded-xl border border-outline-variant/30"
							>
								<div
									class="w-8 h-8 rounded-full bg-secondary-fixed flex items-center justify-center text-on-secondary-fixed font-semibold"
								>
									<span class="material-symbols-outlined text-[18px]">person</span>
								</div>
								<span class="font-medium text-on-surface">{myNickname}</span>
							</div>
						</div>

						<div class="flex flex-col space-y-2">
							<span
								class="font-label-mono text-sm text-outline uppercase tracking-widest text-[12px]"
								>Room Code</span
							>
							<!-- svelte-ignore a11y_click_events_have_key_events -->
							<!-- svelte-ignore a11y_no_static_element_interactions -->
							<div
								class="group relative flex items-center justify-between bg-white p-3 rounded-xl border border-outline-variant/30 cursor-pointer hover:border-primary transition-colors"
								onclick={copyCode}
							>
								<span class="font-label-mono text-primary font-bold tracking-tighter" id="room-code"
									>{roomCode}</span
								>
								<span
									class="material-symbols-outlined text-outline group-hover:text-primary transition-colors text-[20px]"
									>content_copy</span
								>
							</div>
						</div>
					</div>

					<!-- Status Badges -->
					<div class="space-y-stack-sm pt-4 border-t border-outline-variant/20">
						<div class="flex items-center space-x-3 px-1">
							<div class="relative flex h-2 w-2">
								<span
									class="pulse-green absolute inline-flex h-full w-full rounded-full bg-tertiary-container opacity-75"
								></span>
								<span class="relative inline-flex rounded-full h-2 w-2 bg-tertiary-container"
								></span>
							</div>
							<span class="font-label-mono text-[12px] text-tertiary font-bold">
								{#if aesKey}E2EE ACTIVE{:else}E2EE WAITING{/if}
							</span>
						</div>
						<div class="flex items-center space-x-3 px-1">
							<span
								class="material-symbols-outlined {isConnected
									? 'text-primary'
									: 'text-outline'} text-[18px]">wifi_tethering</span
							>
							<span class="font-label-mono text-[12px] text-on-surface-variant uppercase">
								WS: {isConnected ? 'CONNECTED' : 'DISCONNECTED'}
							</span>
						</div>
						<div class="flex items-center space-x-3 px-1">
							<span class="material-symbols-outlined text-outline text-[18px]">timer</span>
							<span class="font-label-mono text-[12px] text-on-surface-variant"
								>TTL: {formatTime(timeLeft)}</span
							>
						</div>
					</div>
				</div>

				<!-- Footer Action -->
				<div class="p-gutter">
					<button
						onclick={leaveRoom}
						class="w-full flex items-center justify-center space-x-2 py-4 bg-[#EF4444]/10 text-[#EF4444] rounded-xl font-bold hover:bg-[#EF4444] hover:text-white transition-all duration-200 active:scale-[0.98]"
					>
						<span class="material-symbols-outlined">logout</span>
						<span>Leave Room</span>
					</button>
				</div>
			</aside>

			<!-- Main Area: Chat Content -->
			<main class="flex-grow flex flex-col relative bg-white">
				<!-- Chat Header (Contextual) -->
				<header
					class="h-16 flex items-center justify-between px-gutter border-b border-outline-variant/10 z-10 bg-white/80 backdrop-blur-md"
				>
					<div class="flex items-center space-x-2 md:space-x-4">
						<!-- Mobile Sidebar Toggle Button -->
						<button
							onclick={() => (isSidebarOpen = !isSidebarOpen)}
							class="md:hidden flex items-center justify-center p-2 rounded-lg bg-surface-container-low border border-outline-variant/30 text-primary hover:bg-surface-container-highest transition-colors cursor-pointer mr-1"
							aria-label="Toggle Sidebar"
						>
							<span class="material-symbols-outlined text-[20px] leading-none">menu</span>
						</button>

						<div class="flex -space-x-2">
							<div
								class="w-8 h-8 rounded-full border-2 border-white bg-primary-fixed flex items-center justify-center text-[10px] font-bold text-on-primary-fixed-variant select-none"
							>
								{myInitials}
							</div>
							{#if aesKey}
								<div
									class="w-8 h-8 rounded-full border-2 border-white bg-tertiary-fixed flex items-center justify-center text-[10px] font-bold text-on-tertiary-fixed-variant select-none"
								>
									{partnerInitials}
								</div>
							{/if}
						</div>
						<h2 class="font-bold text-lg text-on-surface font-['Space_Grotesk']">
							Anonymous Session <span class="text-outline-variant ml-2 font-normal font-sans"
								>#{roomId ? roomId.slice(0, 4).toUpperCase() : '...'}</span
							>
						</h2>
					</div>

					<div class="flex items-center space-x-2 relative">
						<!-- Toggle Search Input or Search Button (Right next to dropdown) -->
						{#if showSearchInput}
							<div
								class="flex items-center bg-surface-container-low px-3 py-1.5 rounded-lg border border-outline-variant/30 max-w-xs transition-all animate-in fade-in slide-in-from-right-2 duration-150 mr-1"
							>
								<span class="material-symbols-outlined text-[16px] text-outline mr-2">search</span>
								<input
									type="text"
									bind:value={searchQuery}
									placeholder="Search message..."
									class="bg-transparent border-none text-xs focus:ring-0 p-0 text-on-surface outline-none w-36 font-body-md"
								/>
								<button
									onclick={() => {
										showSearchInput = false;
										searchQuery = '';
									}}
									class="p-1 hover:text-red-500 transition-colors ml-1 flex items-center"
								>
									<span class="material-symbols-outlined text-[16px]">close</span>
								</button>
							</div>
						{:else}
							<button
								onclick={() => (showSearchInput = true)}
								class="p-2 text-on-surface-variant hover:bg-surface-container rounded-lg transition-colors"
							>
								<span class="material-symbols-outlined">search</span>
							</button>
						{/if}

						<!-- More Options Menu with Dropdown -->
						<div>
							<button
								onclick={() => (showMoreMenu = !showMoreMenu)}
								class="p-2 text-on-surface-variant hover:bg-surface-container rounded-lg transition-colors {showMoreMenu
									? 'bg-surface-container text-primary'
									: ''}"
							>
								<span class="material-symbols-outlined">more_vert</span>
							</button>
							{#if showMoreMenu}
								<!-- svelte-ignore a11y_click_events_have_key_events -->
								<!-- svelte-ignore a11y_no_static_element_interactions -->
								<div
									class="absolute right-0 mt-2 w-48 bg-white border border-outline-variant/30 rounded-xl shadow-premium py-2 z-50 animate-in fade-in slide-in-from-top-2 duration-150"
									onclick={() => (showMoreMenu = false)}
								>
									<button
										onclick={clearChat}
										class="w-full text-left px-4 py-2 text-sm text-on-surface hover:bg-surface-container transition-colors flex items-center space-x-2"
									>
										<span class="material-symbols-outlined text-[18px]">delete_sweep</span>
										<span>Clear Chat</span>
									</button>
									<button
										onclick={() => {
											toastMessage = 'Notifications muted';
											showToast = true;
											setTimeout(() => (showToast = false), 2000);
										}}
										class="w-full text-left px-4 py-2 text-sm text-on-surface hover:bg-surface-container transition-colors flex items-center space-x-2"
									>
										<span class="material-symbols-outlined text-[18px]">volume_off</span>
										<span>Mute Session</span>
									</button>
								</div>
							{/if}
						</div>
					</div>
				</header>

				<!-- ACTIVE CHAT AREA (Filtered by search) -->
				<div
					bind:this={chatContainer}
					class="flex-grow overflow-y-auto p-gutter space-y-6 chat-scroll"
					id="chat-container"
				>
					{#if filteredMessages.length === 0}
						<div
							class="flex flex-col items-center justify-center h-full text-outline py-8 text-center space-y-2"
						>
							<span class="material-symbols-outlined text-[36px]">search_off</span>
							<p class="text-sm">No messages match your search query.</p>
						</div>
					{:else}
						{#each filteredMessages as msg (msg.id)}
							{#if msg.type === 'system'}
								<div class="flex justify-center animate-message-fade-in">
									{#if msg.boxed}
										<p
											class="italic text-on-surface-variant/60 text-sm bg-surface-container-low px-4 py-1 rounded-full border border-outline-variant/10 text-center"
										>
											{msg.text}
										</p>
									{:else}
										<p class="italic text-on-surface-variant/60 text-sm text-center font-semibold">
											{msg.text}
										</p>
									{/if}
								</div>
							{:else if msg.type === 'partner'}
								<div class="flex items-end space-x-2 max-w-[80%] animate-message-fade-in">
									<div
										class="w-8 h-8 rounded-full bg-secondary-container flex-shrink-0 flex items-center justify-center text-[12px] font-bold text-on-secondary-container"
									>
										{msg.sender}
									</div>
									<div
										class="bg-[#F5F8FA] text-on-surface p-4 rounded-2xl shadow-sm bubble-partner"
									>
										<p class="text-base whitespace-pre-wrap select-text">{msg.text}</p>
										<span class="block mt-2 text-[10px] text-outline text-right">{msg.time}</span>
									</div>
								</div>
							{:else if msg.type === 'self'}
								<div class="flex items-end justify-end space-x-2 ml-auto max-w-[80%] animate-message-fade-in">
									<div
										class="bg-primary-container text-white p-4 rounded-2xl shadow-premium bubble-self"
									>
										<p class="text-base whitespace-pre-wrap select-text">{msg.text}</p>
										<span class="block mt-2 text-[10px] opacity-70 text-right">{msg.time}</span>
									</div>
									<div
										class="w-8 h-8 rounded-full bg-primary flex-shrink-0 flex items-center justify-center text-[12px] font-bold text-white font-label-mono"
									>
										{msg.sender}
									</div>
								</div>
							{/if}
						{/each}
					{/if}
				</div>

				<!-- Sticky Input Area -->
				<div class="p-gutter pt-2 bg-white">
					<div
						class="max-w-[1000px] mx-auto bg-white border border-outline-variant/30 rounded-2xl shadow-sm focus-within:ring-4 focus-within:ring-primary-container/10 focus-within:border-primary transition-all"
					>
						<div class="flex items-center p-2 border-b border-outline-variant/10 space-x-2">
							<button class="p-2 text-outline-variant hover:text-primary transition-colors">
								<span class="material-symbols-outlined">format_bold</span>
							</button>
							<button class="p-2 text-outline-variant hover:text-primary transition-colors">
								<span class="material-symbols-outlined">attach_file</span>
							</button>
							<button class="p-2 text-outline-variant hover:text-primary transition-colors">
								<span class="material-symbols-outlined">sentiment_satisfied</span>
							</button>
							<div class="h-4 w-[1px] bg-outline-variant/30 mx-2"></div>
							<span class="text-[10px] font-label-mono text-outline-variant uppercase hidden sm:inline"
								>Ephemeral Encryption Active</span
							>
						</div>
						<div class="flex items-end p-3 space-x-3">
							<textarea
								bind:this={textareaElement}
								bind:value={messageInput}
								onkeydown={handleKeyDown}
								class="flex-grow bg-transparent border-none focus:ring-0 resize-none max-h-32 text-base py-2 px-1 text-on-surface placeholder-outline-variant/60 outline-none"
								id="message-input"
								placeholder="Ketik pesan..."
								rows="1"></textarea>
							<button
								onclick={sendMessage}
								class="w-12 h-12 bg-primary-container rounded-xl flex items-center justify-center text-white shadow-premium hover:bg-primary transition-all active:scale-[0.98] flex-shrink-0 group"
								id="send-btn"
							>
								<span
									class="material-symbols-outlined group-hover:translate-x-0.5 transition-transform"
									>send</span
								>
							</button>
						</div>
					</div>
					<div class="flex justify-center py-4">
						<p
							class="text-[11px] font-label-mono text-outline-variant uppercase tracking-widest flex items-center"
						>
							<span class="material-symbols-outlined text-[14px] mr-1">history_toggle_off</span>
							Messages vanish upon session termination
						</p>
					</div>
				</div>
			</main>
		</div>
	{/if}
</div>

<!-- Copy Toast Notification -->
{#if showToast}
	<div
		class="fixed bottom-24 left-1/2 -translate-x-1/2 bg-inverse-surface text-white px-6 py-3 rounded-full text-sm font-medium shadow-lg z-50 animate-bounce"
	>
		{toastMessage}
	</div>
{/if}
