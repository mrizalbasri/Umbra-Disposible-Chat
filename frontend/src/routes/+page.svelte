<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount, tick } from 'svelte';

	let activeNav = $state('features');
	let hoveredNav = $state<string | null>(null);
	let currentLang = $state<'id' | 'en'>('id');

	function toggleLang() {
		currentLang = currentLang === 'id' ? 'en' : 'id';
	}

	const translations = {
		id: {
			features: 'Fitur',
			security: 'Keamanan',
			docs: 'Dokumentasi',
			start: 'Mulai',
			badge: 'END-TO-END ENCRYPTED',
			h1_1: 'Bicara Bebas,',
			h1_2: 'Tanpa Jejak,',
			h1_3: 'Tanpa Identitas.',
			sub: 'Platform komunikasi paling aman untuk percakapan rahasia. Tidak ada login, tidak ada pelacakan, dan data Anda hilang seketika setelah selesai.',
			btnStart: 'Mulai Chat Sekarang →',
			btnProtocol: 'Pelajari Protokol'
		},
		en: {
			features: 'Features',
			security: 'Security',
			docs: 'Documentation',
			start: 'Start',
			badge: 'END-TO-END ENCRYPTED',
			h1_1: 'Speak Freely,',
			h1_2: 'Without Traces,',
			h1_3: 'Without Identity.',
			sub: 'The most secure communication platform for confidential conversations. No login, no tracking, and your data vanishes immediately when done.',
			btnStart: 'Start Chat Now →',
			btnProtocol: 'Learn Protocol'
		}
	};

	let t = $derived(translations[currentLang]);

	let displayNav = $derived(hoveredNav || activeNav);

	let indicatorLeft = $state(0);
	let indicatorWidth = $state(0);
	let featuresBtn = $state<HTMLButtonElement | null>(null);
	let securityBtn = $state<HTMLButtonElement | null>(null);
	let docsBtn = $state<HTMLButtonElement | null>(null);

	async function updateIndicator(navId: string) {
		await tick();
		let btn: HTMLButtonElement | null = null;
		if (navId === 'features') btn = featuresBtn;
		else if (navId === 'security') btn = securityBtn;
		else if (navId === 'docs') btn = docsBtn;

		if (btn && btn.offsetWidth > 0) {
			indicatorLeft = btn.offsetLeft;
			indicatorWidth = btn.offsetWidth;
		}
	}

	$effect(() => {
		updateIndicator(displayNav);
	});

	function scrollToSecurity() {
		document.querySelector('#security')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
	}

	function scrollTo(id: string) {
		document.querySelector(`#${id}`)?.scrollIntoView({ behavior: 'smooth', block: 'start' });
	}

	let visibleMsgs = $state(0);
	let isTyping = $state(true);
	let typingUser = $state('Anonymous User #294');

	let heroRevealed = $state(false);
	let ctaRevealed = $state(false);
	let stepsRevealed = $state(false);
	let featuresRevealed = $state(false);
	let footerRevealed = $state(false);

	onMount(() => {
		updateIndicator(activeNav);

		// Immediately reveal hero on mount
		setTimeout(() => (heroRevealed = true), 100);

		// Timeline: Typing indicator BEFORE each full pop-up bubble
		const t1 = setTimeout(() => {
			isTyping = false;
			visibleMsgs = 1; // 1.0s: Full Chat 1 pops up!
		}, 1000);

		const t2 = setTimeout(() => {
			typingUser = 'You';
			isTyping = true; // 1.8s: You are typing...
		}, 1800);

		const t3 = setTimeout(() => {
			isTyping = false;
			visibleMsgs = 2; // 2.8s: Full Chat 2 (blue) pops up!
		}, 2800);

		const t4 = setTimeout(() => {
			typingUser = 'Anonymous User #294';
			isTyping = true; // 3.6s: User 294 is typing...
		}, 3600);

		const t5 = setTimeout(() => {
			isTyping = false;
			visibleMsgs = 3; // 4.6s: Full Chat 3 pops up!
		}, 4600);

		const handleResize = () => updateIndicator(activeNav);
		window.addEventListener('resize', handleResize, { passive: true });

		const sectionIds = ['features', 'security', 'docs'];
		const sectionsMap = [
			{ sel: '.cta-section', set: () => (ctaRevealed = true) },
			{ sel: '.steps', set: () => (stepsRevealed = true) },
			{ sel: '.features', set: () => (featuresRevealed = true) },
			{ sel: 'footer', set: () => (footerRevealed = true) }
		];

		const handleScroll = () => {
			const scrollPosition = window.scrollY + 120;
			for (let i = sectionIds.length - 1; i >= 0; i--) {
				const id = sectionIds[i];
				const el = document.querySelector(`#${id}`);
				if (el) {
					const top = (el as HTMLElement).offsetTop;
					if (scrollPosition >= top) {
						if (activeNav !== id) {
							activeNav = id;
						}
						break;
					}
				}
			}

			// Smooth Scroll Reveal Trigger on Scroll
			const triggerBottom = window.innerHeight * 0.88;
			sectionsMap.forEach((s) => {
				const el = document.querySelector(s.sel);
				if (el) {
					const rect = el.getBoundingClientRect();
					if (rect.top < triggerBottom) {
						s.set();
					}
				}
			});
		};

		window.addEventListener('scroll', handleScroll, { passive: true });
		handleScroll(); // Initial check for active section and visible elements

		return () => {
			clearTimeout(t1);
			clearTimeout(t2);
			clearTimeout(t3);
			clearTimeout(t4);
			clearTimeout(t5);
			window.removeEventListener('scroll', handleScroll);
			window.removeEventListener('resize', handleResize);
		};
	});
</script>

<svelte:head>
	<title>UMBRA — Secure Anonymous Chat</title>
	<link
		href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&family=Inter:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600&display=swap"
		rel="stylesheet"
	/>
</svelte:head>

