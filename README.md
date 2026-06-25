# 🌑 Umbra
> **Secure · Anonymous · Ephemeral**
>
> End-to-end encrypted chat application — server never sees plaintext.

---

## 📖 About

Umbra is a web-based secure chat application built for the **Cybersecurity Project** course at President University Pekanbaru. It implements end-to-end encryption (E2EE) using industry-standard cryptography, ensuring that messages can only be read by the intended recipients — not even the server.

All communication is **anonymous** — no registration, no login, no identity stored anywhere.

---

## ✨ Features

### 2 Chat Modes

| Mode | Description | Auth Required |
|------|-------------|---------------|
| 👻 **Private Room** | Create or join a room with a unique code — messages deleted on exit | ❌ No |
| 🎲 **Random Match** | Matched anonymously with a random stranger | ❌ No |

### Security Features
- 🔒 **End-to-End Encryption** — all messages encrypted in browser
- 🛡️ **Anti-MITM** — ECDSA digital signature on every message
- 💨 **Ephemeral Chat** — messages destroyed when users leave room
- 🔑 **Zero Knowledge Server** — server only forwards ciphertext, never stores
- 🌐 **Private Key Never Leaves Browser** — exists in memory only
- 👤 **Fully Anonymous** — no account, no identity, no trace

---

## 🔐 Cryptography

Umbra uses a 3-layer cryptographic system:

```
ECDH (P-256)      →  Key Exchange (generate shared secret)
AES-256-GCM       →  Message Encryption + Integrity Check
ECDSA             →  Digital Signature (anti-MITM, identity verification)
```

### Encryption Flow

```
User A                        Server                        User B
  │                              │                              │
  │─── ECDH Key Exchange ───────►│◄──── ECDH Key Exchange ─────│
  │       (public keys only)     │        (public keys only)   │
  │                              │                              │
  │  generate shared secret      │              generate shared secret
  │  (never sent to server)      │              (never sent to server)
  │                              │                              │
  │── ECDSA sign ───────────────►│                              │
  │── AES-256-GCM encrypt ──────►│──── forward ciphertext ─────►│
  │                              │    (server can't read)       │
  │                              │                              │
  │                              │            AES-256-GCM decrypt
  │                              │            ECDSA verify ─────►✅
  │                              │                              │
  [user leaves room]             │                              │
  │                              │                              │
  └── room empty ───────────────►│── delete all room data ──────┘
                                 │   (nothing remains)
```

---

## 🛠️ Tech Stack

### Frontend
| Technology | Purpose | Why |
|------------|---------|-----|
| **SvelteKit** | UI Framework | Lightweight, high performance, native reactivity — better than Next.js for real-time apps |
| **TypeScript** | Language | Type safety, better DX |
| **WebCrypto API** | Cryptography | Native browser, hardware-accelerated, private key cannot leak |

### Backend
| Technology | Purpose | Why |
|------------|---------|-----|
| **Go (Golang)** | Backend Language | Goroutines handle thousands of connections with low memory |
| **Fiber** | Web Framework | Fastest HTTP framework for Go |
| **Gorilla WebSocket** | WebSocket | Best performance WebSocket for Go — used by Discord, Slack |

### Storage
| Technology | Purpose |
|------------|---------|
| **In-Memory Go** | Ephemeral room data — auto-destroyed when room is empty |
| **Redis** *(optional)* | Alternative ephemeral storage with TTL auto-expire |

> No permanent database — all data lives in memory and disappears automatically.

### Security Layer
| Technology | Purpose |
|------------|---------|
| **TLS 1.3** | Transport encryption |
| **WSS** | Secure WebSocket |
| **HTTPS** | Secure HTTP |

---

## 🏗️ Architecture

