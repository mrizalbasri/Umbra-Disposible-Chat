<script lang="ts">
  import { goto } from '$app/navigation';
</script>

<svelte:head>
  <title>UMBRA — Join Room</title>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Space+Grotesk:wght@500;700&family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="dashboard-layout">
  <!-- === TOP NAVBAR (Sesuai Image 3) === -->
  <nav class="navbar">
    <div class="nav-left">
      <div class="nav-brand">
        <span class="brand-blue">UMBRA</span>
        <span class="brand-divider">|</span>
        <span class="brand-gray">PROTOCOL</span>
      </div>
    </div>
    
    <div class="nav-center">
      <a href="#features">Features</a>
      <a href="#security">Security</a>
      <a href="#docs">Documentation</a>
    </div>

    <div class="nav-right">
      <span class="lang-toggle">ID/EN</span>
      <button class="btn-start">Start</button>
    </div>
  </nav>

  <div class="main-wrapper">
    <!-- === LEFT SIDEBAR === -->
    <aside class="sidebar">
      <div class="sidebar-top">
        <!-- Anonymous Session Badge (Sesuai Image 2) -->
        <div class="session-badge">
          <div class="avatar-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#00AEEF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
          </div>
          <div class="session-info">
            <span class="session-title">Anonymous Session</span>
            <span class="session-status"><span class="dot green"></span> E2EE: READY</span>
          </div>
        </div>

        <!-- Navigation Menu -->
        <ul class="side-menu">
          <li class="active">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="20" viewBox="0 0 16 20" fill="none">
  <path d="M8 20C5.68333 19.4167 3.77083 18.0875 2.2625 16.0125C0.754167 13.9375 0 11.6333 0 9.1V3L8 0L16 3V9.1C16 11.6333 15.2458 13.9375 13.7375 16.0125C12.2292 18.0875 10.3167 19.4167 8 20ZM8 17.9C9.61667 17.4 10.9667 16.4125 12.05 14.9375C13.1333 13.4625 13.7667 11.8167 13.95 10H8V2.125L2 4.375V9.1C2 9.28333 2 9.43333 2 9.55C2 9.66667 2.01667 9.81667 2.05 10H8V17.9Z" fill="#3E4850"/>
</svg>
            Secure Room
          </li>
          <li>
            <svg xmlns="http://www.w3.org/2000/svg" width="22" height="16" viewBox="0 0 22 16" fill="none">
  <path d="M0 16V13.2C0 12.6333 0.145833 12.1125 0.4375 11.6375C0.729167 11.1625 1.11667 10.8 1.6 10.55C2.63333 10.0333 3.68333 9.64583 4.75 9.3875C5.81667 9.12917 6.9 9 8 9C9.1 9 10.1833 9.12917 11.25 9.3875C12.3167 9.64583 13.3667 10.0333 14.4 10.55C14.8833 10.8 15.2708 11.1625 15.5625 11.6375C15.8542 12.1125 16 12.6333 16 13.2V16H0ZM18 16V13C18 12.2667 17.7958 11.5625 17.3875 10.8875C16.9792 10.2125 16.4 9.63333 15.65 9.15C16.5 9.25 17.3 9.42083 18.05 9.6625C18.8 9.90417 19.5 10.2 20.15 10.55C20.75 10.8833 21.2083 11.2542 21.525 11.6625C21.8417 12.0708 22 12.5167 22 13V16H18ZM8 8C6.9 8 5.95833 7.60833 5.175 6.825C4.39167 6.04167 4 5.1 4 4C4 2.9 4.39167 1.95833 5.175 1.175C5.95833 0.391667 6.9 0 8 0C9.1 0 10.0417 0.391667 10.825 1.175C11.6083 1.95833 12 2.9 12 4C12 5.1 11.6083 6.04167 10.825 6.825C10.0417 7.60833 9.1 8 8 8ZM18 4C18 5.1 17.6083 6.04167 16.825 6.825C16.0417 7.60833 15.1 8 14 8C13.8167 8 13.5833 7.97917 13.3 7.9375C13.0167 7.89583 12.7833 7.85 12.6 7.8C13.05 7.26667 13.3958 6.675 13.6375 6.025C13.8792 5.375 14 4.7 14 4C14 3.3 13.8792 2.625 13.6375 1.975C13.3958 1.325 13.05 0.733333 12.6 0.2C12.8333 0.116667 13.0667 0.0625 13.3 0.0375C13.5333 0.0125 13.7667 0 14 0C15.1 0 16.0417 0.391667 16.825 1.175C17.6083 1.95833 18 2.9 18 4ZM2 14H14V13.2C14 13.0167 13.9542 12.85 13.8625 12.7C13.7708 12.55 13.65 12.4333 13.5 12.35C12.6 11.9 11.6917 11.5625 10.775 11.3375C9.85833 11.1125 8.93333 11 8 11C7.06667 11 6.14167 11.1125 5.225 11.3375C4.30833 11.5625 3.4 11.9 2.5 12.35C2.35 12.4333 2.22917 12.55 2.1375 12.7C2.04583 12.85 2 13.0167 2 13.2V14ZM8 6C8.55 6 9.02083 5.80417 9.4125 5.4125C9.80417 5.02083 10 4.55 10 4C10 3.45 9.80417 2.97917 9.4125 2.5875C9.02083 2.19583 8.55 2 8 2C7.45 2 6.97917 2.19583 6.5875 2.5875C6.19583 2.97917 6 3.45 6 4C6 4.55 6.19583 5.02083 6.5875 5.4125C6.97917 5.80417 7.45 6 8 6Z" fill="#3E4850"/>
