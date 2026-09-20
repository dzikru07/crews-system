Ini khusus cheat sheet untuk orang yang baru masuk project.

# Crews Developer Commands

## Infrastructure

Start:

```bash
docker compose up -d

Stop:

docker compose stop

Start again:

docker compose start

Stop and remove containers:

docker compose down

View containers:

docker compose ps

View logs:

docker compose logs

PostgreSQL logs:

docker compose logs postgres

Redis logs:

docker compose logs redis
Go

Format:

go fmt ./...

Build:

go build ./...

Test:

go test ./...

Run API:

go run ./cmd/api

Run Worker:

go run ./cmd/worker
Git

Status:

git status

Diff:

git diff

Branches:

git branch

Create feature:

git checkout -b feature/name

Fetch:

git fetch origin
Docker Debugging

Redis:

docker exec crews-redis redis-cli ping

PostgreSQL:

docker exec crews-postgres pg_isready -U crews -d crews_api
API

Health:

curl http://localhost:8080/health

Ready:

curl http://localhost:8080/ready
PowerShell Environment
$env:DATABASE_URL="postgres://crews:crews_dev_password@localhost:5432/crews_api?sslmode=disable"

$env:REDIS_URL="redis://localhost:6379"

Check:

$env:DATABASE_URL
$env:REDIS_URL
Clean Build
go clean
go build ./...
Important

Do not paste production credentials into terminal history or source code.

Use the appropriate secret management mechanism for each environment.


---

## Struktur dokumentasi kita sekarang

Jadi folder `docs` menjadi:

```text
docs/
│
├── README.md
│
├── getting-started.md
│
├── architecture/
│   ├── overview.md
│   ├── identity.md
│   └── multi-tenancy.md
│
├── development/
│   ├── workflow.md
│   └── commands.md
│
├── database/
│   └── migrations.md
│
├── security/
│   ├── api-security.md
│   └── credential-management.md
│
└── operations/
    └── troubleshooting.md