<main>
	<!-- NAVBAR -->
	<nav>
		<div class="nav-left">
			<img src="/logo.webp" alt="UMBRA Logo" class="logo-img" />
			<span class="logo-text">UMBRA</span>
		</div>
		<div class="nav-center" onmouseleave={() => (hoveredNav = null)}>
			<button
				bind:this={featuresBtn}
				onclick={() => scrollTo('features')}
				onmouseenter={() => (hoveredNav = 'features')}
				class="nav-link {displayNav === 'features' ? 'active' : ''}">{t.features}</button
			>
			<button
				bind:this={securityBtn}
				onclick={() => scrollTo('security')}
				onmouseenter={() => (hoveredNav = 'security')}
				class="nav-link {displayNav === 'security' ? 'active' : ''}">{t.security}</button
			>
			<button
				bind:this={docsBtn}
				onclick={() => scrollTo('docs')}
				onmouseenter={() => (hoveredNav = 'docs')}
				class="nav-link {displayNav === 'docs' ? 'active' : ''}">{t.docs}</button
			>
			<div class="sliding-indicator" style="left: {indicatorLeft}px; width: {indicatorWidth}px;"></div>
		</div>
		<div class="nav-right">
			<button class="lang-toggle" onclick={toggleLang} title="Ubah Bahasa / Switch Language">
				<span class:active-lang={currentLang === 'id'}>ID</span> / <span class:active-lang={currentLang === 'en'}>EN</span>
			</button>
			<button class="btn-mulai" onclick={() => goto(resolve('/create'))}>{t.start}</button>
		</div>
	</nav>

	<!-- HERO -->
	<section class="hero {heroRevealed ? 'revealed' : ''}">
		<div class="hero-glow"></div>
		<div class="hero-content">
			<div class="hero-left">
				<div class="badge">
					<span class="dot-green"></span>
					{t.badge}
				</div>
				<h1>
					{t.h1_1}<br />
					<span class="accent">{t.h1_2}</span><br />
					{t.h1_3}
				</h1>
				<p class="hero-sub">
					{t.sub}
				</p>
				<div class="hero-btns">
					<button class="btn-primary" onclick={() => goto(resolve('/create'))}
						>{t.btnStart}</button
					>
					<button class="btn-ghost" onclick={() => scrollTo('docs')}>{t.btnProtocol}</button>
				</div>
			</div>

			<div class="hero-right">
				<div class="chat-card">
					<div class="chat-header">
						<div class="chat-header-left">
							<div class="chat-avatar">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="19"
									height="20"
									viewBox="0 0 19 20"
									fill="none"
								>
									<path
										d="M8.995 2.45C10.7617 2.45 12.4283 2.82917 13.995 3.5875C15.5617 4.34583 16.87 5.44167 17.92 6.875C18.0367 7.025 18.0742 7.15833 18.0325 7.275C17.9908 7.39167 17.92 7.49167 17.82 7.575C17.72 7.65833 17.6033 7.69583 17.47 7.6875C17.3367 7.67917 17.22 7.60833 17.12 7.475C16.2033 6.175 15.0242 5.17917 13.5825 4.4875C12.1408 3.79583 10.6117 3.45 8.995 3.45C7.37833 3.45 5.86167 3.79583 4.445 4.4875C3.02833 5.17917 1.85333 6.175 0.92 7.475C0.82 7.625 0.703333 7.70833 0.57 7.725C0.436667 7.74167 0.32 7.70833 0.22 7.625C0.103333 7.54167 0.0325 7.4375 0.0075 7.3125C-0.0175 7.1875 0.02 7.05833 0.12 6.925C1.15333 5.50833 2.44917 4.40833 4.0075 3.625C5.56583 2.84167 7.22833 2.45 8.995 2.45ZM8.995 4.8C11.245 4.8 13.1783 5.55 14.795 7.05C16.4117 8.55 17.22 10.4083 17.22 12.625C17.22 13.4583 16.9242 14.1542 16.3325 14.7125C15.7408 15.2708 15.02 15.55 14.17 15.55C13.32 15.55 12.5908 15.2708 11.9825 14.7125C11.3742 14.1542 11.07 13.4583 11.07 12.625C11.07 12.075 10.8658 11.6125 10.4575 11.2375C10.0492 10.8625 9.56167 10.675 8.995 10.675C8.42833 10.675 7.94083 10.8625 7.5325 11.2375C7.12417 11.6125 6.92 12.075 6.92 12.625C6.92 14.2417 7.39917 15.5917 8.3575 16.675C9.31583 17.7583 10.5533 18.5167 12.07 18.95C12.22 19 12.32 19.0833 12.37 19.2C12.42 19.3167 12.4283 19.4417 12.395 19.575C12.3617 19.6917 12.295 19.7917 12.195 19.875C12.095 19.9583 11.97 19.9833 11.82 19.95C10.0867 19.5167 8.67 18.6542 7.57 17.3625C6.47 16.0708 5.92 14.4917 5.92 12.625C5.92 11.7917 6.22 11.0917 6.82 10.525C7.42 9.95833 8.145 9.675 8.995 9.675C9.845 9.675 10.57 9.95833 11.17 10.525C11.77 11.0917 12.07 11.7917 12.07 12.625C12.07 13.175 12.2783 13.6375 12.695 14.0125C13.1117 14.3875 13.6033 14.575 14.17 14.575C14.7367 14.575 15.22 14.3875 15.62 14.0125C16.02 13.6375 16.22 13.175 16.22 12.625C16.22 10.6917 15.5117 9.06667 14.095 7.75C12.6783 6.43333 10.9867 5.775 9.02 5.775C7.05333 5.775 5.36167 6.43333 3.945 7.75C2.52833 9.06667 1.82 10.6833 1.82 12.6C1.82 13 1.8575 13.5 1.9325 14.1C2.0075 14.7 2.18667 15.4 2.47 16.2C2.52 16.35 2.51583 16.4833 2.4575 16.6C2.39917 16.7167 2.30333 16.8 2.17 16.85C2.03667 16.9 1.9075 16.8958 1.7825 16.8375C1.6575 16.7792 1.57 16.6833 1.52 16.55C1.27 15.9 1.09083 15.2542 0.9825 14.6125C0.874167 13.9708 0.82 13.3083 0.82 12.625C0.82 10.4083 1.62417 8.55 3.2325 7.05C4.84083 5.55 6.76167 4.8 8.995 4.8ZM8.995 0C10.0617 0 11.1033 0.129167 12.12 0.3875C13.1367 0.645833 14.12 1.01667 15.07 1.5C15.22 1.58333 15.3075 1.68333 15.3325 1.8C15.3575 1.91667 15.345 2.03333 15.295 2.15C15.245 2.26667 15.1617 2.35833 15.045 2.425C14.9283 2.49167 14.7867 2.48333 14.62 2.4C13.7367 1.95 12.8242 1.60417 11.8825 1.3625C10.9408 1.12083 9.97833 1 8.995 1C8.02833 1 7.07833 1.1125 6.145 1.3375C5.21167 1.5625 4.32 1.91667 3.47 2.4C3.33667 2.48333 3.20333 2.50417 3.07 2.4625C2.93667 2.42083 2.83667 2.33333 2.77 2.2C2.70333 2.06667 2.68667 1.94583 2.72 1.8375C2.75333 1.72917 2.83667 1.63333 2.97 1.55C3.90333 1.05 4.87833 0.666667 5.895 0.4C6.91167 0.133333 7.945 0 8.995 0ZM8.995 7.225C10.545 7.225 11.8783 7.74583 12.995 8.7875C14.1117 9.82917 14.67 11.1083 14.67 12.625C14.67 12.775 14.6242 12.8958 14.5325 12.9875C14.4408 13.0792 14.32 13.125 14.17 13.125C14.0367 13.125 13.92 13.0792 13.82 12.9875C13.72 12.8958 13.67 12.775 13.67 12.625C13.67 11.375 13.2075 10.3292 12.2825 9.4875C11.3575 8.64583 10.2617 8.225 8.995 8.225C7.72833 8.225 6.64083 8.64583 5.7325 9.4875C4.82417 10.3292 4.37 11.375 4.37 12.625C4.37 13.975 4.60333 15.1208 5.07 16.0625C5.53667 17.0042 6.22 17.95 7.12 18.9C7.22 19 7.27 19.1167 7.27 19.25C7.27 19.3833 7.22 19.5 7.12 19.6C7.02 19.7 6.90333 19.75 6.77 19.75C6.63667 19.75 6.52 19.7 6.42 19.6C5.43667 18.5667 4.6825 17.5125 4.1575 16.4375C3.6325 15.3625 3.37 14.0917 3.37 12.625C3.37 11.1083 3.92 9.82917 5.02 8.7875C6.12 7.74583 7.445 7.225 8.995 7.225ZM8.97 12.125C9.12 12.125 9.24083 12.175 9.3325 12.275C9.42417 12.375 9.47 12.4917 9.47 12.625C9.47 13.875 9.92 14.9 10.82 15.7C11.72 16.5 12.77 16.9 13.97 16.9C14.07 16.9 14.2117 16.8917 14.395 16.875C14.5783 16.8583 14.77 16.8333 14.97 16.8C15.12 16.7667 15.2492 16.7875 15.3575 16.8625C15.4658 16.9375 15.5367 17.05 15.57 17.2C15.6033 17.3333 15.5783 17.45 15.495 17.55C15.4117 17.65 15.3033 17.7167 15.17 17.75C14.87 17.8333 14.6075 17.8792 14.3825 17.8875C14.1575 17.8958 14.02 17.9 13.97 17.9C12.4867 17.9 11.1992 17.4 10.1075 16.4C9.01583 15.4 8.47 14.1417 8.47 12.625C8.47 12.4917 8.51583 12.375 8.6075 12.275C8.69917 12.175 8.82 12.125 8.97 12.125Z"
										fill="#4B6171"
									/>
								</svg>
							</div>
							<div>
								<div class="chat-name">Anonymous User #294</div>
								<div class="chat-status"><span class="dot-status"></span> Secure Tunnel Active</div>
							</div>
						</div>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="16"
							height="20"
							viewBox="0 0 16 20"
							fill="none"
						>
							<path
								d="M6.95 13.55L12.6 7.9L11.175 6.475L6.95 10.7L4.85 8.6L3.425 10.025L6.95 13.55ZM8 20C5.68333 19.4167 3.77083 18.0875 2.2625 16.0125C0.754167 13.9375 0 11.6333 0 9.1V3L8 0L16 3V9.1C16 11.6333 15.2458 13.9375 13.7375 16.0125C12.2292 18.0875 10.3167 19.4167 8 20ZM8 17.9C9.73333 17.35 11.1667 16.25 12.3 14.6C13.4333 12.95 14 11.1167 14 9.1V4.375L8 2.125L2 4.375V9.1C2 11.1167 2.56667 12.95 3.7 14.6C4.83333 16.25 6.26667 17.35 8 17.9Z"
								fill="#6E7881"
							/>
						</svg>
					</div>
					<div class="chat-body">
						{#if visibleMsgs >= 1}
							<div class="msg-them anim-pop-them">Apakah data ini benar-benar tidak tersimpan di server?</div>
						{/if}

						{#if visibleMsgs >= 2}
							<div class="msg-me anim-pop-me">
								Tepat sekali. Menggunakan enkripsi X3DH dan Double Ratchet. Segera setelah tab
								didelete.
							</div>
						{/if}

						{#if visibleMsgs >= 3}
							<div class="msg-them anim-pop-them">Luar biasa. Mulai kirim dokumen rahasianya.</div>
						{/if}

						{#if isTyping}
							<div
								class="typing-bar {typingUser === 'You' ? 'me anim-pop-me' : 'them anim-pop-them'}"
							>
								<span class="typing-text">{typingUser} is typing</span>
								<span class="typing-dot"></span>
								<span class="typing-dot"></span>
								<span class="typing-dot"></span>
							</div>
						{/if}
					</div>
					<div class="chat-input-bar">
						<div class="input-fake">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="13"
								height="20"
								viewBox="0 0 13 20"
								fill="none"
							>
								<path
									d="M12.5 13.75C12.5 15.4833 11.8917 16.9583 10.675 18.175C9.45833 19.3917 7.98333 20 6.25 20C4.51667 20 3.04167 19.3917 1.825 18.175C0.608333 16.9583 0 15.4833 0 13.75V4.5C0 3.25 0.4375 2.1875 1.3125 1.3125C2.1875 0.4375 3.25 0 4.5 0C5.75 0 6.8125 0.4375 7.6875 1.3125C8.5625 2.1875 9 3.25 9 4.5V13.25C9 14.0167 8.73333 14.6667 8.2 15.2C7.66667 15.7333 7.01667 16 6.25 16C5.48333 16 4.83333 15.7333 4.3 15.2C3.76667 14.6667 3.5 14.0167 3.5 13.25V4H5.5V13.25C5.5 13.4667 5.57083 13.6458 5.7125 13.7875C5.85417 13.9292 6.03333 14 6.25 14C6.46667 14 6.64583 13.9292 6.7875 13.7875C6.92917 13.6458 7 13.4667 7 13.25V4.5C6.98333 3.8 6.7375 3.20833 6.2625 2.725C5.7875 2.24167 5.2 2 4.5 2C3.8 2 3.20833 2.24167 2.725 2.725C2.24167 3.20833 2 3.8 2 4.5V13.75C1.98333 14.9333 2.39167 15.9375 3.225 16.7625C4.05833 17.5875 5.06667 18 6.25 18C7.41667 18 8.40833 17.5875 9.225 16.7625C10.0417 15.9375 10.4667 14.9333 10.5 13.75V4H12.5V13.75Z"
									fill="#6E7881"
								/>
							</svg>
							<span>Type secure message...</span>
						</div>
						<button class="send-btn" type="button" aria-label="Kirim pesan">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="19"
								height="16"
								viewBox="0 0 19 16"
								fill="none"
							>
								<path
									d="M0 16V0L19 8L0 16ZM2 13L13.85 8L2 3V6.5L8 8L2 9.5V13ZM2 13V8V3V6.5V9.5V13Z"
									fill="#00658D"
								/>
							</svg>
						</button>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- SIAP MEMULAI (FEATURES) -->
	<section class="cta-section {ctaRevealed ? 'revealed' : ''}" id="features">
		<h2>Siap Memulai?</h2>
		<p class="section-sub">Pilih bagaimana Anda ingin terhubung dengan orang lain secara anonim.</p>
		<div class="cta-cards">
			<div class="cta-card featured">
				<div class="cta-icon">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="36"
						height="37"
						viewBox="0 0 36 37"
						fill="none"
					>
						<path
							d="M6.34138e-06 35.9926L0.00734342 31.9926L4.00734 32L4.06603 1.28588e-05L24.066 0.0366982L24.0623 2.03669L32.0623 2.05137L32.0073 32.0513L36.0073 32.0587L35.9999 36.0586L28 36.044L28.055 6.04403L24.055 6.03669L24 36.0366L6.34138e-06 35.9926ZM8.05869 4.00734L8.00733 32.0073L8.05869 4.00734ZM16.0293 20.022C16.596 20.023 17.0713 19.8322 17.4554 19.4496C17.8394 19.067 18.032 18.5923 18.033 18.0257C18.034 17.459 17.8432 16.9836 17.4606 16.5996C17.078 16.2156 16.6033 16.023 16.0367 16.022C15.47 16.021 14.9946 16.2118 14.6106 16.5944C14.2266 16.977 14.034 17.4517 14.033 18.0183C14.032 18.585 14.2228 19.0603 14.6054 19.4444C14.988 19.8284 15.4627 20.021 16.0293 20.022ZM8.00733 32.0073L20.0073 32.0293L20.0587 4.02935L8.05869 4.00734L8.00733 32.0073Z"
							fill="#00658D"
						/>
					</svg>
				</div>
				<h3>Private Room</h3>
				<p>Buat ruang obrolan eksklusif dan bagikan kode rahasia kepada partner bicara Anda.</p>
				<button class="btn-teal full" onclick={() => goto(resolve('/create'))}>Buat Ruangan</button>
			</div>
			<div class="cta-card">
				<div class="cta-icon">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="33"
						height="33"
						viewBox="0 0 33 33"
						fill="none"
					>
						<path
							d="M19.9949 32.0315L20.0022 28.0315L25.2022 28.041L18.8638 21.6794L21.7191 18.8346L28.0075 25.1462L28.0169 20.0462L32.0169 20.0535L31.9948 32.0535L19.9949 32.0315ZM2.79489 32L2.81999e-05 29.1948L25.2462 4.04109L20.0462 4.03155L20.0536 0.0315584L32.0535 0.0535697L32.0315 12.0535L28.0315 12.0462L28.0411 6.84622L2.79489 32ZM10.3795 13.1638L0.0484529 2.79487L2.85358 9.00105e-06L13.1846 10.369L10.3795 13.1638Z"
							fill="#6E7881"
						/>
					</svg>
				</div>
				<h3>Random Match</h3>
				<p>Terhubung secara acak dengan orang lain di seluruh dunia tanpa mengenal identitas.</p>
				<button class="btn-outline-card full" onclick={() => goto(resolve('/match'))}
					>Cari Teman Chat</button
				>
			</div>
		</div>
	</section>

	<!-- 3 LANGKAH -->
	<section class="steps {stepsRevealed ? 'revealed' : ''}" id="security">
		<span class="workflow-label">WORKFLOW</span>
		<h2>3 Langkah Sederhana</h2>
		<div class="steps-row">
			<div class="step">
				<div class="step-circle">01</div>
				<h3>Buka UMBRA</h3>
				<p>Akses web UMBRA dari perangkat mana saja. Tidak ada aplikasi yang perlu diunduh.</p>
			</div>
			<div class="step">
				<div class="step-circle">02</div>
				<h3>Pilih Mode Chat</h3>
				<p>Buat Private Room dengan kode unik atau cari Random Match secara instan.</p>
			</div>
			<div class="step">
				<div class="step-circle">03</div>
				<h3>Chat Aman</h3>
				<p>Mulai berkomunikasi dengan tenang. Semua pesan dienkripsi secara real-time.</p>
			</div>
		</div>
	</section>

	<!-- PRIVASI TANPA KOMPROMI (DOCUMENTATION) -->
	<section class="features {featuresRevealed ? 'revealed' : ''}" id="docs">
		<h2>Privasi Tanpa Kompromi</h2>
		<p class="section-sub">
			UMBRA dirancang oleh pakar keamanan untuk memberikan perlindungan mutlak bagi identitas dan
			pesan Anda.
		</p>
		<div class="feature-grid">
			<div class="feature-card">
				<div class="feature-icon-box">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="16"
						height="21"
						viewBox="0 0 16 21"
						fill="none"
					>
						<path
							d="M2 21C1.45 21 0.979167 20.8042 0.5875 20.4125C0.195833 20.0208 0 19.55 0 19V9C0 8.45 0.195833 7.97917 0.5875 7.5875C0.979167 7.19583 1.45 7 2 7H3V5C3 3.61667 3.4875 2.4375 4.4625 1.4625C5.4375 0.4875 6.61667 0 8 0C9.38333 0 10.5625 0.4875 11.5375 1.4625C12.5125 2.4375 13 3.61667 13 5V7H14C14.55 7 15.0208 7.19583 15.4125 7.5875C15.8042 7.97917 16 8.45 16 9V19C16 19.55 15.8042 20.0208 15.4125 20.4125C15.0208 20.8042 14.55 21 14 21H2ZM2 19H14V9H2V19ZM8 16C8.55 16 9.02083 15.8042 9.4125 15.4125C9.80417 15.0208 10 14.55 10 14C10 13.45 9.80417 12.9792 9.4125 12.5875C9.02083 12.1958 8.55 12 8 12C7.45 12 6.97917 12.1958 6.5875 12.5875C6.19583 12.9792 6 13.45 6 14C6 14.55 6.19583 15.0208 6.5875 15.4125C6.97917 15.8042 7.45 16 8 16ZM5 7H11V5C11 4.16667 10.7083 3.45833 10.125 2.875C9.54167 2.29167 8.83333 2 8 2C7.16667 2 6.45833 2.29167 5.875 2.875C5.29167 3.45833 5 4.16667 5 5V7ZM2 19V9V19Z"
							fill="#00658D"
						/>
					</svg>
				</div>
				<h3>Enkripsi End-to-End</h3>
				<p>
					Hanya Anda dan penerima yang dapat membaca pesan. Tidak ada pihak ketiga yang bisa
					mencegat.
				</p>
			</div>
			<div class="feature-card">
				<div class="feature-icon-box">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
					>
						<path
							d="M18.4 19.825L15.775 17.2H2.625V14.4C2.625 13.8333 2.77083 13.3125 3.0625 12.8375C3.35417 12.3625 3.74167 12 4.225 11.75C4.975 11.3667 5.7375 11.0583 6.5125 10.825C7.2875 10.5917 8.075 10.4167 8.875 10.3L0 1.425L1.425 0L19.825 18.4L18.4 19.825ZM4.625 15.2H13.775L10.775 12.2C10.7417 12.2 10.7167 12.2 10.7 12.2C10.6833 12.2 10.6583 12.2 10.625 12.2C9.69167 12.2 8.76667 12.3125 7.85 12.5375C6.93333 12.7625 6.025 13.1 5.125 13.55C4.975 13.6333 4.85417 13.75 4.7625 13.9C4.67083 14.05 4.625 14.2167 4.625 14.4V15.2ZM17.025 11.75C17.5083 11.9833 17.8917 12.3375 18.175 12.8125C18.4583 13.2875 18.6083 13.8 18.625 14.35L15.275 11C15.575 11.1167 15.8708 11.2333 16.1625 11.35C16.4542 11.4667 16.7417 11.6 17.025 11.75ZM12.825 8.55L11.35 7.075C11.7333 6.925 12.0417 6.67917 12.275 6.3375C12.5083 5.99583 12.625 5.61667 12.625 5.2C12.625 4.65 12.4292 4.17917 12.0375 3.7875C11.6458 3.39583 11.175 3.2 10.625 3.2C10.2083 3.2 9.82917 3.31667 9.4875 3.55C9.14583 3.78333 8.9 4.09167 8.75 4.475L7.275 3C7.65833 2.43333 8.14167 1.99167 8.725 1.675C9.30833 1.35833 9.94167 1.2 10.625 1.2C11.725 1.2 12.6667 1.59167 13.45 2.375C14.2333 3.15833 14.625 4.1 14.625 5.2C14.625 5.88333 14.4667 6.51667 14.15 7.1C13.8333 7.68333 13.3917 8.16667 12.825 8.55ZM13.775 15.2H4.625C4.625 15.2 4.67083 15.2 4.7625 15.2C4.85417 15.2 4.975 15.2 5.125 15.2C5.575 15.2 6.04167 15.2 6.525 15.2C7.00833 15.2 7.57917 15.2 8.2375 15.2C8.89583 15.2 9.67083 15.2 10.5625 15.2C11.4542 15.2 12.525 15.2 13.775 15.2Z"
							fill="#00658D"
						/>
					</svg>
				</div>
				<h3>Tanpa Identitas</h3>
				<p>Tidak perlu nomor telepon, email, atau akun sosial media. Anda anonim secara default.</p>
			</div>
			<div class="feature-card">
				<div class="feature-icon-box">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
					>
						<path
							d="M13.25 14.65L8.95 10.35V4.95H10.95V9.55L14.65 13.25L13.25 14.65ZM0.8 14.025C0.583333 13.5417 0.408333 13.0417 0.275 12.525C0.141667 12.0083 0.05 11.4833 0 10.95H2.025C2.075 11.3 2.14583 11.65 2.2375 12C2.32917 12.35 2.44167 12.6917 2.575 13.025L0.8 14.025ZM0 8.95C0.05 8.41667 0.141667 7.8875 0.275 7.3625C0.408333 6.8375 0.591667 6.33333 0.825 5.85L2.575 6.85C2.44167 7.18333 2.32917 7.525 2.2375 7.875C2.14583 8.225 2.075 8.58333 2.025 8.95H0ZM4.125 18.1C3.675 17.7667 3.25833 17.4042 2.875 17.0125C2.49167 16.6208 2.13333 16.2 1.8 15.75L3.55 14.75C3.78333 15.05 4.02917 15.3292 4.2875 15.5875C4.54583 15.8458 4.825 16.0917 5.125 16.325L4.125 18.1ZM3.575 5.125L1.8 4.125C2.13333 3.675 2.49167 3.25833 2.875 2.875C3.25833 2.49167 3.675 2.13333 4.125 1.8L5.125 3.575C4.84167 3.80833 4.57083 4.05417 4.3125 4.3125C4.05417 4.57083 3.80833 4.84167 3.575 5.125ZM8.95 19.9C8.41667 19.85 7.8875 19.7583 7.3625 19.625C6.8375 19.4917 6.33333 19.3083 5.85 19.075L6.85 17.325C7.18333 17.4583 7.525 17.5708 7.875 17.6625C8.225 17.7542 8.58333 17.825 8.95 17.875V19.9ZM6.85 2.575L5.85 0.825C6.33333 0.591667 6.8375 0.408333 7.3625 0.275C7.8875 0.141667 8.41667 0.05 8.95 0V2.025C8.58333 2.075 8.225 2.14583 7.875 2.2375C7.525 2.32917 7.18333 2.44167 6.85 2.575ZM10.95 19.9V17.875C11.3167 17.825 11.675 17.7542 12.025 17.6625C12.375 17.5708 12.7167 17.4583 13.05 17.325L14.05 19.075C13.5667 19.3083 13.0625 19.4917 12.5375 19.625C12.0125 19.7583 11.4833 19.85 10.95 19.9ZM13.05 2.575C12.7167 2.44167 12.375 2.32917 12.025 2.2375C11.675 2.14583 11.3167 2.075 10.95 2.025V0C11.4833 0.05 12.0125 0.141667 12.5375 0.275C13.0625 0.408333 13.5667 0.591667 14.05 0.825L13.05 2.575ZM15.775 18.1L14.775 16.325C15.0583 16.0917 15.3292 15.8458 15.5875 15.5875C15.8458 15.3292 16.0917 15.0583 16.325 14.775L18.1 15.775C17.7667 16.225 17.4083 16.6458 17.025 17.0375C16.6417 17.4292 16.225 17.7833 15.775 18.1ZM16.325 5.125C16.0917 4.84167 15.8458 4.57083 15.5875 4.3125C15.3292 4.05417 15.0583 3.80833 14.775 3.575L15.775 1.8C16.225 2.11667 16.6417 2.46667 17.025 2.85C17.4083 3.23333 17.7583 3.65 18.075 4.1L16.325 5.125ZM17.875 8.95C17.825 8.58333 17.7542 8.225 17.6625 7.875C17.5708 7.525 17.4583 7.18333 17.325 6.85L19.075 5.825C19.2917 6.325 19.4708 6.8375 19.6125 7.3625C19.7542 7.8875 19.85 8.41667 19.9 8.95H17.875ZM19.075 14.05L17.325 13.05C17.4583 12.7167 17.5708 12.375 17.6625 12.025C17.7542 11.675 17.825 11.3167 17.875 10.95H19.9C19.85 11.4833 19.7583 12.0125 19.625 12.5375C19.4917 13.0625 19.3083 13.5667 19.075 14.05Z"
							fill="#00658D"
						/>
					</svg>
				</div>
				<h3>Ephemeral Chat</h3>
				<p>Pesan akan hancur sendiri setelah waktu yang ditentukan atau setelah sesi berakhir.</p>
			</div>
			<div class="feature-card">
				<div class="feature-icon-box">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="19"
						viewBox="0 0 18 19"
						fill="none"
					>
						<path
							d="M4.5 3C4.08333 3 3.72917 3.14583 3.4375 3.4375C3.14583 3.72917 3 4.08333 3 4.5C3 4.91667 3.14583 5.27083 3.4375 5.5625C3.72917 5.85417 4.08333 6 4.5 6C4.91667 6 5.27083 5.85417 5.5625 5.5625C5.85417 5.27083 6 4.91667 6 4.5C6 4.08333 5.85417 3.72917 5.5625 3.4375C5.27083 3.14583 4.91667 3 4.5 3ZM4.5 13C4.08333 13 3.72917 13.1458 3.4375 13.4375C3.14583 13.7292 3 14.0833 3 14.5C3 14.9167 3.14583 15.2708 3.4375 15.5625C3.72917 15.8542 4.08333 16 4.5 16C4.91667 16 5.27083 15.8542 5.5625 15.5625C5.85417 15.2708 6 14.9167 6 14.5C6 14.0833 5.85417 13.7292 5.5625 13.4375C5.27083 13.1458 4.91667 13 4.5 13ZM1 0H17C17.2833 0 17.5208 0.0958333 17.7125 0.2875C17.9042 0.479167 18 0.716667 18 1V8C18 8.28333 17.9042 8.52083 17.7125 8.7125C17.5208 8.90417 17.2833 9 17 9H1C0.716667 9 0.479167 8.90417 0.2875 8.7125C0.0958333 8.52083 0 8.28333 0 8V1C0 0.716667 0.0958333 0.479167 0.2875 0.2875C0.479167 0.0958333 0.716667 0 1 0ZM2 2V7H16V2H2ZM1 10H17C17.2833 10 17.5208 10.0958 17.7125 10.2875C17.9042 10.4792 18 10.7167 18 11V18C18 18.2833 17.9042 18.5208 17.7125 18.7125C17.5208 18.9042 17.2833 19 17 19H1C0.716667 19 0.479167 18.9042 0.2875 18.7125C0.0958333 18.5208 0 18.2833 0 18V11C0 10.7167 0.0958333 10.4792 0.2875 10.2875C0.479167 10.0958 0.716667 10 1 10ZM2 12V17H16V12H2ZM2 2V7V2ZM2 12V17V12Z"
							fill="#00658D"
						/>
					</svg>
				</div>
				<h3>Zero Knowledge</h3>
				<p>
					Server kami tidak pernah menyimpan data chat atau metadata yang bisa mengidentifikasi
					Anda.
				</p>
			</div>
		</div>
	</section>

	<!-- FOOTER -->
	<footer id="docs">
		<div class="footer-top">
			<div class="footer-brand">
				<div class="footer-logo-row">
					<img src="/logo.webp" alt="UMBRA" class="footer-logo-img" />
					<span class="footer-logo-text">UMBRA</span>
				</div>
				<p class="footer-tagline">Secure · Anonymous · Ephemeral</p>
			</div>
			<div class="footer-links">
				<a href="#audit">Security Audit</a>
				<a href="#privacy">Privacy Policy</a>
				<a href="#github">Github</a>
				<a href="#status">Status</a>
			</div>
		</div>
		<div class="footer-bottom">
			<span class="footer-copy">UMBRA Protocol © 2024. All communications are ephemeral.</span>
			<span class="footer-ws"><span class="dot-green"></span> WebSocket Connected</span>
		</div>
	</footer>
</main>

<style>
	* {
		margin: 0;
		padding: 0;
		box-sizing: border-box;
	}

	main {
		font-family: 'Inter', sans-serif;
		background: #ffffff;
		color: #0f1c2c;
		overflow-x: hidden;
	}

	/* ── NAVBAR ── */
	nav {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 48px;
		height: 80px;
		background: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(8px);
		border-bottom: 1px solid #f1f5f9;
		position: sticky;
		top: 0;
		z-index: 100;
	}

	.nav-left {
		display: flex;
		align-items: center;
		gap: 12px;
	}
	.logo-img {
		height: 40px;
		width: 40px;
		padding: 6px;
		background: #ffffff;
		border: 1px solid #e2e8f0;
		border-radius: 10px;
		object-fit: contain;
	}
	.logo-text {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 24px;
		font-weight: 700;
		color: #00658d;
		letter-spacing: 0.6px;
	}

	.nav-center {
		position: relative;
		display: flex;
		gap: 36px;
		align-items: center;
	}
	.nav-link {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 15px;
		font-weight: 600;
		color: #4b5563;
		text-decoration: none;
		background: none;
		border: none;
		cursor: pointer;
		position: relative;
		padding: 6px 0;
		transition: color 0.2s ease;
	}
	.nav-link.active {
		color: #00658d;
		font-weight: 700;
	}
	.sliding-indicator {
		position: absolute;
		bottom: -4px;
		height: 2.5px;
		background-color: #00658d;
		border-radius: 2px;
		transition: left 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), width 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
		pointer-events: none;
	}
	.nav-link:hover {
		color: #0f1c2c;
	}

	.nav-right {
		display: flex;
		align-items: center;
		gap: 24px;
	}
	.lang-toggle {
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 8px;
		padding: 6px 12px;
		font-family: 'Space Grotesk', sans-serif;
		font-size: 13px;
		font-weight: 600;
		color: #94a3b8;
		cursor: pointer;
		transition: all 0.2s ease;
		display: inline-flex;
		align-items: center;
		gap: 2px;
	}
	.lang-toggle:hover {
		border-color: #00aeef;
		color: #0f1c2c;
		background: #ffffff;
	}
	.active-lang {
		color: #00658d;
		font-weight: 700;
	}
	.btn-mulai {
		background: #00aeef;
		color: #ffffff;
		border: none;
		padding: 10px 24px;
		border-radius: 12px;
		font-weight: 600;
		font-size: 14px;
		cursor: pointer;
		font-family: 'Inter', sans-serif;
		box-shadow: 0 4px 14px rgba(0, 174, 239, 0.4);
		transition: all 0.2s ease;
	}
	.btn-mulai:hover {
		background: #009ce5;
		box-shadow: 0 6px 18px rgba(0, 174, 239, 0.55);
		transform: translateY(-1px);
	}
	.btn-mulai:active {
		transform: translateY(0);
	}

	/* ── HERO ── */
	.hero {
		position: relative;
		padding: 100px 48px;
		background: #ffffff;
		overflow: hidden;
	}

	.hero-glow {
		position: absolute;
		width: 700px;
		height: 700px;
		background: radial-gradient(circle, rgba(0, 174, 239, 0.15) 0%, rgba(255, 255, 255, 0) 70%);
		top: 50%;
		left: 55%;
		transform: translate(-20%, -50%);
		z-index: 0;
	}

	.hero-content {
		position: relative;
		z-index: 1;
		display: flex;
		align-items: center;
		justify-content: space-between;
		max-width: 1200px;
		margin: 0 auto;
		gap: 64px;
	}

	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.hero-left, .hero-right, .feature-card, .cta-card {
		animation: fadeInUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) both;
	}

	.hero-left {
		max-width: 520px;
		animation-delay: 0.1s;
	}

	.hero-right {
		flex: 1;
		display: flex;
		justify-content: flex-end;
		animation-delay: 0.25s;
	}

	.badge {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		background: #ecfdf5;
		color: #006c49;
		padding: 6px 14px;
		border-radius: 999px;
		font-family: 'JetBrains Mono', monospace;
		font-size: 11px;
		font-weight: 600;
		letter-spacing: 0.5px;
		margin-bottom: 28px;
		border: 1px solid #bbf7d0;
	}

	.dot-green {
		width: 7px;
		height: 7px;
		background: #006c49;
		border-radius: 50%;
		flex-shrink: 0;
	}

	h1 {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 64px;
		font-weight: 700;
		line-height: 1.1;
		letter-spacing: -1.28px;
		color: #0f1c2c;
		margin-bottom: 24px;
	}

	.accent {
		color: #00aeef;
	}

	.hero-sub {
		font-size: 18px;
		color: #3e4850;
		line-height: 1.7;
		margin-bottom: 40px;
	}

	.hero-btns {
		display: flex;
		gap: 16px;
	}

	/* ── MICRO-INTERACTIONS ── */
	button {
		transition: transform 0.15s ease, box-shadow 0.15s ease;
		user-select: none;
	}

	button:active {
		transform: scale(0.96) !important;
	}

	.btn-primary {
		background: #00aeef;
		color: #003e58;
		border: none;
		padding: 14px 28px;
		border-radius: 8px;
		font-weight: 700;
		font-size: 15px;
		cursor: pointer;
		font-family: 'Inter', sans-serif;
		box-shadow: 0 4px 14px rgba(0, 174, 239, 0.35);
	}

	.btn-ghost {
		background: white;
		color: #0f1c2c;
		border: 1.5px solid #cbd5e1;
		padding: 14px 28px;
		border-radius: 8px;
		font-weight: 600;
		font-size: 15px;
		cursor: pointer;
		font-family: 'Inter', sans-serif;
	}
	.btn-ghost:hover {
		border-color: #00aeef;
		color: #00658d;
		box-shadow: 0 8px 20px rgba(0, 101, 141, 0.12);
	}

	.btn-teal {
		background: linear-gradient(135deg, #0088b3 0%, #006c8b 100%);
		color: white;
		border: none;
		padding: 14px 24px;
		border-radius: 10px;
		font-weight: 700;
		font-size: 15px;
		cursor: pointer;
		font-family: 'Inter', sans-serif;
		box-shadow: 0 4px 14px rgba(0, 108, 139, 0.3);
	}

	.btn-teal:hover {
		transform: translateY(-3px);
		box-shadow: 0 8px 22px rgba(0, 108, 139, 0.45);
	}

	.btn-outline-card {
		background: white;
		color: #0f1c2c;
		border: 2px solid #cbd5e1;
		padding: 14px 24px;
		border-radius: 10px;
		font-weight: 700;
		font-size: 15px;
		cursor: pointer;
		font-family: 'Inter', sans-serif;
	}

	.btn-outline-card:hover {
		transform: translateY(-3px);
		border-color: #00aeef;
		color: #00658d;
		box-shadow: 0 8px 20px rgba(0, 174, 239, 0.2);
	}

	/* ── CHAT CARD ── */
	.hero-right {
		flex: 1;
		display: flex;
		justify-content: flex-end;
	}

	.chat-card {
		width: 400px;
		background: white;
		border-radius: 16px;
		box-shadow: 0 4px 16px 0 rgba(0, 174, 239, 0.3);
		border: 1px solid #e2e8f0;
		overflow: hidden;
	}

	.chat-header {
		padding: 16px 20px;
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 1px solid #f1f5f9;
	}

	.chat-header-left {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.chat-avatar {
		width: 38px;
		height: 38px;
		background: #e0f2fe;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.chat-name {
		font-size: 13px;
		font-weight: 700;
		color: #0f1c2c;
	}
	.chat-status {
		font-size: 11px;
		color: #10b981;
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 4px;
		margin-top: 2px;
	}

	.dot-status {
		width: 6px;
		height: 6px;
		background: #10b981;
		border-radius: 50%;
	}

	.chat-body {
		padding: 20px;
		background: #f8fafc;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.msg-them,
	.msg-me {
		padding: 12px 16px;
		border-radius: 12px;
		font-size: 13px;
		line-height: 1.5;
		max-width: 86%;
	}

	.msg-them {
		background: white;
		color: #334155;
		align-self: flex-start;
		border: 1px solid #e2e8f0;
		border-bottom-left-radius: 4px;
	}

	.msg-me {
		background: #00aeef;
		color: white;
		align-self: flex-end;
		border-bottom-right-radius: 4px;
	}

	@keyframes popUpBubbleThem {
		0% {
			opacity: 0;
			transform: scale(0.2) translateY(20px);
			transform-origin: bottom left;
		}
		75% {
			opacity: 1;
			transform: scale(1.05) translateY(-2px);
		}
		100% {
			opacity: 1;
			transform: scale(1) translateY(0);
		}
	}

	@keyframes popUpBubbleMe {
		0% {
			opacity: 0;
			transform: scale(0.2) translateY(20px);
			transform-origin: bottom right;
		}
		75% {
			opacity: 1;
			transform: scale(1.05) translateY(-2px);
		}
		100% {
			opacity: 1;
			transform: scale(1) translateY(0);
		}
	}

	@keyframes typingPulse {
		0%, 100% { opacity: 0.3; transform: scale(0.8); }
		50% { opacity: 1; transform: scale(1.15); }
	}

	.anim-pop-them {
		animation: popUpBubbleThem 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) both;
	}
	.anim-pop-me {
		animation: popUpBubbleMe 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) both;
	}

	.typing-bar {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 8px 14px;
		background: #ffffff;
		border: 1px solid #e2e8f0;
		border-radius: 999px;
		width: fit-content;
		font-size: 11px;
		color: #64748b;
		font-family: 'Space Grotesk', sans-serif;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
	}

	.typing-bar.them {
		align-self: flex-start;
	}

	.typing-bar.me {
		align-self: flex-end;
		background: #e0f2fe;
		border-color: #bae6fd;
		color: #00658d;
	}

	.typing-dot {
		width: 6px;
		height: 6px;
		background-color: #00aeef;
		border-radius: 50%;
		display: inline-block;
		animation: dotBounce 1.2s infinite ease-in-out;
	}

	.typing-dot:nth-child(2) {
		animation-delay: 0s;
	}

	.typing-dot:nth-child(3) {
		animation-delay: 0.2s;
	}

	.typing-dot:nth-child(4) {
		animation-delay: 0.4s;
	}

	@keyframes dotBounce {
		0%, 80%, 100% {
			transform: scale(0.6);
			opacity: 0.4;
		}
		40% {
			transform: scale(1.3) translateY(-3px);
			opacity: 1;
		}
	}

	@keyframes blinkCursor {
		0%, 100% { opacity: 1; }
		50% { opacity: 0; }
	}

	.chat-input-bar {
		padding: 14px 20px;
		background: white;
		border-top: 1px solid #f1f5f9;
		display: flex;
		gap: 10px;
		align-items: center;
	}

	.input-fake {
		flex: 1;
		display: flex;
		gap: 8px;
		align-items: center;
		background: #f1f5f9;
		padding: 10px 14px;
		border-radius: 20px;
	}

	.input-fake span {
		font-size: 13px;
		color: #94a3b8;
	}

	.send-btn {
		width: 36px;
		height: 36px;
		border: none;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}

	/* ── SIAP MEMULAI ── */
	.cta-section {
		padding: 100px 48px;
		background: #f0f9ff;
		text-align: center;
	}

	.cta-section h2 {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 40px;
		font-weight: 600;
		color: #0f1c2c;
		margin-bottom: 12px;
	}

	.section-sub {
		font-size: 16px;
		color: #3e4850;
		margin-bottom: 56px;
		line-height: 24px;
		max-width: 560px;
		margin-left: auto;
		margin-right: auto;
	}

	.cta-cards {
		display: flex;
		gap: 24px;
		justify-content: center;
		max-width: 760px;
		margin: 0 auto;
	}

	.cta-card {
		background: white;
		border: 2px solid #cbd5e1;
		border-radius: 16px;
		padding: 40px 40px 44.01px 40px;
		width: 348px;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		transition: transform 0.25s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.25s ease, border-color 0.25s ease;
	}

	.cta-card:hover {
		transform: translateY(-6px);
		box-shadow: 0 16px 32px -6px rgba(0, 101, 141, 0.15);
		border-color: #00aeef;
	}

	.cta-card.featured {
		border: 2px solid #00aeef;
		box-shadow: 0 0 0 4px rgba(0, 174, 239, 0.08);
	}

	.cta-icon {
		margin-bottom: 20px;
	}

	.cta-card h3 {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 20px;
		font-weight: 700;
		color: #0f1c2c;
		margin-bottom: 12px;
	}

	.cta-card p {
		font-size: 14px;
		color: #475569;
		line-height: 1.7;
		margin-bottom: 28px;
	}

	.btn-teal {
		background: #066c8b;
		color: white;
		border: none;
		padding: 14px 24px;
		border-radius: 8px;
		font-weight: 700;
		font-size: 15px;
		cursor: pointer;
		font-family: 'Inter', sans-serif;
	}

	.btn-outline-card {
		background: transparent;
		color: #0f1c2c;
		border: 2px solid #94a3b8;
		padding: 14px 24px;
		border-radius: 8px;
		font-weight: 700;
		font-size: 15px;
		cursor: pointer;
		font-family: 'Inter', sans-serif;
	}

	.full {
		width: 100%;
	}

	/* ── 3 LANGKAH ── */
	.steps {
		padding: 100px 48px;
		background: white;
		text-align: center;
	}

	.workflow-label {
		display: block;
		font-family: 'JetBrains Mono', monospace;
		font-size: 14px;
		font-weight: 500;
		color: #00658d;
		letter-spacing: 2px;
		text-transform: uppercase;
		margin-bottom: 14px;
	}

	.steps h2 {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 40px;
		font-weight: 700;
		color: #0f1c2c;
		margin-bottom: 64px;
	}

	.steps-row {
		display: flex;
		gap: 48px;
		justify-content: center;
		max-width: 960px;
		margin: 0 auto;
	}

	/* ── SLOW & MAJESTIC SCROLL REVEALS (1.6s Duration) ── */

	/* Section Titles & Subtitles Scroll Reveal */
	.features h2,
	.features .section-sub,
	.steps h2,
	.steps .workflow-label,
	.cta-section h2,
	.cta-section .section-sub {
		opacity: 0;
		transform: translateY(40px);
		transition: opacity 1.4s cubic-bezier(0.16, 1, 0.3, 1), transform 1.4s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.features.revealed h2,
	.steps.revealed h2,
	.cta-section.revealed h2 {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.1s;
	}

	.features.revealed .section-sub,
	.steps.revealed .workflow-label,
	.cta-section.revealed .section-sub {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.25s;
	}

	/* 1. Siap Memulai? (CTA Cards - From Bottom) */
	.cta-section .cta-card {
		opacity: 0;
		transform: translateY(70px);
		transition: opacity 1.6s cubic-bezier(0.16, 1, 0.3, 1), transform 1.6s cubic-bezier(0.16, 1, 0.3, 1);
	}
	.cta-section.revealed .cta-card:nth-child(1) {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.4s;
	}
	.cta-section.revealed .cta-card:nth-child(2) {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.7s;
	}

	/* 2. 3 Langkah Sederhana (Workflow - From Sides: Left, Bottom, Right) */
	.steps .step:nth-child(1) {
		opacity: 0;
		transform: translateX(-75px);
		transition: opacity 1.6s cubic-bezier(0.16, 1, 0.3, 1), transform 1.6s cubic-bezier(0.16, 1, 0.3, 1);
	}
	.steps .step:nth-child(2) {
		opacity: 0;
		transform: translateY(70px);
		transition: opacity 1.6s cubic-bezier(0.16, 1, 0.3, 1), transform 1.6s cubic-bezier(0.16, 1, 0.3, 1);
	}
	.steps .step:nth-child(3) {
		opacity: 0;
		transform: translateX(75px);
		transition: opacity 1.6s cubic-bezier(0.16, 1, 0.3, 1), transform 1.6s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.steps.revealed .step:nth-child(1) {
		opacity: 1;
		transform: translateX(0);
		transition-delay: 0.4s;
	}
	.steps.revealed .step:nth-child(2) {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.7s;
	}
	.steps.revealed .step:nth-child(3) {
		opacity: 1;
		transform: translateX(0);
		transition-delay: 1.0s;
	}

	/* 3. Privasi Tanpa Kompromi (4 Feature Cards - From Bottom) */
	.features .feature-card {
		opacity: 0;
		transform: translateY(70px);
		transition: opacity 1.6s cubic-bezier(0.16, 1, 0.3, 1), transform 1.6s cubic-bezier(0.16, 1, 0.3, 1);
	}
	.features.revealed .feature-card:nth-child(1) {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.4s;
	}
	.features.revealed .feature-card:nth-child(2) {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.65s;
	}
	.features.revealed .feature-card:nth-child(3) {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 0.9s;
	}
	.features.revealed .feature-card:nth-child(4) {
		opacity: 1;
		transform: translateY(0);
		transition-delay: 1.15s;
	}

	.step {
		flex: 1;
		text-align: center;
		transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.7s ease;
		cursor: pointer;
		padding: 16px;
		border-radius: 16px;
	}

	.step:hover {
		transform: translateY(-10px);
	}

	.step-circle {
		width: 60px;
		height: 60px;
		background: #e0f2fe;
		border-radius: 9999px;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 auto 24px;
		font-family: 'Space Grotesk', sans-serif;
		font-size: 22px;
		font-weight: 700;
		color: #00658d;
		transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1), background 0.35s ease, color 0.35s ease, box-shadow 0.35s ease;
		box-shadow: 0 4px 12px rgba(0, 101, 141, 0.08);
	}

	.step:hover .step-circle {
		transform: scale(1.25) translateY(-4px);
		background: linear-gradient(135deg, #00c6ff 0%, #00aeef 100%);
		color: #ffffff;
		box-shadow: 0 10px 24px rgba(0, 174, 239, 0.45);
	}

	.step h3 {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 24px;
		font-weight: 600;
		color: #0f1c2c;
		margin-bottom: 12px;
	}

	.step p {
		font-size: 16px;
		color: #3e4850;
		font-family: 'Inter';
		font-size: 16px;
		font-style: normal;
		font-weight: 400;
		line-height: 24px;
	}

	/* ── PRIVASI ── */
	.features {
		padding: 100px 48px;
		background: #f0f9ff;
		text-align: center;
	}

	.features h2 {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 40px;
		font-weight: 700;
		color: #0f1c2c;
		margin-bottom: 20px;
	}

	.feature-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 24px;
		max-width: 1200px;
		margin: 0 auto;
	}

	.feature-card {
		background: white;
		border: 1.5px solid #e2e8f0;
		border-radius: 20px;
		padding: 40px 32px;
		text-align: left;
		display: flex;
		flex-direction: column;
		transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.35s ease, border-color 0.35s ease;
		cursor: pointer;
	}

	.feature-card:hover {
		transform: translateY(-12px) scale(1.02);
		box-shadow: 0 20px 40px -8px rgba(0, 174, 239, 0.28);
		border-color: #00aeef;
	}

	.feature-icon-box {
		width: 52px;
		height: 52px;
		background: #e0f2fe;
		border-radius: 14px;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 24px;
		transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1), background 0.35s ease, box-shadow 0.35s ease;
	}

	.feature-card:hover .feature-icon-box {
		transform: scale(1.25) rotate(6deg) translateY(-4px);
		background: linear-gradient(135deg, #00c6ff 0%, #00aeef 100%);
		box-shadow: 0 8px 20px rgba(0, 174, 239, 0.4);
	}

	.feature-card:hover .feature-icon-box svg path {
		fill: #ffffff;
		transition: fill 0.25s ease;
	}box svg path {
		fill: #ffffff;
		transition: fill 0.25s ease;
	}

	.feature-card h3 {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 24px;
		font-weight: 700;
		color: #0f1c2c;
		margin-bottom: 14px;
		line-height: 1.25;
	}

	.feature-card p {
		font-size: 15px;
		color: #475569;
		line-height: 1.7;
	}

	/* ── FOOTER ── */
	footer {
		background: #0d1b2a;
		padding: 64px 64px 0;
		color: white;
	}

	.footer-top {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding-bottom: 40px;
		border-bottom: 1px solid #1e293b;
		max-width: 1152px;
		margin: 0 auto;
	}

	.footer-brand {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.footer-logo-row {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.footer-logo-img {
		height: 32px;
		width: auto;
		background: white;
		padding: 4px;
		border-radius: 4px;
	}

	.footer-logo-text {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 18px;
		font-weight: 700;
		color: white;
		letter-spacing: 1px;
	}

	.footer-tagline {
		font-family: 'JetBrains Mono', monospace;
		font-size: 13px;
		color: rgba(198, 231, 255, 0.6);
		letter-spacing: 0.7px;
	}

	.footer-links {
		display: flex;
		gap: 32px;
		align-items: center;
	}

	.footer-links a {
		font-family: 'JetBrains Mono', monospace;
		font-size: 14px;
		font-weight: 500;
		color: #94a3b8;
		text-decoration: none;
		transition: color 0.2s;
	}

	.footer-links a:hover {
		color: white;
	}

	.footer-bottom {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 24px 0 40px;
		max-width: 1152px;
		margin: 0 auto;
	}

	.footer-copy {
		font-family: 'JetBrains Mono', monospace;
		font-size: 12px;
		color: #475569;
	}

	.footer-ws {
		font-family: 'JetBrains Mono', monospace;
		font-size: 12px;
		color: #005236;
		display: flex;
		align-items: center;
		gap: 6px;
	}

	/* ── RESPONSIVENESS ── */
	@media (max-width: 1024px) {
		.feature-grid {
			grid-template-columns: repeat(2, 1fr);
			gap: 20px;
		}

		.cta-cards {
			flex-wrap: wrap;
		}
	}

	@media (max-width: 768px) {
		nav {
			padding: 0 20px;
		}

		.nav-center {
			display: none;
		}

		.hero {
			padding: 60px 20px;
		}

		.hero-content {
			flex-direction: column;
			gap: 40px;
			text-align: center;
		}

		.hero-left {
			max-width: 100%;
			display: flex;
			flex-direction: column;
			align-items: center;
		}

		h1 {
			font-size: 38px;
			letter-spacing: -0.8px;
		}

		.hero-sub {
			font-size: 15px;
			margin-bottom: 24px;
		}

		.hero-btns {
			width: 100%;
			flex-direction: column;
			gap: 12px;
		}

		.hero-btns button {
			width: 100%;
		}

		.hero-right {
			width: 100%;
			justify-content: center;
		}

		.chat-card {
			width: 100%;
			max-width: 360px;
		}

		.cta-section {
			padding: 60px 20px;
		}

		.cta-section h2 {
			font-size: 30px;
		}

		.cta-cards {
			flex-direction: column;
			align-items: center;
			gap: 20px;
		}

		.cta-card {
			width: 100%;
			max-width: 348px;
			padding: 30px 24px;
		}

		.steps {
			padding: 60px 20px;
		}

		.steps h2 {
			font-size: 30px;
			margin-bottom: 40px;
		}

		.steps-row {
			flex-direction: column;
			gap: 32px;
		}

		.features {
			padding: 60px 20px;
		}

		.features h2 {
			font-size: 30px;
		}

		.feature-grid {
			grid-template-columns: 1fr;
			gap: 16px;
		}

		.feature-card {
			padding: 30px 24px;
		}

		footer {
			padding: 40px 20px 0;
		}

		.footer-top {
			flex-direction: column;
			gap: 32px;
			align-items: center;
			text-align: center;
		}

		.footer-links {
			flex-direction: column;
			gap: 16px;
		}

		.footer-bottom {
			flex-direction: column;
			gap: 16px;
			text-align: center;
			padding: 24px 0 30px;
		}
	}
</style>