</svg>
            Participants
          </li>
          <li>
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="16" viewBox="0 0 20 16" fill="none">
  <path d="M9 13H17V12.45C17 11.7 16.6333 11.1042 15.9 10.6625C15.1667 10.2208 14.2 10 13 10C11.8 10 10.8333 10.2208 10.1 10.6625C9.36667 11.1042 9 11.7 9 12.45V13ZM13 9C13.55 9 14.0208 8.80417 14.4125 8.4125C14.8042 8.02083 15 7.55 15 7C15 6.45 14.8042 5.97917 14.4125 5.5875C14.0208 5.19583 13.55 5 13 5C12.45 5 11.9792 5.19583 11.5875 5.5875C11.1958 5.97917 11 6.45 11 7C11 7.55 11.1958 8.02083 11.5875 8.4125C11.9792 8.80417 12.45 9 13 9ZM2 16C1.45 16 0.979167 15.8042 0.5875 15.4125C0.195833 15.0208 0 14.55 0 14V2C0 1.45 0.195833 0.979167 0.5875 0.5875C0.979167 0.195833 1.45 0 2 0H8L10 2H18C18.55 2 19.0208 2.19583 19.4125 2.5875C19.8042 2.97917 20 3.45 20 4V14C20 14.55 19.8042 15.0208 19.4125 15.4125C19.0208 15.8042 18.55 16 18 16H2ZM2 14H18V4H9.175L7.175 2H2V14ZM2 14V2V4V14Z" fill="#3E4850"/>
</svg>
            Media Vault
          </li>
          <li>
            <svg class="menu-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
            Settings
          </li>
        </ul>
      </div>

      <div class="sidebar-bottom">
        <button class="btn-end" on:click={() => goto('/')}>
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
          End Session
        </button>
        <div class="ws-status">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
          WEBSOCKET: CONNECTED
        </div>
      </div>
    </aside>

    <!-- === MAIN CONTENT (Sesuai Image 1) === -->
    <main class="content-area">
      <div class="form-card">
        <div class="header-text">
          <h1>Masuk ke Room</h1>
          <p>Masukkan kode akses 6 digit untuk bergabung ke dalam sesi enkripsi end-to-end yang aman.</p>
        </div>

        <!-- KODE ROOM (OTP Boxes) -->
        <div class="input-group">
          <label>KODE ROOM</label>
          <div class="otp-container">
            <input type="text" maxlength="1" class="otp-box" />
            <input type="text" maxlength="1" class="otp-box" />
            <input type="text" maxlength="1" class="otp-box" />
            <input type="text" maxlength="1" class="otp-box" />
            <input type="text" maxlength="1" class="otp-box" />
            <input type="text" maxlength="1" class="otp-box" />
          </div>
        </div>

        <!-- NAMA SAMARAN -->
        <div class="input-group">
          <label>NAMA SAMARAN (OPSIONAL)</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <input type="text" placeholder="Contoh: Ghost_99" />
          </div>
        </div>

        <button class="btn-join" on:click={() => goto('/chatroom')}>
          Gabung Room 
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
        </button>

        <div class="security-badge">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
          ZERO-KNOWLEDGE ARCHITECTURE ACTIVE
        </div>
      </div>
    </main>
  </div>
</div>

