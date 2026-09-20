# Crews API Platform

Crews API Platform is a centralized API platform designed to manage:

- Internal APIs
- Public APIs
- Private APIs
- B2B APIs
- External provider integrations
- Mobile application integrations
- Service-to-service communication
- Transactions
- Deposits
- Reconciliation
- Webhooks
- Credentials
- API access control
- Audit logs
- Observability

The system is designed to support multiple environments:

- Development
- Staging
- Production

---

## Architecture Principles

The main principles of Crews are:

1. Security first
2. Environment isolation
3. Multi-tenancy
4. Strong authentication
5. Fine-grained authorization
6. Credential encryption
7. Full auditability
8. Provider abstraction
9. Sync and async processing
10. Idempotent transactions
11. Reconciliation support
12. Observable systems
13. Failure isolation

---

## Repository Structure

```text
crews-system/
│
├── cmd/
│   ├── api/
│   │   └── main.go
│   │
│   └── worker/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── server/
│   │
│   ├── middleware/
│   │
│   ├── auth/
│   ├── identity/
│   ├── credential/
│   ├── authorization/
│   ├── api_registry/
│   ├── transaction/
│   ├── provider/
│   ├── webhook/
│   ├── reconciliation/
│   ├── idempotency/
│   ├── queue/
│   ├── audit/
│   └── observability/
│
├── pkg/
│   ├── crypto/
│   ├── httpclient/
│   ├── logger/
│   ├── response/
│   └── validator/
│
├── migrations/
├── configs/
├── deployments/
├── docs/
└── tests/

Main Components
API

Handles HTTP requests from:

Mobile applications
Web applications
Internal services
B2B partners
Worker

Handles asynchronous jobs such as:

Provider requests
Reconciliation
Webhook processing
Retry
Background transaction processing
PostgreSQL

Primary persistent database.

Used for:

Identity
Applications
Partners
Transactions
API registry
Audit
Reconciliation
Configuration metadata
Redis

Used for temporary/high-speed data such as:

Cache
Rate limiting
Distributed locks
Idempotency
Queue support
Short-lived authentication state

Redis must not become the primary source of truth for business data.

Environments

Crews supports:

development
staging
production

Credentials and resources must never be shared between environments.

Example:

Loket Development
    ↓
Development API
    ↓
Development Provider

Loket Production
    ↓
Production API
    ↓
Production Provider
Health Endpoints
GET /health

Checks whether the application process is alive.

Expected response:

{
  "status": "ok"
}
GET /ready

Checks whether required dependencies are available.

Expected response:

{
  "status": "ready",
  "database": "ok",
  "redis": "ok"
}

A ready response does not mean every external provider is available.

It only indicates that the application can accept requests.

Development

Start infrastructure:

docker compose up -d

Check containers:

docker compose ps

Stop infrastructure:

docker compose down

Stop without removing containers:

docker compose stop

Start again:

docker compose start
Go Commands

Format:

go fmt ./...

Build:

go build ./...

Test:

go test ./...

Run API:

go run ./cmd/api

Run worker:

go run ./cmd/worker
Development Environment

PowerShell example:

$env:DATABASE_URL="postgres://crews:crews_dev_password@localhost:5432/crews_api?sslmode=disable"
$env:REDIS_URL="redis://localhost:6379"

Run:

go run ./cmd/api
Important Security Rule

Never commit real credentials.

Never commit:

.env
API keys
API secrets
JWT secrets
provider passwords
private keys
certificates
production credentials

Only example values may be committed.

Use:

.env.example

for documentation.

Development Checklist

Before pushing code:

go fmt ./...
go build ./...
go test ./...

Check:

git status

Make sure no secret files are staged.

Architecture Documentation

See:

architecture/overview.md
architecture/identity.md
architecture/multi-tenancy.md
architecture/security.md
Security Documentation

See:

security/credential-management.md
security/api-security.md
security/incident-response.md
Operational Documentation

See:

operations/deployment.md
operations/monitoring.md
operations/backup.md

---

# 2. `docs/getting-started.md`

Ini yang dibaca developer baru.

```md
# Getting Started

