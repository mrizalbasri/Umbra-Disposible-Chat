<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let nickname = $state('');
  let errorMessage = $state('');

  onMount(() => {
    const params = new URLSearchParams(window.location.search);
    errorMessage = params.get('error') || '';
  });

  function generateRandomNickname() {
    nickname = `Ghost_${Math.floor(1000 + Math.random() * 9000)}`;
  }

  function handleCreate() {
    const finalNickname = nickname.trim() || `Ghost_${Math.floor(1000 + Math.random() * 9000)}`;
    goto(`/chatroom?tab=create&nickname=${encodeURIComponent(finalNickname)}`);
  }
</script>

<svelte:head>
  <title>UMBRA Create Room</title>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Space+Grotesk:wght@500;700&family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="layout">
  <!-- SIDEBAR KIRI -->
  <aside class="sidebar">
    <div class="sidebar-top">
      <!-- Logo -->
      <div class="brand">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="20" viewBox="0 0 16 20" fill="none">
  <path d="M8 20C5.68333 19.4167 3.77083 18.0875 2.2625 16.0125C0.754167 13.9375 0 11.6333 0 9.1V3L8 0L16 3V9.1C16 11.6333 15.2458 13.9375 13.7375 16.0125C12.2292 18.0875 10.3167 19.4167 8 20ZM8 17.9C9.61667 17.4 10.9667 16.4125 12.05 14.9375C13.1333 13.4625 13.7667 11.8167 13.95 10H8V2.125L2 4.375V9.1C2 9.28333 2 9.43333 2 9.55C2 9.66667 2.01667 9.81667 2.05 10H8V17.9Z" fill="#00AEEF"/>
</svg>
        <span class="brand-text">UMBRA</span>
      </div>

      <!-- Tombol Kembali -->
      <button class="btn-back" onclick={() => goto('/')}>
        ← KEMBALI
      </button>

      <!-- Badges -->
      <div class="badges">
        <div class="badge badge-green">
          <span class="dot green"></span> E2EE ACTIVE
        </div>
        <div class="badge badge-outline">
          <svg xmlns="http://www.w3.org/2000/svg" width="15" height="13" viewBox="0 0 15 13" fill="none">
  <path d="M2.325 12.9375C1.6125 12.25 1.04688 11.4406 0.628125 10.5094C0.209375 9.57812 0 8.575 0 7.5C0 6.4625 0.196875 5.4875 0.590625 4.575C0.984375 3.6625 1.51875 2.86875 2.19375 2.19375C2.86875 1.51875 3.6625 0.984375 4.575 0.590625C5.4875 0.196875 6.4625 0 7.5 0C8.5375 0 9.5125 0.196875 10.425 0.590625C11.3375 0.984375 12.1312 1.51875 12.8062 2.19375C13.4812 2.86875 14.0156 3.6625 14.4094 4.575C14.8031 5.4875 15 6.4625 15 7.5C15 8.575 14.7906 9.58125 14.3719 10.5188C13.9531 11.4563 13.3875 12.2625 12.675 12.9375L11.625 11.8875C12.2 11.3375 12.6562 10.6844 12.9937 9.92813C13.3312 9.17188 13.5 8.3625 13.5 7.5C13.5 5.825 12.9188 4.40625 11.7563 3.24375C10.5938 2.08125 9.175 1.5 7.5 1.5C5.825 1.5 4.40625 2.08125 3.24375 3.24375C2.08125 4.40625 1.5 5.825 1.5 7.5C1.5 8.3625 1.66875 9.16875 2.00625 9.91875C2.34375 10.6687 2.80625 11.3187 3.39375 11.8687L2.325 12.9375ZM4.44375 10.8188C4.00625 10.4062 3.65625 9.91563 3.39375 9.34688C3.13125 8.77812 3 8.1625 3 7.5C3 6.25 3.4375 5.1875 4.3125 4.3125C5.1875 3.4375 6.25 3 7.5 3C8.75 3 9.8125 3.4375 10.6875 4.3125C11.5625 5.1875 12 6.25 12 7.5C12 8.1625 11.8687 8.78125 11.6062 9.35625C11.3438 9.93125 10.9937 10.4188 10.5562 10.8188L9.4875 9.75C9.8 9.4625 10.0469 9.125 10.2281 8.7375C10.4094 8.35 10.5 7.9375 10.5 7.5C10.5 6.675 10.2062 5.96875 9.61875 5.38125C9.03125 4.79375 8.325 4.5 7.5 4.5C6.675 4.5 5.96875 4.79375 5.38125 5.38125C4.79375 5.96875 4.5 6.675 4.5 7.5C4.5 7.95 4.59062 8.36562 4.77187 8.74687C4.95312 9.12813 5.2 9.4625 5.5125 9.75L4.44375 10.8188ZM7.5 9C7.0875 9 6.73438 8.85312 6.44063 8.55937C6.14688 8.26562 6 7.9125 6 7.5C6 7.0875 6.14688 6.73438 6.44063 6.44063C6.73438 6.14688 7.0875 6 7.5 6C7.9125 6 8.26562 6.14688 8.55937 6.44063C8.85312 6.73438 9 7.0875 9 7.5C9 7.9125 8.85312 8.26562 8.55937 8.55937C8.26562 8.85312 7.9125 9 7.5 9Z" fill="#10B981"/>
