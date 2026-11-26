# Eneba Marketplace

An unofficial digital game marketplace built for the **Eneba Internship Assignment 2026**. The application recreates the Eneba shopping experience, from browsing and comparing seller offers to checkout and digital key delivery.

**Live demo:** [eneba-marketplace-psi.vercel.app](https://eneba-marketplace-psi.vercel.app/)

**Backend API:** [Health check](https://eneba-marketplace-api-ecja.onrender.com/api/health) · [Product catalog](https://eneba-marketplace-api-ecja.onrender.com/api/products)

> The application is hosted on free-tier services. The first request after inactivity may take around a minute while the backend starts up. All purchases are simulated, and digital keys are fictional and cannot be redeemed.

## Features

- Backend-driven homepage with promotional banners, best sellers and deals.
- Product catalog with platform, region and price filters, sorting and global search.
- Partial and typo-tolerant search using PostgreSQL `pg_trgm`.
- Product pages with multiple seller offers and price comparison.
- Shopping cart with quantity management and item removal.
- User registration, login and logout.
- Protected checkout with simulated payment and digital key allocation.
- Order history with purchased digital keys.
- Responsive, Eneba-inspired interface.

## Technology stack

| Layer | Technology | Hosting |
| --- | --- | --- |
| Frontend | React, TypeScript, Vite, React Router | Vercel |
| Backend | Go, REST API, pgx | Render |
| Database | PostgreSQL | Neon |

## Architecture

The backend is a modular monolith: a single Go application with separate feature packages. The React frontend communicates with the REST API, which handles business logic and database access.

```text
React + TypeScript (Vercel)
          | HTTPS / REST
          v
     Go API (Render)
          | SQL / TLS
          v
   PostgreSQL (Neon)
```

```text
eneba-marketplace/
├── backend/
│   ├── cmd/api/                 # Application entry point
│   ├── internal/
│   │   ├── config/              # Environment configuration
│   │   ├── database/            # PostgreSQL connection
│   │   ├── middleware/          # HTTP middleware
│   │   ├── httpx/               # HTTP response helpers
│   │   ├── server/              # API routing and integration tests
│   │   └── modules/
│   │       ├── auth/
│   │       ├── products/
│   │       ├── sellers/
│   │       ├── offers/
│   │       ├── homepage/
│   │       ├── search/
│   │       ├── cart/
│   │       ├── checkout/
│   │       └── orders/
│   ├── migrations/             # Database schema
│   └── seeds/                  # Demo catalog and inventory
├── frontend/
│   ├── public/
│   └── src/
│       ├── api/
│       ├── components/
│       ├── features/
│       └── styles/
└── docs/                      # Architecture, API and testing notes
```

Passwords are hashed with bcrypt, and authenticated routes use expiring session tokens. Prices and inventory are validated on the backend. Checkout uses a database transaction to allocate digital keys and prevent duplicate allocation. Monetary values are stored as integer euro cents.

## Local development

Docker is not required. The project can run locally on Windows using VS Code and PowerShell.

Prerequisites: Git, Go 1.23 or newer, Node.js 22 or newer, and a PostgreSQL database (Neon or local PostgreSQL).

### 1. Clone the repository

```powershell
git clone https://github.com/Mantyhas/eneba-marketplace.git
cd eneba-marketplace
```

### 2. Set up PostgreSQL

Create a separate development database in [Neon](https://console.neon.tech/), or use a local PostgreSQL instance. Execute the following scripts in order using the Neon SQL Editor or a PostgreSQL client:

1. `backend/migrations/001_init.sql`
2. `backend/seeds/001_demo.sql`

Use the same database and branch for both scripts and the backend connection. Do not use the production database for destructive tests.

### 3. Start the backend

```powershell
cd backend
Copy-Item .env.example .env
```

Set the following values in `backend/.env`:

```env
DATABASE_URL=postgresql://USER:PASSWORD@HOST/DATABASE?sslmode=require
PORT=8080
FRONTEND_URL=http://localhost:5173
```

Replace the example URL with your own database connection string. For Neon, preserve its TLS parameters. If using a pooled connection and encountering prepared-statement conflicts, use a direct connection or append `&default_query_exec_mode=exec` to the URL.

```powershell
go mod tidy
go run ./cmd/api
```

The API runs at `http://localhost:8080/api`. Check `/api/health` and `/api/products` to verify connectivity.

### 4. Start the frontend

Open a second terminal from the repository root:

```powershell
cd frontend
Copy-Item .env.example .env
npm install
npm run dev
```

Set `frontend/.env` to:

```env
VITE_API_URL=http://localhost:8080/api
```

Restart Vite if you changed the environment file after starting it. Open [http://localhost:5173](http://localhost:5173).

Never commit `.env` files or database credentials. Use fictional account details when testing the demo.

## Deployment

The application is deployed using Vercel, Render and Neon.

| Service | Configuration |
| --- | --- |
| Vercel | Root: `frontend`; preset: Vite; build: `npm run build`; output: `dist` |
| Render | Root: `backend`; build: `go build -o app ./cmd/api`; start: `./app` |
| Neon | Managed PostgreSQL database with migrations and demo data applied |

Vercel environment variable:

```env
VITE_API_URL=https://eneba-marketplace-api-ecja.onrender.com/api
```

Render environment variables:

- `DATABASE_URL`: private Neon connection string.
- `FRONTEND_URL`: `https://eneba-marketplace-psi.vercel.app` (exact origin for CORS).
- `PORT`: provided by Render.

The frontend includes a Vercel rewrite configuration for client-side routing. The Render root URL is not a website; use the `/api` endpoints instead.

## Testing

From `backend`, run `go test ./...`. From `frontend`, run `npm run build` to check TypeScript and build the application.

The main manual test flow is: browse or search for a game, compare sellers, add an offer to the cart, register or sign in, complete a simulated purchase, and view the digital key in order history.

See [architecture notes](docs/architecture.md), [API documentation](docs/api.md) and [test plan](docs/testing.md) for more detail.

## Scope and AI-assisted development

This is an internship demo, not an official Eneba service. The design is inspired by eneba.com. Sellers, prices and keys are demonstration data; no real payments, redeemable keys or seller administration are supported.

AI assistance was used during planning, implementation, debugging and deployment. The conversation history is provided separately as part of the assignment submission.
