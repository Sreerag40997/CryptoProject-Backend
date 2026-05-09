# Crytinox — Crypto Exchange Simulation

A full-stack cryptocurrency exchange backend built in **Go + Fiber**, simulating real-world trading infrastructure with live Binance market data, advanced order types, and automated order execution.

> **Status:** In Progress &nbsp;|&nbsp; **Team:** 2 developers

---

## What This Is

Crytinox is not a price tracker or a simple CRUD app. It is a simulation of a real crypto exchange backend — with an order book, automated order matching engine, live market data pipeline, and a payment-verified wallet system.

---

## Key Features

**Order Engine**
- Supports 7 order types: Market, Limit, Stop-Market, Stop-Limit, Take-Profit, OCO, Trailing Stop
- Atomic database transactions on every order execution — no partial state
- Price Watcher Goroutine runs every 3 seconds, evaluating all pending orders concurrently against live Binance prices from Redis

**Real-Time Price Pipeline**
- Binance WebSocket → Redis → Fiber SSE broadcaster → React UI
- Live price updates pushed to all connected clients every 2 seconds
- Zero polling — fully event-driven

**Payment & Wallet**
- Razorpay integration with HMAC-SHA256 webhook signature verification
- Idempotency guards on wallet credit — funds added only on verified server-side event, never on frontend callback
- Prevents double-credit on duplicate webhook delivery

**Auth & Security**
- JWT-based authentication with RBAC
- Secure session management

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go (Golang) |
| Framework | Fiber |
| Database | PostgreSQL |
| Caching | Redis |
| Real-Time | WebSockets |
| Payments | Razorpay |
| Frontend | React |
| DevOps | Docker, Docker Compose |
| Architecture | Clean Architecture |

---

## Architecture Overview

```
Binance WebSocket
      │
      ▼
  Redis Cache  ◄──── Price Watcher Goroutine (every 3s)
      │                      │
      │              Evaluates pending orders
      │              Executes matches atomically
      ▼
Fiber Broadcaster
      │
      ▼
  React UI (live prices every 2s)
```

---

## Project Structure

```
CryptoProject-Backend/
├── cmd/                  # Entry point
├── internal/
│   ├── modules/           # modules all (with repo,service,handlers,routes)
├── pkg/                  # Shared utilities (JWT, HMAC, etc.)
├── docker-compose.yml
└── .env.example
```

---

## Getting Started

**Prerequisites:** Go 1.21+, PostgreSQL, Redis, Docker (optional)

```bash
# Clone the repo
git clone https://github.com/tibin-peter/CryptoProject-Backend.git
cd CryptoProject-Backend

# Copy environment variables
cp .env.example .env
# Fill in your PostgreSQL, Redis, Razorpay, and Binance API credentials

# Run with Docker
docker-compose up --build

# Or run locally
go run cmd/main.go
```

---

## Environment Variables

```env
DB_URL=postgres://user:password@localhost:5432/crytinox
REDIS_URL=localhost:6379
JWT_SECRET=your_jwt_secret
RAZORPAY_KEY_ID=your_key_id
RAZORPAY_KEY_SECRET=your_key_secret
BINANCE_WS_URL=wss://stream.binance.com:9443/ws
```

---

## Authors

1.
**Tibin Peter** — Backend Developer  
[GitHub](https://github.com/tibin-peter) · [LinkedIn](https://www.linkedin.com/in/tibin-peter-b9496a282/) · [Portfolio](https://tibin-peter-portfolio.vercel.app/)