## Requirements

Install:

- Git
- Go
- Docker
- Docker Compose

Verify:

```bash
git --version
go version
docker --version
docker compose version
Clone Repository
git clone <repository-url>
cd crews-system
Start Infrastructure
docker compose up -d

Check:

docker compose ps

Expected:

crews-postgres   Up (healthy)
crews-redis      Up (healthy)
Configure Development Environment

PowerShell:

$env:DATABASE_URL="postgres://crews:crews_dev_password@localhost:5432/crews_api?sslmode=disable"
$env:REDIS_URL="redis://localhost:6379"

Linux/macOS:

export DATABASE_URL="postgres://crews:crews_dev_password@localhost:5432/crews_api?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
Run API
go run ./cmd/api
Verify API

Health:

curl http://localhost:8080/health

Readiness:

curl http://localhost:8080/ready
Development Workflow
Create or checkout branch
Start Docker infrastructure
Configure environment
Run API
Make changes
Run formatter
Run tests
Build
Review git diff
Commit

Commands:

git checkout -b feature/example

go fmt ./...
go test ./...
go build ./...

git diff
git status
Troubleshooting
PostgreSQL unavailable

Check:

docker compose ps

Logs:

docker compose logs postgres
Redis unavailable

Check:

docker compose ps

Test:

docker exec crews-redis redis-cli ping

Expected:

PONG
PostgreSQL connection test
docker exec crews-postgres pg_isready -U crews -d crews_api

Expected:

accepting connections
API does not start

Run:

go run ./cmd/api

and inspect the first error.

Do not hide startup errors.

Important

Do not manually modify production configuration while debugging development.

Always identify the environment first.


---

# 3. `docs/architecture/overview.md`

```md
# Architecture Overview

Crews is a centralized API platform.

The system separates:

- Identity
- Application
- Client
- Credential
- Authorization
- API
- Provider
- Transaction
- Reconciliation

---

# High-Level Architecture

```text
                     Clients
                        │
        ┌───────────────┼────────────────┐
        │               │                │
      Mobile           Web              B2B
        │               │                │
        └───────────────┼────────────────┘
                        ▼
                  API Gateway
                        │
                Security Layer
                        │
          ┌─────────────┼─────────────┐
          │             │             │
      Authentication Authorization Rate Limit
          │             │             │
          └─────────────┼─────────────┘
                        ▼
                   API Registry
                        │
                        ▼
                   API Services
                        │
        ┌───────────────┼────────────────┐
        │               │                │
    PostgreSQL        Redis            Queue
        │                                │
        │                                ▼
        │                             Worker
        │                                │
        │                                ▼
        │                            Providers
        │
        ▼
    Audit / Logs
Synchronous Processing

Used when the caller requires an immediate response.

Client
  ↓
Crews API
  ↓
Provider
  ↓
Response
Asynchronous Processing

Used for long-running or retryable operations.

Client
  ↓
Crews API
  ↓
Transaction
  ↓
Queue
  ↓
Worker
  ↓
Provider
  ↓
Webhook / Result
  ↓
Transaction Update
Source of Truth

PostgreSQL is the source of truth for business state.

Redis is not the source of truth.

Queues are not the source of truth.

Provider responses may be stored for traceability and reconciliation.

Provider Isolation

External providers must be accessed through the provider layer.

Business services must not directly depend on provider-specific implementations.

Example:

Transaction Service
        │
        ▼
Provider Interface
        │
   ┌────┼────┐
   ▼    ▼    ▼
Provider A B  C

This allows providers to be replaced without changing the transaction domain.


---

# 4. `docs/architecture/identity.md`

```md
# Identity Architecture

Crews separates human identities from machine identities.

---

