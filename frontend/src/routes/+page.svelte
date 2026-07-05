<script lang="ts">
  import { goto } from '$app/navigation';
</script>

<svelte:head>
  <title>UMBRA — Secure Anonymous Chat</title>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
</svelte:head>

<main>
  <!-- NAVBAR -->
  <nav>
    <div class="nav-left">
      <!-- Menggunakan tag img untuk logo sesuai permintaan -->
      <img src="/logo.png" alt="UMBRA Logo" class="logo-img" />
      <span class="logo">UMBRA</span>
    </div>
    <div class="nav-center">
      <a href="#features" class="active">Features</a>
      <a href="#security">Security</a>
      <a href="#docs">Documentation</a>
    </div>
    <div class="nav-right">
      <span class="lang">ID/EN</span>
      <button class="btn-start" on:click={() => goto('/create')}>Mulai</button>
    </div>
  </nav>

  <!-- HERO -->
  <section class="hero">
    <div class="hero-left">
      <span class="badge">● END-TO-END ENCRYPTED</span>
      <h1>Bicara Bebas,<br><span class="accent">Tanpa Jejak,</span><br>Tanpa Identitas.</h1>
      <p>Platform komunikasi paling aman untuk percakapan rahasia. Tidak ada login, tidak ada pelacakan, dan data Anda hilang seketika setelah selesai.</p>
      <div class="hero-buttons">
        <button class="btn-primary" on:click={() => goto('/create')}>Mulai Chat Sekarang →</button>
        <button class="btn-outline-hero" on:click={() => goto('/docs')}>Pelajari Protokol</button>
      </div>
    </div>
    <div class="hero-right">
      <div class="chat-preview">
        <div class="chat-header">
          <div class="chat-header-left">
            <div class="avatar">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#0f3460" stroke-width="2">
                <circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
              </svg>
            </div>
            <div>
              <div class="chat-name">Anonymous User #294</div>
              <div class="chat-status">● Secure Tunnel Active</div>
            </div>
          </div>
          <div class="shield-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
          </div>
        </div>
        <div class="chat-body">
          <div class="msg-left">Apakah data ini benar-benar tidak tersimpan di server?</div>
          <div class="msg-right">Tepat sekali. Menggunakan enkripsi ECDH dan Double Ratchet. Segera setelah tab ditutup, kunci dibuang.</div>
          <div class="msg-left">Luar biasa. Mulai kirim dokumen rahasianya.</div>
        </div>
        <div class="chat-input">
          <div class="input-fake">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M21.44 11.05l-9.19 9.19a6 6 0 01-8.49-8.49l9.19-9.19a4 4 0 015.66 5.66l-9.2 9.19a2 2 0 01-2.83-2.83l8.49-8.48"/></svg>
            <span>Type secure message...</span>
          </div>
          <button class="send-btn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>
          </button>
        </div>
      </div>
    </div>
  </section>

  <!-- SIAP MEMULAI (Posisi disesuaikan urutan Figma) -->
  <section class="cta" id="cta">
    <h2>Siap Memulai?</h2>
    <p class="section-subtitle">Pilih bagaimana Anda ingin terhubung dengan orang lain secara anonim.</p>
    <div class="cta-cards">
      <!-- Card 1 (Active) -->
      <div class="cta-card active-card">
        <div class="cta-icon blue-icon">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#00b4d8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg>
        </div>
        <h3>Private Room</h3>
        <p>Buat ruang obrolan eksklusif dan bagikan kode rahasia kepada partner bicara Anda.</p>
        <button class="btn-primary full" on:click={() => goto('/create')}>Buat Ruangan</button>
      </div>
      <!-- Card 2 -->
      <div class="cta-card">
        <div class="cta-icon dark-icon">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 3 21 3 21 8"></polyline><line x1="4" y1="20" x2="21" y2="3"></line><polyline points="21 16 21 21 16 21"></polyline><line x1="15" y1="15" x2="21" y2="21"></line><line x1="4" y1="4" x2="9" y2="9"></line></svg>
        </div>
        <h3>Random Match</h3>
        <p>Terhubung secara acak dengan orang lain di seluruh dunia tanpa mengenal identitas.</p>
        <button class="btn-outline full" on:click={() => goto('/random')}>Cari Teman Chat</button>
      </div>
    </div>
  </section>

  <!-- 3 LANGKAH -->
  <section class="steps" id="security">
    <span class="workflow-label">WORKFLOW</span>
    <h2>3 Langkah Sederhana</h2>
    <div class="steps-row">
      <div class="step">
        <div class="step-num">01</div>
        <h3>Buka UMBRA</h3>
        <p>Akses web UMBRA dari perangkat mana saja. Tidak ada aplikasi yang perlu diunduh.</p>
      </div>
      <div class="step">
        <div class="step-num">02</div>
        <h3>Pilih Mode Chat</h3>
        <p>Buat Private Room dengan kode unik atau cari Random Match secara instan.</p>
      </div>
      <div class="step">
        <div class="step-num">03</div>
        <h3>Chat Aman</h3>
        <p>Mulai berkomunikasi dengan tenang. Semua pesan dienkripsi secara real-time.</p>
      </div>
    </div>
  </section>

  <!-- PRIVASI TANPA KOMPROMI (Layout 4 Kolom) -->
  <section class="features" id="features">
    <h2>Privasi Tanpa Kompromi</h2>
    <p class="section-subtitle">UMBRA dirancang oleh pakar keamanan untuk memberikan perlindungan mutlak bagi identitas dan pesan Anda.</p>
    <div class="feature-cards">
      <div class="feature-card">
        <div class="feature-icon-box">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#00b4d8" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
        </div>
        <h3>Enkripsi End-to-End</h3>
        <p>Hanya Anda dan penerima yang dapat membaca pesan. Tidak ada pihak ketiga yang bisa mencegat.</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon-box">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#00b4d8" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
        </div>
        <h3>Tanpa Identitas</h3>
        <p>Tidak perlu nomor telepon, email, atau akun sosial media. Anda anonim secara default.</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon-box">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#00b4d8" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
        </div>
        <h3>Ephemeral Chat</h3>
        <p>Pesan akan hancur sendiri setelah waktu yang ditentukan atau setelah sesi berakhir.</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon-box">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#00b4d8" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
        </div>
        <h3>Zero Knowledge</h3>
        <p>Server kami tidak pernah menyimpan data chat atau metadata yang bisa mengidentifikasi Anda.</p>
      </div>
    </div>
  </section>

  <!-- FOOTER -->
  <footer>
    <div class="footer-left">
      <img src="/logo.png" alt="UMBRA Logo" class="logo-img-small" />
      <span class="logo small">UMBRA</span>
    </div>
    <div class="footer-links">
      <a href="#audit">Security Audit</a>
      <a href="#privacy">Privacy Policy</a>
      <a href="#github">Github</a>
      <a href="#status">Status</a>
    </div>
    <div class="footer-right">
      <span>© 2026 UMBRA. All rights reserved.</span>
    </div>
  </footer>