<style>
  /* RESET & BASE */
  :global(body) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    font-family: 'Inter', sans-serif;
    background: #F8F9FF; /* Background slightly bluish based on Figma */
  }

  .dashboard-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100%;
    overflow: hidden;
  }

  /* === NAVBAR (Image 3) === */
  .navbar {
    display: flex;
    width: 100%;
    height: 64px;
    padding: 0 24px;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid rgba(189, 196, 209, 0.3);
    background: rgba(248, 249, 255, 0.85);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    box-sizing: border-box;
    z-index: 10;
  }

  .nav-brand {
    font-family: 'Space Grotesk', sans-serif;
    font-size: 24px;
    font-weight: 700;
    display: flex;
    gap: 8px;
    align-items: center;
  }
  .brand-blue { color: #00AEEF; }
  .brand-divider { color: #CBD5E1; font-weight: 400; }
  .brand-gray { color: #64748B; }

  .nav-center {
    display: flex;
    gap: 32px;
  }
  .nav-center a {
    text-decoration: none;
    color: #475569;
    font-size: 14px;
    font-weight: 500;
    transition: color 0.2s;
  }
  .nav-center a:hover { color: #00AEEF; }

  .nav-right {
    display: flex;
    align-items: center;
    gap: 24px;
  }
  .lang-toggle {
    font-size: 14px;
    font-weight: 600;
    color: #64748B;
    cursor: pointer;
  }
  .btn-start {
    background: #00AEEF;
    color: white;
    border: none;
    padding: 8px 20px;
    border-radius: 6px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
  }

  /* === MAIN WRAPPER === */
  .main-wrapper {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  /* === SIDEBAR (Image 2) === */
  .sidebar {
    width: 260px;
    background: #F8F9FF;
    border-right: 1px solid #E2E8F0;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 24px;
    flex-shrink: 0;
  }

  .sidebar-top { display: flex; flex-direction: column; gap: 32px; }

  /* Anonymous Session Typography (Image 2) */
  .session-badge {
    display: flex;
    align-items: center;
    gap: 12px;
    background: #E8F4F8;
    padding: 12px;
    border-radius: 12px;
    border: 1px solid #D6E8ED;
  }
  .avatar-icon {
    width: 32px;
    height: 32px;
    background: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .session-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .session-title {
    color: #00658D;
    font-family: 'JetBrains Mono', monospace;
    font-size: 14px;
    font-weight: 700;
    line-height: 16.8px;
    letter-spacing: 0.7px;
  }
  .session-status {
    font-family: 'Inter', sans-serif;
    font-size: 10px;
    font-weight: 600;
    color: #64748B;
    display: flex;
    align-items: center;
    gap: 4px;
  }
  .dot.green { width: 6px; height: 6px; background: #16A34A; border-radius: 50%; }

  .side-menu {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .side-menu li {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    border-radius: 8px;
    color: #475569;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }
  .side-menu li:hover { background: #F1F5F9; color: #0F1C2C; }
  .side-menu li.active { background: #FFFFFF; font-weight: 600; color: #0F1C2C; box-shadow: 0 1px 2px rgba(0,0,0,0.05); }
  .menu-icon { width: 18px; height: 18px; color: #64748B; }

  .sidebar-bottom { display: flex; flex-direction: column; gap: 16px; }

  .btn-end {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    width: 100%;
    padding: 10px;
    background: transparent;
    border: 1px solid #EF4444;
    border-radius: 8px;
    color: #EF4444;
    font-family: 'JetBrains Mono', monospace;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
  }
  
  .ws-status {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    color: #94A3B8;
    font-family: 'JetBrains Mono', monospace;
    font-size: 9px;
    font-weight: 600;
  }

  /* === MAIN CONTENT (Image 1) === */
  .content-area {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    background: #FFFFFF;
    position: relative;
  }

  .form-card {
    background: #FFFFFF;
    border: 1px solid #E2E8F0;
    border-radius: 16px;
    padding: 48px;
    width: 100%;
    max-width: 460px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.03);
  }

  .header-text { margin-bottom: 32px; }
  .header-text h1 {
    font-family: 'Inter', sans-serif;
    font-size: 24px;
    font-weight: 700;
    color: #0F1C2C;
    margin: 0 0 12px 0;
  }
  /* Typography Subtitle (Image 1) */
  .header-text p {
    color: #3E4850;
    font-family: 'Inter', sans-serif;
    font-size: 16px;
    font-style: normal;
    font-weight: 400;
    line-height: 24px;
    margin: 0;
  }

  .input-group { margin-bottom: 24px; }
  .input-group label {
    display: block;
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    font-weight: 600;
    color: #64748B;
    margin-bottom: 12px;
    letter-spacing: 0.5px;
  }

  /* OTP Boxes */
  .otp-container {
    display: flex;
    gap: 12px;
    justify-content: space-between;
  }
  .otp-box {
    width: 48px;
    height: 56px;
    border-radius: 10px;
    border: 1px solid #E2E8F0;
    background: #FFFFFF;
    text-align: center;
    font-family: 'Inter', sans-serif;
    font-size: 20px;
    font-weight: 600;
    color: #0F1C2C;
    outline: none;
    transition: border-color 0.2s;
  }
  .otp-box:focus { border-color: #00AEEF; box-shadow: 0 0 0 3px rgba(0, 174, 239, 0.1); }

  /* Input Biasa */
  .input-wrapper {
    display: flex;
    height: 56px;
    padding: 0 16px;
    align-items: center;
    border-radius: 10px;
    border: 1px solid #E2E8F0;
    background: #F8FAFC;
    gap: 12px;
  }
  .input-icon { width: 18px; height: 18px; color: #94A3B8; }
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

  /* Button Gabung */
  .btn-join {
    width: 100%;
    height: 56px;
    background: #9BBCC4; /* Warna muted teal/blue sesuai Figma */
    color: white;
    border: none;
    border-radius: 10px;
    font-family: 'Inter', sans-serif;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    transition: background 0.2s;
    margin-bottom: 24px;
    margin-top: 16px;
  }
  .btn-join:hover { background: #88A9B1; }

  /* Bottom Security Text */
  .security-badge {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    color: #10B981;
    font-family: 'JetBrains Mono', monospace;
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.5px;
  }
</style>