```
┌──────────────────────────────────────────────────┐
│              Client (SvelteKit)                   │
│                                                  │
│  ┌──────────────────┐  ┌───────────────────────┐ │
│  │  Private Room    │  │    Random Match       │ │
│  │  create/join     │  │    auto-paired        │ │
│  │  via room code   │  │    anonymous          │ │
│  └────────┬─────────┘  └──────────┬────────────┘ │
│           └──────────┬────────────┘              │
│                      │                           │
│        ┌─────────────▼────────────┐              │
│        │     WebCrypto Module     │              │
│        │  ECDH · AES-256-GCM      │              │
│        │  ECDSA · Key Management  │              │
│        │  private key: RAM only   │              │
│        └─────────────┬────────────┘              │
└──────────────────────┼──────────────────────────┘
                       │ WSS (TLS 1.3)
┌──────────────────────▼──────────────────────────┐
│           Go Backend (single binary)             │
│                                                 │
│         ┌───────────────────────────┐           │
│         │       WebSocket Hub       │           │
│         │     ciphertext relay      │           │
│         └──────────────┬────────────┘           │
│                        │                        │
│  ┌─────────────────────▼──────────────────────┐ │
│  │              Room Manager                  │ │
│  │   storage: In-Memory Go / Redis (optional) │ │
│  │   auto-destroy when all users disconnect   │ │
│  │   no DB · no auth · no trace               │ │
│  └────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘
```

---

## 🚀 Getting Started

### Prerequisites

```bash
# Required
go >= 1.21
node >= 18

# Optional (for Redis storage mode)
redis >= 7
```

### Installation

```bash
# Clone repository
git clone https://github.com/yourusername/umbra.git
cd umbra
```

#### Backend Setup

```bash
cd backend

# Install dependencies
go mod tidy

# Setup environment
cp .env.example .env

# Start server
go run cmd/server/main.go
```

#### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Setup environment
cp .env.example .env.local

# Start development server
npm run dev
```

### Environment Variables

#### Backend `.env`
```env
PORT=8080
STORAGE_MODE=memory    # "memory" or "redis"
REDIS_URL=redis://localhost:6379   # only if STORAGE_MODE=redis
```

#### Frontend `.env.local`
```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

---

## 📁 Project Structure

```
umbra/
├── frontend/                    # SvelteKit
│   ├── src/
│   │   ├── routes/
│   │   │   ├── +page.svelte    # Landing page (mode selection)
│   │   │   ├── room/[id]/      # Private Room
│   │   │   └── random/         # Random Match
│   │   └── lib/
│   │       ├── crypto.ts       # WebCrypto helpers (ECDH, AES, ECDSA)
│   │       └── socket.ts       # WebSocket client
│   └── package.json
│
├── backend/                     # Go
│   ├── cmd/
│   │   └── server/main.go      # Entry point
│   ├── internal/
│   │   ├── crypto/             # Key management
│   │   ├── socket/             # WebSocket hub
│   │   ├── room/               # Room manager
│   │   └── storage/            # In-Memory + Redis adapter
│   └── go.mod
│
├── docs/
│   ├── BRD_Umbra.docx          # Business Requirement Document
│   └── system-design.png       # Architecture diagram
│
└── README.md
```

---

## 🔒 Security Notes

- **Private keys** are generated in browser memory and never sent to the server
- **Messages** are encrypted before leaving the browser — server only sees ciphertext
- **Room data** is completely destroyed from memory when all users disconnect
- **ECDSA signatures** verify message authenticity and prevent MITM attacks
- **TLS 1.3** encrypts all transport layer communication
- **No database** — nothing is ever written to disk

---

## 👥 Team

| Name | Role |
|------|------|
| **Rizal** | Project Manager & Developer |
| **Anggie** | Frontend Developer |
| **Masya** | Backend & Cryptography Developer |
| **Geysa** | QA & Documentation |

**Supervisor:** Gilang Gumelar, S.Tr.Kom., M.Kom.

**Course:** Cybersecurity Project
**University:** President University Pekanbaru

---

## 📄 License

This project is developed for academic purposes at President University Pekanbaru.

---

<p align="center">
  <strong>Umbra</strong> — <em>Chat in the shadows.</em>
</p>