</main>

<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }

  main {
    font-family: 'Inter', sans-serif;
    background: #f0f9ff;
    color: #1a1a2e;
  }

  h1, h2, h3 { color: #0f3460; }

  /* NAVBAR */
  nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 48px;
    height: 72px;
    background: white;
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
    height: 32px;
    width: auto;
  }

  .logo {
    font-size: 20px;
    font-weight: 800;
    color: #0f3460;
    letter-spacing: 2px;
  }

  .nav-center {
    display: flex;
    gap: 24px;
  }

  .nav-center a {
    text-decoration: none;
    color: #64748b;
    font-size: 14px;
    font-weight: 600;
  }

  .nav-center a.active {
    color: #00b4d8;
    border-bottom: 2px solid #00b4d8;
    padding-bottom: 4px;
  }

  .lang {
    margin-right: 20px;
    font-size: 14px;
    color: #64748b;
    font-weight: 600;
    cursor: pointer;
  }

  .btn-start {
    background: #00b4d8;
    color: white;
    border: none;
    padding: 10px 24px;
    border-radius: 8px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
    font-family: 'Inter', sans-serif;
  }

  /* HERO */
  .hero {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 80px 48px;
    gap: 48px;
    max-width: 1200px;
    margin: 0 auto;
  }

  .hero-left {
    max-width: 540px;
  }

  .badge {
    background: #e0f2fe;
    color: #00b4d8;
    padding: 6px 16px;
    border-radius: 999px;
    font-size: 12px;
    font-weight: 700;
    display: inline-block;
    margin-bottom: 24px;
    letter-spacing: 0.5px;
  }

  h1 {
    font-size: 56px;
    font-weight: 800;
    line-height: 1.15;
    margin-bottom: 24px;
  }

  .accent { color: #00b4d8; }

  .hero-left p {
    font-size: 16px;
    color: #64748b;
    line-height: 1.7;
    margin-bottom: 40px;
  }

  .hero-buttons {
    display: flex;
    gap: 16px;
  }

  .btn-primary {
    background: #00b4d8;
    color: white;
    border: none;
    padding: 14px 28px;
    border-radius: 8px;
    font-weight: 600;
    font-size: 15px;
    cursor: pointer;
    font-family: 'Inter', sans-serif;
  }

  .btn-outline-hero {
    background: transparent;
    color: #0f3460;
    border: 1px solid #cbd5e1;
    padding: 14px 28px;
    border-radius: 8px;
    font-weight: 600;
    font-size: 15px;
    cursor: pointer;
    font-family: 'Inter', sans-serif;
    background-color: white;
  }

  /* CHAT PREVIEW WIDGET */
  .hero-right {
    flex: 1;
    display: flex;
    justify-content: flex-end;
  }

  .chat-preview {
    width: 380px;
    background: white;
    border-radius: 16px;
    box-shadow: 0 20px 40px rgba(15, 52, 96, 0.08);
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

  .chat-header-left { display: flex; align-items: center; gap: 12px; }
  .avatar { width: 36px; height: 36px; background: #f0f9ff; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
  .chat-name { font-size: 13px; font-weight: 700; color: #0f3460; }
  .chat-status { font-size: 11px; color: #10b981; font-weight: 600; }

  .chat-body {
    padding: 24px 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    background: #f8fafc;
  }

  .msg-left, .msg-right {
    padding: 12px 16px;
    border-radius: 12px;
    font-size: 13px;
    line-height: 1.5;
    max-width: 85%;
  }

  .msg-left { background: white; color: #334155; align-self: flex-start; border: 1px solid #e2e8f0; border-bottom-left-radius: 4px; }
  .msg-right { background: #00b4d8; color: white; align-self: flex-end; border-bottom-right-radius: 4px; }

  .chat-input {
    padding: 16px 20px;
    border-top: 1px solid #f1f5f9;
    display: flex;
    gap: 12px;
    align-items: center;
    background: white;
  }

  .input-fake { flex: 1; display: flex; gap: 8px; background: #f1f5f9; padding: 10px 16px; border-radius: 20px; align-items: center; }
  .input-fake span { font-size: 13px; color: #94a3b8; }
  .send-btn { background: #00b4d8; border: none; width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; }

  /* COMMON SECTION TITLES */
  .section-subtitle {
    font-size: 15px;
    color: #64748b;
    margin-bottom: 48px;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
  }

  /* SIAP MEMULAI CTA */
  .cta {
    padding: 80px 48px;
    text-align: center;
    background: #f0f9ff; /* Sesuai warna background figma */
  }

  .cta h2 {
    font-size: 36px;
    font-weight: 800;
    margin-bottom: 16px;
  }

  .cta-cards {
    display: flex;
    gap: 24px;
    justify-content: center;
    max-width: 1000px;
    margin: 0 auto;
  }

  .cta-card {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    padding: 40px 32px;
    width: 380px;
  }

  .active-card {
    border: 2px solid #00b4d8;
    box-shadow: 0 10px 25px rgba(0, 180, 216, 0.1);
  }

  .cta-icon { margin-bottom: 24px; display: flex; justify-content: center; }

  .cta-card h3 {
    font-size: 20px;
    font-weight: 700;
    margin-bottom: 12px;
  }

  .cta-card p {
    font-size: 14px;
    color: #64748b;
    line-height: 1.6;
    margin-bottom: 32px;
  }

  .btn-outline {
    background: transparent;
    color: #0f3460;
    border: 1px solid #94a3b8;
    padding: 14px 24px;
    border-radius: 8px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
    font-family: 'Inter', sans-serif;
  }

  .full { width: 100%; }

  /* WORKFLOW / 3 LANGKAH */
  .steps {
    padding: 80px 48px;
    text-align: center;
    background: white;
  }

  .workflow-label {
    font-size: 12px;
    font-weight: 700;
    color: #64748b;
    letter-spacing: 2px;
    margin-bottom: 12px;
    display: block;
  }

  .steps h2 {
    font-size: 36px;
    font-weight: 800;
    margin-bottom: 48px;
  }

  .steps-row {
    display: flex;
    justify-content: space-between;
    gap: 24px;
    max-width: 900px;
    margin: 0 auto;
  }

  .step {
    flex: 1;
    text-align: center;
    padding: 0 20px;
  }

  .step-num {
    width: 56px;
    height: 56px;
    background: #e0f2fe;
    color: #00b4d8;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    font-weight: 800;
    margin: 0 auto 24px;
  }

  .step h3 {
    font-size: 18px;
    font-weight: 700;
    margin-bottom: 12px;
  }

  .step p {
    font-size: 14px;
    color: #64748b;
    line-height: 1.6;
  }

  /* PRIVASI (4 Kolom Sesuai Figma) */
  .features {
    padding: 80px 48px;
    text-align: center;
    background: #f0f9ff;
  }

  .features h2 {
    font-size: 36px;
    font-weight: 800;
    margin-bottom: 16px;
  }

  .feature-cards {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    max-width: 1200px;
    margin: 0 auto;
  }

  .feature-card {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    padding: 32px 24px;
    text-align: left;
  }

  .feature-icon-box {
    width: 44px;
    height: 44px;
    background: #e0f2fe;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 20px;
  }

  .feature-card h3 {
    font-size: 18px;
    font-weight: 700;
    margin-bottom: 12px;
  }

  .feature-card p {
    font-size: 13.5px;
    color: #64748b;
    line-height: 1.6;
  }

  /* FOOTER */
  footer {
    background: #0f3460;
    padding: 40px 48px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: white;
  }

  .footer-left { display: flex; align-items: center; gap: 12px; }
  .logo-img-small { height: 24px; width: auto; filter: brightness(0) invert(1); /* Membuat logo img jadi putih jika logonya gelap */ }
  .logo.small { color: white; font-size: 16px; letter-spacing: 1px; }

  .footer-links { display: flex; gap: 32px; }
  .footer-links a { color: #94a3b8; text-decoration: none; font-size: 14px; }
  .footer-links a:hover { color: white; }

  .footer-right span { color: #94a3b8; font-size: 13px; }
</style>