</svg>
          WS: CONNECTED
        </div>
      </div>

      <!-- Panduan Sesi -->
      <div class="guide">
        <h3>PANDUAN SESI</h3>
        <p>Umbra Protocol menjamin enkripsi end-to-end pada setiap transmisi. Room akan otomatis hancur setelah semua peserta keluar atau sesi berakhir.</p>
      </div>
    </div>

    <div class="sidebar-bottom">
      <!-- Tombol End Session -->
      <button class="btn-end" onclick={() => goto('/')}>
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
        End Session
      </button>
      
      <!-- Copyright Text -->
      <p class="copyright">
        UMBRA PROTOCOL © 2024.<br>ALL COMMUNICATIONS ARE EPHEMERAL.
      </p>
    </div>
  </aside>

  <!-- KONTEN UTAMA KANAN -->
  <main class="main-content">
    <div class="watermark">SECURE</div>

    <div class="header-text">
      <h1>Buat Room Baru</h1>
      <p>Mulai sesi komunikasi terenkripsi secara instan.</p>
    </div>

    <div class="form-card">
      <div class="anon-notice">IDENTITAS ANDA AKAN TETAP ANONIM DI JARINGAN</div>

      {#if errorMessage}
        <div class="error-box">
          {errorMessage}
        </div>
      {/if}
      
      <div class="input-group">
        <div class="input-labels">
          <label>NICKNAME (OPSIONAL)</label>
          <span>{nickname.length}/12</span>
        </div>
        <!-- Input Wrapper (Sesuai desain height 56px) -->
        <div class="input-wrapper">
          <input type="text" placeholder="Ghost_42" maxlength="12" bind:value={nickname} />
          <button class="btn-random" onclick={generateRandomNickname}>
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
  <path d="M7 12.25C6.65 12.25 6.35417 12.1292 6.1125 11.8875C5.87083 11.6458 5.75 11.35 5.75 11C5.75 10.65 5.87083 10.3542 6.1125 10.1125C6.35417 9.87083 6.65 9.75 7 9.75C7.35 9.75 7.64583 9.87083 7.8875 10.1125C8.12917 10.3542 8.25 10.65 8.25 11C8.25 11.35 8.12917 11.6458 7.8875 11.8875C7.64583 12.1292 7.35 12.25 7 12.25ZM13 12.25C12.65 12.25 12.3542 12.1292 12.1125 11.8875C11.8708 11.6458 11.75 11.35 11.75 11C11.75 10.65 11.8708 10.3542 12.1125 10.1125C12.3542 9.87083 12.65 9.75 13 9.75C13.35 9.75 13.6458 9.87083 13.8875 10.1125C14.1292 10.3542 14.25 10.65 14.25 11C14.25 11.35 14.1292 11.6458 13.8875 11.8875C13.6458 12.1292 13.35 12.25 13 12.25ZM10 18C12.2333 18 14.125 17.225 15.675 15.675C17.225 14.125 18 12.2333 18 10C18 9.6 17.975 9.2125 17.925 8.8375C17.875 8.4625 17.7833 8.1 17.65 7.75C17.3 7.83333 16.95 7.89583 16.6 7.9375C16.25 7.97917 15.8833 8 15.5 8C13.9833 8 12.55 7.675 11.2 7.025C9.85 6.375 8.7 5.46667 7.75 4.3C7.21667 5.6 6.45417 6.72917 5.4625 7.6875C4.47083 8.64583 3.31667 9.36667 2 9.85C2 9.88333 2 9.90833 2 9.925C2 9.94167 2 9.96667 2 10C2 12.2333 2.775 14.125 4.325 15.675C5.875 17.225 7.76667 18 10 18ZM10 20C8.61667 20 7.31667 19.7375 6.1 19.2125C4.88333 18.6875 3.825 17.975 2.925 17.075C2.025 16.175 1.3125 15.1167 0.7875 13.9C0.2625 12.6833 0 11.3833 0 10C0 8.61667 0.2625 7.31667 0.7875 6.1C1.3125 4.88333 2.025 3.825 2.925 2.925C3.825 2.025 4.88333 1.3125 6.1 0.7875C7.31667 0.2625 8.61667 0 10 0C11.3833 0 12.6833 0.2625 13.9 0.7875C15.1167 1.3125 16.175 2.025 17.075 2.925C17.975 3.825 18.6875 4.88333 19.2125 6.1C19.7375 7.31667 20 8.61667 20 10C20 11.3833 19.7375 12.6833 19.2125 13.9C18.6875 15.1167 17.975 16.175 17.075 17.075C16.175 17.975 15.1167 18.6875 13.9 19.2125C12.6833 19.7375 11.3833 20 10 20ZM8.65 2.125C9.35 3.29167 10.3 4.22917 11.5 4.9375C12.7 5.64583 14.0333 6 15.5 6C15.7333 6 15.9583 5.9875 16.175 5.9625C16.3917 5.9375 16.6167 5.90833 16.85 5.875C16.15 4.70833 15.2 3.77083 14 3.0625C12.8 2.35417 11.4667 2 10 2C9.76667 2 9.54167 2.0125 9.325 2.0375C9.10833 2.0625 8.88333 2.09167 8.65 2.125ZM2.425 7.475C3.275 6.99167 4.01667 6.36667 4.65 5.6C5.28333 4.83333 5.75833 3.975 6.075 3.025C5.225 3.50833 4.48333 4.13333 3.85 4.9C3.21667 5.66667 2.74167 6.525 2.425 7.475Z" fill="#6E7881"/>
</svg>
          </button>
        </div>
      </div>

      <button class="btn-create" onclick={handleCreate}>
        Buat Room 
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
  <path d="M9 15H11V11H15V9H11V5H9V9H5V11H9V15ZM10 20C8.61667 20 7.31667 19.7375 6.1 19.2125C4.88333 18.6875 3.825 17.975 2.925 17.075C2.025 16.175 1.3125 15.1167 0.7875 13.9C0.2625 12.6833 0 11.3833 0 10C0 8.61667 0.2625 7.31667 0.7875 6.1C1.3125 4.88333 2.025 3.825 2.925 2.925C3.825 2.025 4.88333 1.3125 6.1 0.7875C7.31667 0.2625 8.61667 0 10 0C11.3833 0 12.6833 0.2625 13.9 0.7875C15.1167 1.3125 16.175 2.025 17.075 2.925C17.975 3.825 18.6875 4.88333 19.2125 6.1C19.7375 7.31667 20 8.61667 20 10C20 11.3833 19.7375 12.6833 19.2125 13.9C18.6875 15.1167 17.975 16.175 17.075 17.075C16.175 17.975 15.1167 18.6875 13.9 19.2125C12.6833 19.7375 11.3833 20 10 20ZM10 18C12.2333 18 14.125 17.225 15.675 15.675C17.225 14.125 18 12.2333 18 10C18 7.76667 17.225 5.875 15.675 4.325C14.125 2.775 12.2333 2 10 2C7.76667 2 5.875 2.775 4.325 4.325C2.775 5.875 2 7.76667 2 10C2 12.2333 2.775 14.125 4.325 15.675C5.875 17.225 7.76667 18 10 18Z" fill="white"/>
</svg>
      </button>

      <div class="join-link">
        Join ke room teman ? <a href="/join" class="nama-class-css-kamu">Klik disini</a>
      </div>
    </div>
  </main>
</div>

<style>
  /* RESET & BASE */
  :global(body) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    font-family: 'Inter', sans-serif;
    background: #F8FAFC;
  }

  .layout {
    display: flex;
    height: 100vh;
    width: 100%;
    overflow: hidden;
  }

  /* --- SIDEBAR KIRI --- */
  .sidebar {
    width: 280px;
    background: #F8FAFC; 
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 40px 32px;
    flex-shrink: 0;
  }