# Platform User

A human who manages the Crews platform.

Examples:

- Platform Admin
- Operations
- Finance
- Support

---

# Service Account

A machine identity.

Examples:

- loket-backend
- reconciliation-worker
- payment-worker

Service accounts must not use human passwords.

---

# Application

An application consuming Crews.

Examples:

- Loket Digital
- Atensi
- Ayo Gas

An application is not a user.

---

# Partner

A business entity consuming B2B services.

Examples:

- Travel partner
- PPOB partner
- Corporate partner

---

# Partner User

Human users belonging to a partner.

Used for dashboard access.

---

# API Client

A machine client used to call APIs.

Examples:

```text
loket-production
partner-abc-production
reconciliation-worker
Credential

Authentication material associated with a client.

Examples:

API Key
Client ID / Secret
HMAC Secret
OAuth credential
mTLS certificate

Credentials are not identities.

Relationship
Application
    │
    ▼
Client
    │
    ▼
Credential
    │
    ▼
Authorization
    │
    ▼
API
Mobile Users

Mobile users belong to the application.

Example:

Loket Digital
    │
    └── Mobile Users

Their authentication is handled by the application's backend.

Crews does not need to store their passwords.

Environment Isolation

Each environment has separate clients and credentials.

Development Client
Staging Client
Production Client

Production credentials must never be used in development.


---

# 5. `docs/architecture/multi-tenancy.md`

```md
# Multi-Tenancy

Crews is designed as a multi-tenant platform.

A tenant represents an isolated organizational boundary.

---

# Tenant Isolation

Tenant-specific resources must contain an ownership relationship to the tenant.

Examples:

- Applications
- Clients
- Credentials
- Users
- API access
- Transactions
- Provider configuration

---

# Example

```text
Tenant: Loket Digital
│
├── Application
│   └── Loket Backend
│
├── Clients
│   ├── Development
│   ├── Staging
│   └── Production
│
└── API Access

Another tenant:

Tenant: Atensi
│
├── Application
│
├── Clients
│
└── API Access
Isolation Rule

A request authenticated for Tenant A must not be able to access Tenant B resources unless explicitly authorized.

Tenant filtering must be enforced at the service/repository layer.

Never rely only on frontend filtering.

Environment Isolation

Tenant resources are additionally separated by environment.

Tenant
  │
  ├── Development
  ├── Staging
  └── Production

---

# 6. `docs/security/credential-management.md`

Ini nanti akan menjadi **dokumen sangat penting**.

```md
# Credential Management

Credentials are sensitive security material.

Examples:

- API keys
- API secrets
- Client secrets
- Provider passwords
- Provider API keys
- JWT signing secrets
- Encryption keys
- Private keys

---

# Rules

Credentials must:

1. Never be committed to Git
2. Never be logged
3. Never be returned through normal API responses
4. Never be stored as plaintext when encryption is required
5. Be isolated by environment
6. Be revocable
7. Be auditable
8. Have controlled access

---

# Storage

Sensitive credentials should be encrypted at rest.

Concept:

```text
Plain Credential
      │
      ▼
Encryption Service
      │
      ▼
Encrypted Credential
      │
      ▼
PostgreSQL

Decryption should happen only when required.

Logging

Never log:

Authorization
API-Key
Client-Secret
Password
Private-Key

Use masking:

sk_live_************abcd

instead of:

sk_live_12345678901234567890
Credential Lifecycle
Created
   ↓
Active
   ↓
Rotating
   ↓
Revoked

Credentials must have status and audit history.

Environment

Development, staging and production credentials are completely separate.

Never copy production credentials into development.


---

# 7. `docs/development/workflow.md`

```md
# Development Workflow

## Branches

Recommended:

