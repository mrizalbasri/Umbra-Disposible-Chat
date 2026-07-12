<script lang="ts">
	import { onMount, tick } from 'svelte';
	import {
		generateKeyPair,
		exportPublicKey,
		importPublicKey,
		deriveSharedSecret,
		deriveAESKey,
		encrypt,
		decrypt,
		generateSigningKeyPair,
		exportSigningPublicKey,
		importSigningPublicKey,
		sign,
		verify
	} from '$lib/crypto';

	// API and WS Configuration (for future integration)
	const BASE_URL = 'http://localhost:8080/v1/api';
	const WS_BASE_URL = 'ws://localhost:8080/ws';

	// roomState: set to 'setup' by default to render the real flow
	let roomState = $state<'setup' | 'waiting' | 'chat'>('setup');
	let activeTab = $state<'create' | 'join'>('create');
	let nicknameInput = $state('');
	let roomCodeInput = $state('');
	let isLoading = $state(false);
	let errorMessage = $state('');

	// Interactive features state
	let showSearchInput = $state(false);
	let searchQuery = $state('');
	let showMoreMenu = $state(false);

	// Session variables (initialized to mockup values)
	let myNickname = $state('Neon_Specter');
	let roomCode = $state('UX-882-KLA');
	let roomId = $state('550e8400-e29b-41d4-a716-446655440421');
	let memberId = $state('my-member-id');

	// Deriving user initials dynamically
	let myInitials = $derived(myNickname === 'Neon_Specter' ? 'ME' : myNickname.slice(0, 2).toUpperCase());

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

	// Initial mockup messages exactly matching the screenshot
	let messages = $state<any[]>([
		{
			id: 'm1',
			type: 'system',
			text: 'Session established. Messages will self-destruct on disconnect.',
			boxed: true
		},
		{
			id: 'm2',
			type: 'partner',
			sender: 'JS',
			text: 'Data integrity checks complete. Are we ready to initiate the transfer?',
			time: '10:42 AM'
		},
		{
			id: 'm3',
			type: 'system',
			text: 'Cipher updated by system.',
			boxed: false
		},
		{
			id: 'm4',
			type: 'self',
			sender: 'ME',
			text: 'Confirming. The keys have been rotated and I am monitoring the egress points now.',
			time: '10:43 AM'
		},
		{
			id: 'm5',
			type: 'partner',
			sender: 'JS',
			text: 'Understood. Initiating in T-minus 60 seconds. Keep the tunnel open.',
			time: '10:45 AM'
		},
		{
			id: 'm6',
			type: 'self',
			sender: 'ME',
			text: 'Monitoring heartbeat...',
			time: '10:45 AM'
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
	let timerInterval: any = null;

	// Start TTL countdown
	function startTimer() {
		timeLeft = 899;
		if (timerInterval) clearInterval(timerInterval);
		timerInterval = setInterval(() => {
			if (timeLeft > 0) {
				timeLeft -= 1;
			} else {
				clearInterval(timerInterval);
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
		if (timerInterval) {
			clearInterval(timerInterval);
			timerInterval = null;
		}
	}

	function formatTime(seconds: number): string {
		const m = Math.floor(seconds / 60);
		const s = seconds % 60;
		return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}s`;
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

	// Auto-format Room Code to XXXX-XX format
	function handleRoomCodeInput(e: Event) {
		const input = e.target as HTMLInputElement;
		let val = input.value.replace(/[^a-zA-Z0-9]/g, '').toUpperCase().slice(0, 6);
		if (val.length > 4) {
			val = val.slice(0, 4) + '-' + val.slice(4);
		}
		roomCodeInput = val;
		input.value = val;
	}

	// Generate ECDH & ECDSA key pairs
	async function generateKeys() {
		myEcdhKeyPair = await generateKeyPair();
		myEcdsaKeyPair = await generateSigningKeyPair();
		myEcdsaPkBase64 = await exportSigningPublicKey(myEcdsaKeyPair.publicKey);
	}

	// Connect to WebSocket and bind message handlers (deferred hook)
	function connectWebSocket(myEcdhPkBase64: string) {
		const wsUrl = `${WS_BASE_URL}?roomId=${roomId}&memberId=${memberId}&ecdhPublicKey=${encodeURIComponent(myEcdhPkBase64)}&ecdsaPublicKey=${encodeURIComponent(myEcdsaPkBase64)}`;
		socket = new WebSocket(wsUrl);

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
					peerEcdhPublicKey = await importPublicKey(data.payload.ecdhPublicKey);
					peerEcdsaPublicKey = await importSigningPublicKey(data.payload.ecdsaPublicKey);

					const sharedSecret = await deriveSharedSecret(myEcdhKeyPair!.privateKey, peerEcdhPublicKey);
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
						const timeString = timeObj.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

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

	// Encrypted Handshake over WebSocket
	async function sendHandshake() {
		if (!aesKey || !myEcdsaKeyPair) return;

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
				timestamp,
				publicKey: myEcdsaPkBase64
			}
		};

		socket?.send(JSON.stringify(msg));
	}

	// Call POST /v1/api/room/create (deferred hook)
	async function handleCreateRoom() {
		if (!nicknameInput.trim()) {
			errorMessage = 'Nickname cannot be empty';
			return;
		}
		isLoading = true;
		errorMessage = '';
		try {
			await generateKeys();
			const myEcdhPkBase64 = await exportPublicKey(myEcdhKeyPair!.publicKey);

			const res = await fetch(`${BASE_URL}/room/create`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ecdhPublicKey: myEcdhPkBase64, ecdsaPublicKey: myEcdsaPkBase64 })
			});

			const result = await res.json();
			if (result.responseCode !== '00') {
				throw new Error(result.responseMessage || 'Failed to create room');
			}

			myNickname = nicknameInput.trim();
			roomCode = result.data.roomCode;
			roomId = result.data.roomId;
			memberId = result.data.memberId;

			messages = [
				{
					id: 'init-sys',
					type: 'system',
					text: 'Session established. Messages will self-destruct on disconnect.',
					boxed: true
				}
			];

			roomState = 'waiting';
			connectWebSocket(myEcdhPkBase64);
		} catch (err: any) {
			errorMessage = err.message || 'Network error, please start/check the backend server.';
		} finally {
			isLoading = false;
		}
	}

	// Call POST /v1/api/room/join (deferred hook)
	async function handleJoinRoom() {
		if (!roomCodeInput.trim() || roomCodeInput.length < 7) {
			errorMessage = 'Enter a valid Room Code (XXXX-XX)';
			return;
		}
		if (!nicknameInput.trim()) {
			errorMessage = 'Nickname cannot be empty';
			return;
		}
		isLoading = true;
		errorMessage = '';
		try {
			await generateKeys();
			const myEcdhPkBase64 = await exportPublicKey(myEcdhKeyPair!.publicKey);

			const res = await fetch(`${BASE_URL}/room/join`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					roomCode: roomCodeInput.trim(),
					ecdhPublicKey: myEcdhPkBase64,
					ecdsaPublicKey: myEcdsaPkBase64
				})
			});

			const result = await res.json();
			if (result.responseCode !== '00') {
				throw new Error(result.responseMessage || 'Failed to join room');
			}

			myNickname = nicknameInput.trim();
			roomCode = roomCodeInput.trim();
			roomId = result.data.roomId;
			memberId = result.data.memberId;

			// Kunci ECDH peer untuk enkripsi + kunci ECDSA peer untuk verifikasi tanda tangan
			peerEcdhPublicKey = await importPublicKey(result.data.peerPublicKey);
			peerEcdsaPublicKey = await importSigningPublicKey(result.data.peerSignKey);

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
		} catch (err: any) {
			errorMessage = err.message || 'Network error, please start/check the backend server.';
		} finally {
			isLoading = false;
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
						timestamp,
						publicKey: myEcdsaPkBase64
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
		aesKey = null;
		nicknameInput = '';
		roomCodeInput = '';
		stopTimer();
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

	// Initialize scroll and timer on load
	onMount(() => {
		const params = new URLSearchParams(window.location.search);
		const tab = params.get('tab');
		if (tab === 'create' || tab === 'join') {
			activeTab = tab;
		}
		scrollToBottom();
	});
</script>

<div class="h-screen bg-background select-none">
	{#if roomState === 'setup'}
		<!-- SETUP WELCOME SCREEN -->
		<div class="flex items-center justify-center h-full p-6">
			<div class="w-full max-w-md bg-white rounded-2xl shadow-premium border border-outline-variant/30 overflow-hidden flex flex-col p-6 space-y-6">
				<!-- Brand -->
				<div class="flex items-center space-x-3 justify-center">
					<img src="/logo.webp" alt="UMBRA Logo" class="w-12 h-12 object-contain rounded-xl" />
					<span class="font-headline-md text-2xl font-bold text-primary tracking-tight">UMBRA</span>
				</div>

				<div class="text-center space-y-1">
					<h3 class="text-on-surface font-semibold text-lg">Secure Ephemeral Session</h3>
					<p class="text-sm text-outline">End-to-End Encrypted & Zero-Persistence</p>
				</div>

				<!-- Setup Tabs -->
				<div class="flex border-b border-outline-variant/20">
					<button
						onclick={() => { activeTab = 'create'; errorMessage = ''; }}
						class="flex-1 pb-3 text-center font-medium text-sm transition-colors {activeTab === 'create' ? 'border-b-2 border-primary text-primary font-bold' : 'text-outline hover:text-on-surface'}"
					>
						Create Room
					</button>
					<button
						onclick={() => { activeTab = 'join'; errorMessage = ''; }}
						class="flex-1 pb-3 text-center font-medium text-sm transition-colors {activeTab === 'join' ? 'border-b-2 border-primary text-primary font-bold' : 'text-outline hover:text-on-surface'}"
					>
						Join Room
					</button>
				</div>

				<!-- Input Form -->
				<div class="space-y-4 pt-2">
					{#if errorMessage}
						<div class="p-3 bg-error-container text-on-error-container rounded-xl text-sm border border-error/20">
							{errorMessage}
						</div>
					{/if}

					<div class="flex flex-col space-y-2">
						<label for="nickname" class="font-label-mono text-sm text-outline uppercase tracking-widest">Nickname</label>
						<input
							type="text"
							id="nickname"
							bind:value={nicknameInput}
							placeholder="e.g. Neon_Specter"
							disabled={isLoading}
							class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-3 text-on-surface focus:outline-none focus:border-primary-container focus:ring-4 focus:ring-primary-container/10 transition-all font-body-md"
						/>
					</div>

					{#if activeTab === 'join'}
						<div class="flex flex-col space-y-2">
							<label for="room-code-input" class="font-label-mono text-sm text-outline uppercase tracking-widest">Room Code</label>
							<input
								type="text"
								id="room-code-input"
								placeholder="XXXX-XX"
								disabled={isLoading}
								oninput={handleRoomCodeInput}
								class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-3 text-primary font-label-mono font-bold tracking-widest focus:outline-none focus:border-primary-container focus:ring-4 focus:ring-primary-container/10 transition-all text-center"
							/>
						</div>
					{/if}

					<div class="pt-2">
						{#if activeTab === 'create'}
							<button
								onclick={handleCreateRoom}
								disabled={isLoading}
								class="w-full py-4 bg-primary-container text-white font-bold rounded-xl shadow-premium hover:bg-primary transition-all active:scale-[0.98] flex items-center justify-center space-x-2"
							>
								{#if isLoading}
									<div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
									<span>Generating Room...</span>
								{:else}
									<span class="material-symbols-outlined text-[20px]">add_circle</span>
									<span>Create Room</span>
								{/if}
							</button>
						{:else}
							<button
								onclick={handleJoinRoom}
								disabled={isLoading}
								class="w-full py-4 bg-primary-container text-white font-bold rounded-xl shadow-premium hover:bg-primary transition-all active:scale-[0.98] flex items-center justify-center space-x-2"
							>
								{#if isLoading}
									<div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
									<span>Connecting...</span>
								{:else}
									<span class="material-symbols-outlined text-[20px]">login</span>
									<span>Join Room</span>
								{/if}
							</button>
						{/if}
					</div>
				</div>
			</div>
		</div>
	{:else}
		<!-- CHAT ROOM INTERFACE -->
		<div class="flex h-screen overflow-hidden w-full bg-surface-container-lowest">
			<!-- Sidebar (Left, 300px) -->
			<aside class="w-[300px] flex-shrink-0 bg-surface-container-low border-r border-outline-variant/20 flex flex-col justify-between sidebar-transition">
				<div class="p-gutter space-y-8">
					<!-- Brand Anchor -->
					<div class="flex items-center space-x-3">
						<img src="/logo.webp" alt="UMBRA Logo" class="w-10 h-10 object-contain rounded-xl" />
						<span class="font-bold text-2xl text-primary tracking-tight font-['Space_Grotesk']">UMBRA</span>
					</div>
					
					<!-- Session Identity -->
					<div class="space-y-stack-md">
						<div class="flex flex-col space-y-2">
							<span class="font-label-mono text-sm text-outline uppercase tracking-widest text-[12px]">Active Alias</span>
							<div class="flex items-center space-x-3 bg-surface-container-highest p-3 rounded-xl border border-outline-variant/30">
								<div class="w-8 h-8 rounded-full bg-secondary-fixed flex items-center justify-center text-on-secondary-fixed font-semibold">
									<span class="material-symbols-outlined text-[18px]">person</span>
								</div>
								<span class="font-medium text-on-surface">{myNickname}</span>
							</div>
						</div>
						
						<div class="flex flex-col space-y-2">
							<span class="font-label-mono text-sm text-outline uppercase tracking-widest text-[12px]">Room Code</span>
							<!-- svelte-ignore a11y_click_events_have_key_events -->
							<!-- svelte-ignore a11y_no_static_element_interactions -->
							<div class="group relative flex items-center justify-between bg-white p-3 rounded-xl border border-outline-variant/30 cursor-pointer hover:border-primary transition-colors" onclick={copyCode}>
								<span class="font-label-mono text-primary font-bold tracking-tighter" id="room-code">{roomCode}</span>
								<span class="material-symbols-outlined text-outline group-hover:text-primary transition-colors text-[20px]">content_copy</span>
							</div>
						</div>
					</div>
					
					<!-- Status Badges -->
					<div class="space-y-stack-sm pt-4 border-t border-outline-variant/20">
						<div class="flex items-center space-x-3 px-1">
							<div class="relative flex h-2 w-2">
								<span class="pulse-green absolute inline-flex h-full w-full rounded-full bg-tertiary-container opacity-75"></span>
								<span class="relative inline-flex rounded-full h-2 w-2 bg-tertiary-container"></span>
							</div>
							<span class="font-label-mono text-[12px] text-tertiary font-bold">
								{#if aesKey}E2EE ACTIVE{:else}E2EE WAITING{/if}
							</span>
						</div>
						<div class="flex items-center space-x-3 px-1">
							<span class="material-symbols-outlined {isConnected ? 'text-primary' : 'text-outline'} text-[18px]">wifi_tethering</span>
							<span class="font-label-mono text-[12px] text-on-surface-variant uppercase">
								WS: {isConnected ? 'CONNECTED' : 'DISCONNECTED'}
							</span>
						</div>
						<div class="flex items-center space-x-3 px-1">
							<span class="material-symbols-outlined text-outline text-[18px]">timer</span>
							<span class="font-label-mono text-[12px] text-on-surface-variant">TTL: {formatTime(timeLeft)}</span>
						</div>
					</div>
				</div>
				
				<!-- Footer Action -->
				<div class="p-gutter">
					<button onclick={leaveRoom} class="w-full flex items-center justify-center space-x-2 py-4 bg-[#EF4444]/10 text-[#EF4444] rounded-xl font-bold hover:bg-[#EF4444] hover:text-white transition-all duration-200 active:scale-[0.98]">
						<span class="material-symbols-outlined">logout</span>
						<span>Leave Room</span>
					</button>
				</div>
			</aside>

			<!-- Main Area: Chat Content -->
			<main class="flex-grow flex flex-col relative bg-white">
				<!-- Chat Header (Contextual) -->
				<header class="h-16 flex items-center justify-between px-gutter border-b border-outline-variant/10 z-10 bg-white/80 backdrop-blur-md">
					<div class="flex items-center space-x-4">
						<div class="flex -space-x-2">
							<div class="w-8 h-8 rounded-full border-2 border-white bg-primary-fixed flex items-center justify-center text-[10px] font-bold text-on-primary-fixed-variant select-none">
								{myInitials}
							</div>
							{#if aesKey}
								<div class="w-8 h-8 rounded-full border-2 border-white bg-tertiary-fixed flex items-center justify-center text-[10px] font-bold text-on-tertiary-fixed-variant select-none">
									JS
								</div>
							{/if}
						</div>
						<h2 class="font-bold text-lg text-on-surface font-['Space_Grotesk']">
							Anonymous Session <span class="text-outline-variant ml-2 font-normal font-sans">#421</span>
						</h2>
					</div>
					
					<div class="flex items-center space-x-2 relative">
						<!-- Toggle Search Input or Search Button (Right next to dropdown) -->
						{#if showSearchInput}
							<div class="flex items-center bg-surface-container-low px-3 py-1.5 rounded-lg border border-outline-variant/30 max-w-xs transition-all animate-in fade-in slide-in-from-right-2 duration-150 mr-1">
								<span class="material-symbols-outlined text-[16px] text-outline mr-2">search</span>
								<input
									type="text"
									bind:value={searchQuery}
									placeholder="Search message..."
									class="bg-transparent border-none text-xs focus:ring-0 p-0 text-on-surface outline-none w-36 font-body-md"
								/>
								<button onclick={() => { showSearchInput = false; searchQuery = ''; }} class="p-1 hover:text-red-500 transition-colors ml-1 flex items-center">
									<span class="material-symbols-outlined text-[16px]">close</span>
								</button>
							</div>
						{:else}
							<button 
								onclick={() => showSearchInput = true} 
								class="p-2 text-on-surface-variant hover:bg-surface-container rounded-lg transition-colors"
							>
								<span class="material-symbols-outlined">search</span>
							</button>
						{/if}
						
						<!-- More Options Menu with Dropdown -->
						<div>
							<button 
								onclick={() => showMoreMenu = !showMoreMenu} 
								class="p-2 text-on-surface-variant hover:bg-surface-container rounded-lg transition-colors {showMoreMenu ? 'bg-surface-container text-primary' : ''}"
							>
								<span class="material-symbols-outlined">more_vert</span>
							</button>
							{#if showMoreMenu}
								<!-- svelte-ignore a11y_click_events_have_key_events -->
								<!-- svelte-ignore a11y_no_static_element_interactions -->
								<div class="absolute right-0 mt-2 w-48 bg-white border border-outline-variant/30 rounded-xl shadow-premium py-2 z-50 animate-in fade-in slide-in-from-top-2 duration-150" onclick={() => showMoreMenu = false}>
									<button onclick={clearChat} class="w-full text-left px-4 py-2 text-sm text-on-surface hover:bg-surface-container transition-colors flex items-center space-x-2">
										<span class="material-symbols-outlined text-[18px]">delete_sweep</span>
										<span>Clear Chat</span>
									</button>
									<button onclick={() => { toastMessage = 'Notifications muted'; showToast = true; setTimeout(() => showToast = false, 2000); }} class="w-full text-left px-4 py-2 text-sm text-on-surface hover:bg-surface-container transition-colors flex items-center space-x-2">
										<span class="material-symbols-outlined text-[18px]">volume_off</span>
										<span>Mute Session</span>
									</button>
								</div>
							{/if}
						</div>
					</div>
				</header>

				<!-- ACTIVE CHAT AREA (Filtered by search) -->
				<div bind:this={chatContainer} class="flex-grow overflow-y-auto p-gutter space-y-6 chat-scroll" id="chat-container">
					{#if filteredMessages.length === 0}
						<div class="flex flex-col items-center justify-center h-full text-outline py-8 text-center space-y-2">
							<span class="material-symbols-outlined text-[36px]">search_off</span>
							<p class="text-sm">No messages match your search query.</p>
						</div>
					{:else}
						{#each filteredMessages as msg (msg.id)}
							{#if msg.type === 'system'}
								<div class="flex justify-center">
									{#if msg.boxed}
										<p class="italic text-on-surface-variant/60 text-sm bg-surface-container-low px-4 py-1 rounded-full border border-outline-variant/10 text-center">
											{msg.text}
										</p>
									{:else}
										<p class="italic text-on-surface-variant/60 text-sm text-center font-semibold">
											{msg.text}
										</p>
									{/if}
								</div>
							{:else if msg.type === 'partner'}
								<div class="flex items-end space-x-2 max-w-[80%]">
									<div class="w-8 h-8 rounded-full bg-secondary-container flex-shrink-0 flex items-center justify-center text-[12px] font-bold text-on-secondary-container">
										{msg.sender}
									</div>
									<div class="bg-[#F5F8FA] text-on-surface p-4 rounded-2xl shadow-sm bubble-partner">
										<p class="text-base whitespace-pre-wrap select-text">{msg.text}</p>
										<span class="block mt-2 text-[10px] text-outline text-right">{msg.time}</span>
									</div>
								</div>
							{:else if msg.type === 'self'}
								<div class="flex items-end justify-end space-x-2 ml-auto max-w-[80%]">
									<div class="bg-primary-container text-white p-4 rounded-2xl shadow-premium bubble-self">
										<p class="text-base whitespace-pre-wrap select-text">{msg.text}</p>
										<span class="block mt-2 text-[10px] opacity-70 text-right">{msg.time}</span>
									</div>
									<div class="w-8 h-8 rounded-full bg-primary flex-shrink-0 flex items-center justify-center text-[12px] font-bold text-white font-label-mono">
										{msg.sender}
									</div>
								</div>
							{/if}
						{/each}
					{/if}
				</div>

				<!-- Sticky Input Area -->
				<div class="p-gutter pt-2 bg-white">
					<div class="max-w-[1000px] mx-auto bg-white border border-outline-variant/30 rounded-2xl shadow-sm focus-within:ring-4 focus-within:ring-primary-container/10 focus-within:border-primary transition-all">
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
							<span class="text-[10px] font-label-mono text-outline-variant uppercase">Ephemeral Encryption Active</span>
						</div>
						<div class="flex items-end p-3 space-x-3">
							<textarea
								bind:this={textareaElement}
								bind:value={messageInput}
								onkeydown={handleKeyDown}
								class="flex-grow bg-transparent border-none focus:ring-0 resize-none max-h-32 text-base py-2 px-1 text-on-surface placeholder-outline-variant/60 outline-none"
								id="message-input"
								placeholder="Ketik pesan..."
								rows="1"
							></textarea>
							<button
								onclick={sendMessage}
								class="w-12 h-12 bg-primary-container rounded-xl flex items-center justify-center text-white shadow-premium hover:bg-primary transition-all active:scale-[0.98] flex-shrink-0 group"
								id="send-btn"
							>
								<span class="material-symbols-outlined group-hover:translate-x-0.5 transition-transform">send</span>
							</button>
						</div>
					</div>
					<div class="flex justify-center py-4">
						<p class="text-[11px] font-label-mono text-outline-variant uppercase tracking-widest flex items-center">
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
	<div class="fixed bottom-24 left-1/2 -translate-x-1/2 bg-inverse-surface text-white px-6 py-3 rounded-full text-sm font-medium shadow-lg z-50 animate-bounce">
		{toastMessage}
	</div>
{/if}