.btn-create {
    width: 100%;
    height: 56px;
    background: #00AEEF;
    color: white;
    border: none;
    border-radius: 10px;
    font-family: 'Inter', sans-serif;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s;
    margin-bottom: 24px;
    
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px; 
  }
  .sidebar-top { display: flex; flex-direction: column; gap: 32px; }

  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .brand-text {
    color: #00658D;
    font-family: 'Space Grotesk', sans-serif;
    font-size: 24px;
    font-weight: 700;
    line-height: 31.2px;
    letter-spacing: -0.6px;
  }

  .btn-back {
    background: transparent;
    border: none;
    color: #475569;
    font-family: 'JetBrains Mono', monospace;
    font-size: 12px;
    font-weight: 600;
    text-align: left;
    cursor: pointer;
    padding: 0;
    letter-spacing: 0.5px;
  }

  .badges { display: flex; flex-direction: column; gap: 12px; align-items: flex-start; }
  
  .badge {
    padding: 6px 12px;
    border-radius: 99px;
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    font-weight: 700;
    display: flex;
    align-items: center;
    gap: 8px;
    letter-spacing: 0.5px;
  }
  .badge-green { background: rgba(255, 255, 255, 0.50); color: #16A34A; }
  .badge-outline { border: 1px solid #CBD5E1; color: #64748B; }
  .dot.green { width: 6px; height: 6px; background: #16A34A; border-radius: 50%; }
  .ws-icon { color: #CBD5E1; }

  .guide h3 {
    color: #00658D;
    font-family: 'Space Grotesk', sans-serif;
    font-size: 12px; 
    font-weight: 700;
    margin-bottom: 8px;
    letter-spacing: 0.5px;
  }
  .guide p {
    color: #64748B;
    font-size: 12px;
    line-height: 1.6;
    margin: 0;
  }

  .sidebar-bottom { display: flex; flex-direction: column; gap: 24px; }

  .btn-end {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    width: 100%;
    padding: 12px;
    background: transparent;
    border: 1px solid #BA1A1A;
    border-radius: 8px;
    color: #BA1A1A;
    font-family: 'JetBrains Mono', monospace;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
  }

  .copyright {
    color: #6E7881;
    font-family: 'JetBrains Mono', monospace;
    font-size: 10px;
    font-weight: 400;
    line-height: 15px;
    letter-spacing: -0.5px;
    text-transform: uppercase;
    margin: 0;
  }

  /* --- KONTEN UTAMA KANAN --- */
  .main-content {
    flex: 1;
    background: #FFFFFF; /* Di Image 1 area kanan itu putih bersih */
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  /* Posisi Watermark SECURE */
  .watermark {
    position: absolute;
    top: 24px;   
    right: 24px; 
    font-family: 'JetBrains Mono', monospace;
    font-size: 120px;
    font-weight: 700;
    line-height: 120px;
    color: #BDC8D1; 
    opacity: 0.2; 
    z-index: 0;
  }

  .header-text {
    text-align: center;
    margin-bottom: 40px;
    z-index: 1;
  }
  .header-text h1 {
    font-family: 'Space Grotesk', sans-serif;
    font-size: 32px;
    font-weight: 700;
    color: #0F1C2C;
    margin: 0 0 8px 0;
  }
  .header-text p {
    color: #64748B;
    font-size: 15px;
    margin: 0;
  }

  /* REVISI BORDER: Form Card */
  .form-card {
    background: #FFFFFF;
    border: 1px solid #E2E8F0; 
    border-radius: 12px; 
    padding: 40px;
    width: 100%;
    max-width: 440px;
    z-index: 1;
  }

  .anon-notice {
    color: #6E7881;
    text-align: center;
    font-family: 'JetBrains Mono', monospace;
    font-size: 12px;
    font-weight: 400;
    line-height: 18px;
    margin-bottom: 32px;
  }

  .input-group { margin-bottom: 24px; }
  
  .input-labels {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    font-weight: 600;
    color: #475569;
  }

  .input-wrapper {
    display: flex;
    height: 56px;
    padding: 17px 16px;
    align-items: center;
    border-radius: 10px;
    border: 1px solid #E2E8F0;
    background: #F5F8FA;
    box-sizing: border-box;
  }

  .input-wrapper input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    font-family: 'Inter', sans-serif;
    font-size: 15px;
    color: #0F1C2C;
  }
  .input-wrapper input::placeholder { color: #94A3B8; }
  
  .btn-random {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0;
    display: flex;
    align-items: center;
  }

  .btn-create {
    width: 100%;
    height: 56px;
    background: #00AEEF;
    color: white;
    border: none;
    border-radius: 10px;
    font-family: 'Inter', sans-serif;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s;
    margin-bottom: 24px;
  }
  .btn-create:hover { background: #0096D1; }

  .join-link {
    text-align: center;
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    color: #6E7881;
  }
  .join-link a {
    color: #6E7881;
    text-decoration: none;
    font-weight: 600;
  }

  .error-box {
    background: #FEE2E2;
    border: 1px solid #FCA5A5;
    color: #DC2626;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 13px;
    margin-bottom: 24px;
    text-align: center;
    font-weight: 500;
    font-family: 'Inter', sans-serif;
  }
</style>