```text
main
develop
feature/*
fix/*
hotfix/*

Feature example:

git checkout -b feature/api-registry
Before Commit

Run:

go fmt ./...
go test ./...
go build ./...

Check:

git status
git diff
Secret Check

Before commit verify that no sensitive files are staged:

git status

Never commit:

.env
.env.production
*.pem
*.key
production credentials
provider credentials
Commit

Example:

git add .
git commit -m "feat: add api registry foundation"
Pull Latest Changes
git fetch origin
git pull --rebase
Important Rule

Database migrations and security changes must be reviewed carefully before merging.


---

# 8. `docs/database/migrations.md`

```md
# Database Migrations

Database schema changes must use migrations.

Never modify production tables manually unless there is an approved emergency procedure.

---

# Migration Naming

Example:

```text
000001_create_tenants.up.sql
000001_create_tenants.down.sql

000002_create_applications.up.sql
000002_create_applications.down.sql
Migration Rules

Every migration should have:

.up.sql
.down.sql

Example:

000001_create_tenants.up.sql
000001_create_tenants.down.sql
Migration Principles
Small changes
Reversible when possible
Tested locally
Tested in staging
Reviewed before production
Never silently destroy data
Production

Production migrations must be executed using a controlled deployment process.

Never run experimental migrations directly against production.


---

# 9. `docs/operations/troubleshooting.md`

```md
# Troubleshooting

## Check Containers

```bash
docker compose ps
PostgreSQL

Logs:

docker compose logs postgres

Connection test:

docker exec crews-postgres pg_isready -U crews -d crews_api
Redis

Logs:

docker compose logs redis

Connection test:

docker exec crews-redis redis-cli ping

Expected:

PONG
API

Run:

go run ./cmd/api
Health
curl http://localhost:8080/health
Readiness
curl http://localhost:8080/ready

Expected:

{
  "status": "ready",
  "database": "ok",
  "redis": "ok"
}
API Returns 500

First check application logs.

Do not immediately modify database or Redis configuration.

Check:

Application panic
Database connection
Redis connection
Configuration
Middleware
Request ID
Redis Connection Error

Correct URL:

redis://localhost:6379

Not:

localhost

when using URL parsing.

Production may use:

rediss://...

---

# 10. `docs/security/api-security.md`

```md
# API Security

All external API requests must pass through the security layer.

---

# Request Flow

```text
Request
  ↓
TLS
  ↓
Request ID
  ↓
IP Policy
  ↓
Authentication
  ↓
Authorization
  ↓
API Access
  ↓
Rate Limit
  ↓
Validation
  ↓
Handler
Authentication

Supported mechanisms may include:

API Key
Client ID / Secret
HMAC
OAuth
JWT
mTLS

Authentication identifies the caller.

Authentication does not automatically grant permission.

Authorization

Authorization determines what the authenticated client is allowed to do.

Example:

transaction:create
transaction:read
transaction:cancel
IP Whitelist

IP restrictions may be configured for B2B clients.

Example:

Partner ABC
    │
    ├── 103.x.x.x
    └── 103.x.x.x

IP whitelist must not replace authentication.

Both should be evaluated when required.

Rate Limiting

Rate limits may be applied per:

Tenant
Client
API
Endpoint
IP
Idempotency

Financial or transactional endpoints should support idempotency where appropriate.

Example:

Idempotency-Key: abc-123

The same request must not accidentally create multiple transactions.

Audit

Security-sensitive actions must be auditable.

Examples:

Credential creation
Credential rotation
Credential revocation
Permission changes
IP whitelist changes
API access changes

---

# 11. `.env.example`

Ini juga **wajib kita simpan di repository**.

```env
# Application
APP_NAME=crews-system
APP_ENV=development
APP_PORT=8080

# Database
DATABASE_URL=postgres://crews:crews_dev_password@localhost:5432/crews_api?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# Logging
LOG_LEVEL=debug

# Security
JWT_SECRET=change-me
ENCRYPTION_KEY=change-me

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000

Catatan: nilai di atas hanya contoh development. Jangan gunakan nilai tersebut untuk production.