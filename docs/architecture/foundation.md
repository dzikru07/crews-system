# Crews API Platform

## Technical Architecture & Integration Guide

**Version:** 0.1
**Status:** Foundation / Development
**Platform:** Crews API Platform
**Repository:** `github.com/dzikru07/crews-system`

---

# 1. Tujuan Sistem

Crews API Platform adalah sebuah **API management dan integration platform internal** yang menjadi pusat pengelolaan:

* aplikasi internal
* aplikasi mobile
* API public
* API private
* API B2B
* credential
* user dan identity
* authentication
* authorization
* environment
* provider / aggregator
* transaction
* webhook
* reconciliation
* audit log
* rate limit
* IP whitelist
* API access control

Tujuan utamanya adalah agar aplikasi seperti:

* Loket Digital
* aplikasi mobile lainnya
* sistem B2B
* dashboard internal
* partner/mitra
* service internal

tidak perlu mengelola credential, authentication, authorization, logging, dan integrasi API secara terpisah.

Crews menjadi **central API platform**.

---

# 2. Prinsip Utama

Arsitektur Crews menggunakan beberapa prinsip:

### 2.1 Separation of Concerns

Crews tidak menganggap semua user sebagai satu jenis user.

Ada perbedaan antara:

* manusia
* machine/service
* aplikasi
* client
* partner
* credential

Hal ini penting karena kebutuhan authentication dan authorization masing-masing berbeda.

---

### 2.2 Environment Isolation

Setiap aplikasi dapat memiliki environment:

```text
development
staging
production
```

Credential dan konfigurasi antar environment tidak boleh dianggap sama.

Contoh:

```text
Loket Digital
│
├── Development
│   └── Mobile Client
│       └── Development Credential
│
├── Staging
│   └── Mobile Client
│       └── Staging Credential
│
└── Production
    └── Mobile Client
        └── Production Credential
```

Dengan demikian credential development tidak boleh digunakan untuk production.

---

### 2.3 Credential Tidak Sama Dengan User

Credential adalah identitas teknis untuk melakukan komunikasi antar sistem.

Contoh:

```text
Loket Digital Backend
        │
        │ client_id + client_secret
        ▼
Crews API Platform
```

Sedangkan user aplikasi adalah user milik aplikasi tersebut.

Contoh:

```text
Loket Digital
│
├── User A
├── User B
├── User C
└── User D
```

User aplikasi tersebut **tidak harus menjadi user langsung di database Crews**.

---

# 3. Konsep Tenant

`Tenant` adalah pemilik atau boundary utama dari resource di Crews.

Contoh:

```text
Tenant
└── Loket Digital
```

Database:

```text
tenants
```

Contoh data:

```text
id:   0a17fa46-3421-42c8-a150-3ec8e7daa91d
name: Loket Digital
slug: loket-digital
type: organization
```

Tenant digunakan sebagai boundary isolation.

Artinya resource milik:

```text
Tenant A
```

tidak boleh sembarangan diakses oleh:

```text
Tenant B
```

---

# 4. Application

Application merepresentasikan sebuah sistem/aplikasi yang menggunakan Crews.

Contoh:

```text
Tenant
└── Loket Digital

Application
└── Loket Digital API
```

Database:

```text
applications
```

Contoh:

```text
name:
Loket Digital API

slug:
loket-digital-api
```

Satu tenant dapat mempunyai banyak application.

Contoh:

```text
Loket Digital
│
├── Loket Digital API
├── Loket Digital Mobile
├── Loket Digital Admin
├── Loket Digital B2B
└── Loket Digital Partner API
```

Namun pada tahap sekarang kita menggunakan `application` sebagai logical container, bukan sebagai user.

---

# 5. Environment

Setiap application dapat memiliki beberapa environment.

Contoh:

```text
Loket Digital API
│
├── Development
├── Staging
└── Production
```

Database:

```text
environments
```

Jenis environment yang direncanakan:

```text
development
staging
production
```

Environment berfungsi sebagai isolation boundary.

Contoh credential:

```text
Development
client_id = crw_dev_xxxxx

Staging
client_id = crw_stg_xxxxx

Production
client_id = crw_live_xxxxx
```

Credential tersebut tidak boleh digunakan silang environment.

---

# 6. Client

Client merepresentasikan consumer teknis dalam suatu environment.

Contoh:

```text
Loket Digital
└── Loket Digital API
    └── Production
        ├── Mobile
        ├── Backend
        ├── Admin
        └── B2B
```

Database:

```text
clients
```

Contoh:

```text
name:
Loket Digital Mobile

slug:
loket-digital-mobile

type:
mobile
```

Client adalah salah satu komponen penting untuk menentukan:

* siapa consumer-nya
* environment yang digunakan
* identity yang digunakan
* credential yang digunakan
* API apa yang boleh diakses

---

# 7. Identity

Identity merupakan representasi aktor yang melakukan autentikasi.

Database:

```text
identities
```

Saat ini identity memiliki dua tipe utama:

```text
human
machine
```

## 7.1 Human Identity

Digunakan untuk manusia.

Contoh:

```text
Admin
Operator
Developer
Finance
Support
```

Human identity nantinya dapat memiliki:

```text
user_accounts
```

Contoh:

```text
Identity
└── Human
    └── User Account
        ├── email
        ├── username
        └── password
```

---

# 8. Machine Identity

Machine identity digunakan untuk sistem/service.

Contoh:

```text
Loket Digital Backend
```

atau:

```text
B2B Partner ABC
```

Contoh:

```text
identity:
Loket Digital Mobile

type:
machine
```

Machine identity **tidak mempunyai password manusia**.

Machine identity menggunakan machine credential.

---

# 9. User Account

`user_accounts` digunakan untuk akun manusia.

Strukturnya:

```text
Identity
    │
    └── User Account
```

Database:

```text
user_accounts
```

Field utama:

```text
identity_id
email
username
password_hash
email_verified_at
last_login_at
```

Password tidak pernah disimpan plaintext.

Contoh:

```text
password:
Rahasia123

```

tidak boleh disimpan sebagai:

```text
Rahasia123
```

Tetapi menjadi:

```text
password_hash
```

menggunakan password hashing yang aman.

---

# 10. Machine Credential

Machine credential digunakan oleh aplikasi/service untuk melakukan authentication.

Database:

```text
machine_credentials
```

Credential terdiri dari:

```text
client_id
client_secret
```

Namun database hanya menyimpan:

```text
client_secret_hash
```

Bukan plaintext secret.

Contoh response ketika credential dibuat:

```json
{
  "id": "...",
  "client_id": "crw_dev_xxxxx",
  "client_secret": "crs_dev_xxxxx",
  "status": "active"
}
```

`client_secret` hanya boleh diberikan ketika credential dibuat.

Setelah itu secret tidak boleh dapat dibaca kembali dari database.

---

# 11. Credential Flow

Alur pembuatan credential saat ini:

```text
Client
   │
   ▼
Identity
   │
   ▼
Credential Creation
   │
   ├── Generate client_id
   │
   ├── Generate client_secret
   │
   ├── Hash client_secret
   │
   └── Store hash
```

Database:

```text
machine_credentials
```

menyimpan:

```text
client_id
client_secret_hash
```

bukan:

```text
client_secret
```

---

# 12. Prefix Credential

Credential dibuat berdasarkan environment.

Contoh:

### Development

```text
crw_dev_xxxxx
crs_dev_xxxxx
```

### Staging

```text
crw_stg_xxxxx
crs_stg_xxxxx
```

### Production

```text
crw_live_xxxxx
crs_live_xxxxx
```

Prefix ini membantu:

* identifikasi environment
* debugging
* audit
* menghindari kesalahan penggunaan credential

Prefix **bukan security mechanism utama**.

Security tetap bergantung pada authentication dan authorization.

---

# 13. Hubungan Data

Relasi saat ini:

```text
Tenant
│
├── Applications
│   │
│   └── Environments
│       │
│       └── Clients
│           │
│           └── Identity
│
└── Identities
    │
    ├── User Accounts
    │
    └── Machine Credentials
```

Secara database:

```text
tenants
   │
   ├── applications
   │       │
   │       └── environments
   │               │
   │               └── clients
   │                       │
   │                       └── identities
   │
   └── identities
           │
           ├── user_accounts
           │
           └── machine_credentials
```

---

# 14. Contoh Arsitektur Loket Digital

Misalnya kita mempunyai:

```text
Tenant:
Loket Digital
```

Application:

```text
Loket Digital API
```

Environment:

```text
Development
Staging
Production
```

Client:

```text
Loket Digital Mobile
Loket Digital Backend
Loket Digital B2B
```

Maka:

```text
Loket Digital
│
└── Loket Digital API
    │
    ├── Development
    │   │
    │   ├── Mobile Client
    │   └── Backend Client
    │
    ├── Staging
    │   │
    │   ├── Mobile Client
    │   └── Backend Client
    │
    └── Production
        │
        ├── Mobile Client
        ├── Backend Client
        └── B2B Client
```

---

# 15. Mobile Application

Mobile application membutuhkan perlakuan khusus.

Jangan menganggap:

```text
Flutter App
```

sama dengan:

```text
User Account
```

Flutter application adalah consumer.

Sedangkan user aplikasi adalah pengguna akhir.

Contoh:

```text
Loket Digital Mobile
│
├── Mobile Client
│
└── Users
    ├── User A
    ├── User B
    └── User C
```

User A/B/C dapat disimpan di database aplikasi Loket Digital.

Crews tidak perlu mengetahui seluruh detail user aplikasi tersebut.

---

# 16. Backend Application

Backend Loket Digital dapat menggunakan machine credential.

Contoh:

```text
Loket Digital Backend
        │
        │ client_id
        │ client_secret
        ▼
   Crews Authentication
        │
        ▼
   Access Token
        │
        ▼
     API Request
```

Backend dapat menggunakan credential tanpa mengekspos secret kepada Flutter.

---

# 17. Hal Penting Untuk Mobile

Untuk aplikasi mobile, **jangan menaruh client_secret permanen di dalam APK/IPA**.

Alasannya:

APK/IPA dapat dibongkar atau dianalisis.

Jadi pola yang lebih aman adalah:

```text
Flutter
   │
   ▼
Loket Digital Backend
   │
   ▼
Crews
```

bukan:

```text
Flutter
   │
   │ client_secret
   ▼
Crews
```

Credential sensitif idealnya berada di backend.

---

# 18. B2B Partner

Crews juga akan digunakan untuk partner B2B.

Contoh:

```text
Crews
│
├── Partner A
├── Partner B
└── Partner C
```

Partner dapat diberikan credential masing-masing.

Contoh:

```text
Partner A
client_id = crw_live_partner_a_xxx

Partner B
client_id = crw_live_partner_b_xxx
```

Dengan demikian kita dapat melakukan:

* revoke credential
* disable partner
* rate limiting
* IP whitelist
* API permission
* audit
* monitoring
* transaction tracking

secara individual.

---

# 19. Partner Tidak Boleh Berbagi Credential

Credential harus bersifat individual.

Tidak boleh:

```text
Partner A
    │
    └── Credential bersama
          │
          ├── Partner B
          └── Partner C
```

Yang benar:

```text
Partner A → Credential A

Partner B → Credential B

Partner C → Credential C
```

Hal ini sangat penting untuk audit dan security.

---

# 20. API Registry

Tahap berikutnya adalah membuat API Registry.

Tujuannya adalah Crews mengetahui API apa saja yang tersedia.

Contoh:

```text
API Registry
│
├── PPOB API
├── Flight API
├── Train API
├── Hotel API
├── Payment API
└── User API
```

Setiap API dapat memiliki:

```text
API
├── endpoint
├── method
├── version
├── environment
├── authentication
├── permission
├── rate limit
└── status
```

---

# 21. Public API dan Private API

Crews akan mendukung dua kategori besar.

## Public API

API yang dapat digunakan consumer melalui API gateway.

Contoh:

```text
POST /v1/products
```

## Private API

API internal/provider yang tidak boleh langsung diakses oleh public.

Contoh:

```text
Crews
  │
  ▼
Provider API
```

Credential provider disimpan di Crews.

Consumer tidak mengetahui credential provider.

---

# 22. Provider Integration

Arsitektur yang direncanakan:

```text
Client
   │
   ▼
Crews
   │
   ├── Authentication
   ├── Authorization
   ├── Rate Limit
   ├── IP Whitelist
   ├── Idempotency
   ├── Audit
   │
   ▼
Provider Adapter
   │
   ▼
Aggregator / Biller
```

Provider credential tidak diberikan kepada consumer.

Contoh:

```text
Partner
   │
   ▼
Crews
   │
   ▼
Aggregator A
```

Partner hanya mengetahui credential Crews.

---

# 23. Provider Credential

Provider credential akan dipisahkan dari credential consumer.

Contoh:

```text
Consumer Credential
    │
    ▼
Crews
    │
    ▼
Provider Credential
```

Consumer:

```text
client_id
client_secret
```

Provider:

```text
api_key
secret
username
password
signature_key
private_key
```

Credential provider akan disimpan secara aman dan dienkripsi.

---

# 24. Authentication

Authentication menjawab:

> "Siapa yang melakukan request?"

Contoh:

```text
client_id
client_secret
```

diverifikasi oleh Crews.

Jika valid:

```text
Authenticated
```

Jika tidak:

```text
401 Unauthorized
```

---

# 25. Authorization

Authorization menjawab:

> "Apakah identity tersebut boleh melakukan operasi ini?"

Contoh:

```text
Partner A
```

boleh:

```text
POST /v1/tickets/order
```

tetapi tidak boleh:

```text
GET /admin/users
```

Authorization akan dibangun menggunakan:

* roles
* permissions
* scopes
* API policies

---

# 26. IP Whitelist

Untuk B2B, kita akan mendukung IP restriction.

Contoh:

```text
Partner A

Allowed IP:
103.x.x.x
103.x.x.x
```

Request:

```text
103.x.x.x
```

→ allowed

Request:

```text
8.x.x.x
```

→ rejected

Contoh response:

```json
{
  "error": "ip_not_allowed"
}
```

IP whitelist sebaiknya menjadi salah satu layer security, bukan satu-satunya authentication.

---

# 27. Rate Limiting

Setiap client/partner dapat mempunyai rate limit.

Contoh:

```text
Partner A
100 requests / second
```

atau:

```text
Partner B
20 requests / second
```

Rate limit dapat diterapkan berdasarkan:

```text
client
identity
API
IP
tenant
```

Redis akan digunakan sebagai salah satu komponen untuk kebutuhan rate limiting dan distributed state.

---

# 28. Idempotency

Untuk transaksi penting, terutama B2B, request harus mendukung idempotency.

Contoh:

```http
Idempotency-Key: ORDER-20260921-00001
```

Jika partner mengirim request dua kali:

```text
Request 1 → processed
Request 2 → same result
```

bukan:

```text
Request 1 → create transaction
Request 2 → create transaction again
```

Hal ini sangat penting untuk:

* payment
* booking
* PPOB
* ticketing
* order
* topup

---

# 29. Transaction

Setiap transaksi nantinya mempunyai lifecycle.

Contoh:

```text
RECEIVED
   │
   ▼
PROCESSING
   │
   ├── SUCCESS
   │
   ├── FAILED
   │
   └── UNKNOWN
```

Status `UNKNOWN` penting ketika provider tidak memberikan response yang jelas.

Contoh:

```text
Crews → Provider
        │
        ▼
     timeout
```

Belum tentu transaksi gagal.

Karena itu diperlukan reconciliation.

---

# 30. Sync dan Async

Platform dirancang mendukung dua model.

## Synchronous

```text
Client
  │
  ▼
Crews
  │
  ▼
Provider
  │
  ▼
Response
  │
  ▼
Client
```

Cocok untuk:

* inquiry
* price
* availability
* simple validation

---

## Asynchronous

```text
Client
  │
  ▼
Crews
  │
  ▼
Queue
  │
  ▼
Worker
  │
  ▼
Provider
  │
  ▼
Webhook / Result
```

Cocok untuk:

* transaksi
* retry
* long-running operation
* reconciliation
* batch process

Redis nantinya digunakan untuk queue/cache/rate limiting sesuai kebutuhan arsitektur.

---

# 31. Raw Provider Response

Salah satu prinsip penting platform:

> Response asli dari provider harus dapat disimpan secara terpisah.

Contoh:

```text
Provider Response
       │
       ├── Raw Response
       │
       └── Normalized Response
```

Raw response berguna untuk:

* debugging
* audit
* dispute
* reconciliation
* investigasi provider
* debugging mapping

Contoh:

```text
raw_provider_responses
```

tidak boleh dicampur dengan business response yang sudah dinormalisasi.

---

# 32. Logging

Logging dipisahkan dari transaction data.

Contoh:

```text
Application Log
Audit Log
Transaction Log
Provider Log
Security Log
```

Tujuannya agar ketika terjadi error kita dapat mengetahui:

```text
Request masuk
     ↓
Authentication
     ↓
Authorization
     ↓
Transaction
     ↓
Provider request
     ↓
Provider response
     ↓
Final response
```

---

# 33. Audit Log

Audit log digunakan untuk mengetahui:

> siapa melakukan apa, kapan, dan terhadap resource apa.

Contoh:

```text
Actor:
Partner A

Action:
CREATE_TRANSACTION

Resource:
transaction/xxxx

Timestamp:
2026-09-21T...
```

Audit log berbeda dari application log.

---

# 34. Security Architecture

Security akan dibuat berlapis.

```text
                    Internet
                       │
                       ▼
                Reverse Proxy
                       │
                       ▼
                Crews Gateway
                       │
              ┌────────┴────────┐
              ▼                 ▼
        Authentication     IP Filter
              │                 │
              └────────┬────────┘
                       ▼
                 Authorization
                       │
                       ▼
                  Rate Limit
                       │
                       ▼
                  Idempotency
                       │
                       ▼
                   Business
                       │
                       ▼
                    Provider
```

Tidak ada satu security mechanism yang dianggap cukup.

---

# 35. Secret Management

Credential sensitif tidak boleh:

* commit ke Git
* ditulis di source code
* disimpan plaintext
* dikirim melalui log
* dikirim ke client mobile
* dimasukkan ke Postman collection production

Contoh yang harus dihindari:

```go
password := "mypassword123"
```

atau:

```env
PROVIDER_SECRET=my-secret
```

di repository.

Development boleh menggunakan `.env` lokal, tetapi production harus menggunakan secret management yang sesuai.

---

# 36. Environment Configuration

Konfigurasi dipisahkan:

```text
development
staging
production
```

Contoh:

```text
DATABASE_URL
REDIS_URL
JWT_SECRET
ENCRYPTION_KEY
```

Environment production harus mempunyai credential dan secret yang berbeda dari development.

---

# 37. Current Database Schema

Saat ini database sudah mempunyai migration:

```text
000001_create_tenants
000002_create_applications
000003_create_environments
000004_create_clients
000005_create_identities
000006_create_user_accounts
000007_create_machine_credentials
000008_add_identity_to_clients
```

Urutan ini merupakan bagian penting dari evolution database.

Migration tidak boleh diubah sembarangan setelah digunakan bersama.

Migration berikutnya harus menggunakan nomor:

```text
000009_...
```

---

# 38. Current Tables

Saat ini tabel utama:

```text
tenants
applications
environments
clients
identities
user_accounts
machine_credentials
schema_migrations
```

---

# 39. Current Foundation Flow

Flow yang sudah berhasil dibuat:

```text
Tenant
  │
  ▼
Application
  │
  ▼
Environment
  │
  ▼
Client
  │
  ▼
Identity
  │
  ▼
Machine Credential
```

Contoh nyata yang sudah dibuat:

```text
Loket Digital
│
└── Loket Digital API
    │
    └── Production
        │
        └── Loket Digital Mobile
            │
            └── Loket Digital Mobile Identity
                │
                └── Machine Credential
```

---

# 40. Endpoint Foundation Saat Ini

## Health

```http
GET /health
```

Digunakan untuk mengetahui apakah application process hidup.

Response:

```json
{
  "status": "ok"
}
```

---

## Readiness

```http
GET /ready
```

Memeriksa dependency:

```text
PostgreSQL
Redis
```

Contoh:

```json
{
  "status": "ready",
  "database": "ok",
  "redis": "ok"
}
```

---

## Tenant

```http
POST /v1/tenants
GET  /v1/tenants/:tenant_id
```

---

## Application

```http
POST /v1/tenants/:tenant_id/applications

GET /v1/tenants/:tenant_id/applications/:application_id
```

---

## Environment

```http
POST /v1/tenants/:tenant_id/applications/:application_id/environments
```

---

## Client

```http
POST /v1/tenants/:tenant_id/applications/:application_id/environments/:environment_id/clients
```

---

## Identity

```http
POST /v1/tenants/:tenant_id/identities

GET /v1/tenants/:tenant_id/identities/:identity_id
```

---

## Machine Credential

```http
POST /v1/tenants/:tenant_id/identities/:identity_id/credentials
```

---

# 41. Contoh Request Credential

```http
POST /v1/tenants/{tenant_id}/identities/{identity_id}/credentials

X-Crews-Environment: development

Content-Type: application/json
```

Body:

```json
{
  "expires_at": null
}
```

Response akan memberikan credential baru.

Contoh:

```json
{
  "id": "...",
  "client_id": "crw_dev_xxxxx",
  "client_secret": "crs_dev_xxxxx",
  "status": "active"
}
```

Secret harus diamankan oleh consumer.

---

# 42. Linking Client dan Identity

Saat ini `clients.identity_id` sudah tersedia.

Relasinya:

```text
Client
   │
   └── identity_id
          │
          ▼
       Identity
```

Contoh:

```sql
UPDATE clients
SET identity_id = 'IDENTITY_ID'
WHERE id = 'CLIENT_ID';
```

Namun ini merupakan tahap sementara.

Ke depan sebaiknya linking tersebut dilakukan melalui API/service, bukan SQL manual.

---

# 43. Apa Yang Belum Final

Foundation saat ini belum dianggap sebagai final security architecture.

Komponen berikut masih akan dibangun:

```text
Authentication
Authorization
Access Token
JWT / token strategy
API Registry
API Policy
Scope
Permission
Role
IP Whitelist
Rate Limit
Idempotency
Provider Registry
Provider Credential
Encryption at Rest
Webhook
Transaction
Queue Worker
Retry
Reconciliation
Audit
Observability
```

Karena itu credential yang saat ini dibuat belum boleh dianggap sebagai mekanisme authentication production final.

---

# 44. Roadmap Arsitektur

Tahapan berikutnya:

```text
FOUNDATION
    │
    ├── Tenant
    ├── Application
    ├── Environment
    ├── Client
    ├── Identity
    └── Credential
            │
            ▼
AUTHENTICATION
    │
    ├── Client Authentication
    ├── Access Token
    ├── Token Validation
    └── Credential Rotation
            │
            ▼
AUTHORIZATION
    │
    ├── Role
    ├── Permission
    ├── Scope
    └── API Policy
            │
            ▼
API MANAGEMENT
    │
    ├── API Registry
    ├── Versioning
    ├── Rate Limit
    └── IP Whitelist
            │
            ▼
INTEGRATION
    │
    ├── Provider
    ├── Provider Credential
    ├── Sync
    ├── Async
    └── Webhook
            │
            ▼
TRANSACTION
    │
    ├── Idempotency
    ├── Transaction
    ├── Raw Response
    ├── Retry
    └── Reconciliation
            │
            ▼
OBSERVABILITY
    │
    ├── Audit
    ├── Logs
    ├── Metrics
    └── Monitoring
```

---

# 45. Aturan Untuk Developer

Developer yang bekerja pada repository harus memahami beberapa aturan:

### Database

Jangan mengubah database production secara manual.

Gunakan migration.

```text
000009_xxx.up.sql
000009_xxx.down.sql
```

---

### Credential

Jangan log:

```text
client_secret
password
provider_secret
private_key
```

---

### Environment

Jangan menggunakan credential:

```text
development
```

untuk:

```text
production
```

---

### Tenant Isolation

Semua query resource tenant harus memastikan:

```text
resource.tenant_id == authenticated_tenant
```

Jangan hanya menggunakan:

```sql
WHERE id = $1
```

jika resource bersifat tenant-scoped.

Gunakan:

```sql
WHERE id = $1
AND tenant_id = $2
```

---

### Error

Jangan mengembalikan error internal kepada consumer.

Hindari:

```json
{
  "error": "pq: duplicate key value..."
}
```

Gunakan error yang terkontrol:

```json
{
  "error": "resource_already_exists"
}
```

Detail internal masuk ke log.

---

# 46. Aturan Untuk Partner

Partner B2B akan menerima:

```text
client_id
client_secret
```

atau mekanisme authentication yang disepakati.

Partner tidak diberikan:

```text
provider credential
database credential
Redis credential
internal API credential
encryption key
JWT signing key
```

Partner hanya berinteraksi dengan public B2B API yang telah diizinkan.

---

# 47. Contoh Arsitektur Final B2B

```text
                    PARTNER
                       │
                       │ HTTPS
                       ▼
                ┌──────────────┐
                │ Crews Gateway│
                └──────┬───────┘
                       │
             ┌─────────┼─────────┐
             │         │         │
             ▼         ▼         ▼
        Auth        IP ACL    Rate Limit
             │         │         │
             └─────────┼─────────┘
                       ▼
                 Authorization
                       │
                       ▼
                  API Registry
                       │
                       ▼
                  Transaction
                       │
                       ▼
                Provider Adapter
                       │
                       ▼
                 Aggregator
                       │
                       ▼
                  Provider
```

Partner hanya melihat:

```text
Crews API
```

sedangkan kompleksitas provider berada di belakang Crews.

---

# 48. Contoh Arsitektur Mobile

```text
Flutter
   │
   │ HTTPS
   ▼
Loket Digital Backend
   │
   │ Machine Credential
   ▼
Crews
   │
   ▼
Provider / Aggregator
```

User:

```text
Flutter
   │
   ▼
Loket Digital User System
```

Dengan demikian:

```text
Crews Identity
```

tidak harus menjadi:

```text
Loket Digital User
```

Ini merupakan separation yang disengaja.

---

# 49. Kesimpulan

Crews bukan sekadar backend API.

Crews dirancang sebagai:

> **Central API, Identity, Credential, Security, Integration, and Transaction Platform.**

Model dasarnya:

```text
Tenant
  ↓
Application
  ↓
Environment
  ↓
Client
  ↓
Identity
  ↓
Credential
  ↓
Authentication
  ↓
Authorization
  ↓
API
  ↓
Provider
  ↓
Transaction
```

Dengan struktur ini, satu platform Crews dapat digunakan oleh:

```text
Flutter App
Backend Service
Internal Service
B2B Partner
Public API
Private API
Provider Integration
Aggregator Integration
```

tanpa mencampurkan credential, user, environment, dan permission masing-masing.

---

# 50. Status Implementasi

| Component                      | Status  |
| ------------------------------ | ------- |
| PostgreSQL                     | ✅       |
| Redis                          | ✅       |
| Health Check                   | ✅       |
| Readiness Check                | ✅       |
| Tenant                         | ✅       |
| Application                    | ✅       |
| Environment                    | ✅       |
| Client                         | ✅       |
| Identity                       | ✅       |
| Human Identity Schema          | ✅       |
| Machine Identity               | ✅       |
| User Account Schema            | ✅       |
| Machine Credential Schema      | ✅       |
| Credential Generation          | ✅       |
| Credential Hashing             | ✅       |
| Client → Identity relation     | ✅       |
| Authentication                 | 🔄 Next |
| Access Token                   | 🔄 Next |
| Authorization                  | ⏳       |
| Role                           | ⏳       |
| Permission                     | ⏳       |
| Scope                          | ⏳       |
| API Registry                   | ⏳       |
| IP Whitelist                   | ⏳       |
| Rate Limit                     | ⏳       |
| Provider Registry              | ⏳       |
| Provider Credential Encryption | ⏳       |
| Transaction                    | ⏳       |
| Idempotency                    | ⏳       |
| Sync / Async                   | ⏳       |
| Queue Worker                   | ⏳       |
| Webhook                        | ⏳       |
| Reconciliation                 | ⏳       |
| Audit Log                      | ⏳       |
| Observability                  | ⏳       |

---

# 51. Terminology

| Istilah             | Arti                                         |
| ------------------- | -------------------------------------------- |
| Tenant              | Pemilik/boundary utama resource              |
| Application         | Sistem/aplikasi yang menggunakan platform    |
| Environment         | Development, staging, atau production        |
| Client              | Consumer teknis dalam environment            |
| Identity            | Aktor yang melakukan authentication          |
| Human Identity      | Identity untuk manusia                       |
| Machine Identity    | Identity untuk service/system                |
| User Account        | Akun manusia                                 |
| Machine Credential  | Credential service-to-service                |
| Client ID           | Identifier credential                        |
| Client Secret       | Secret credential                            |
| API Registry        | Daftar API yang tersedia                     |
| Provider            | Sistem eksternal yang dipanggil Crews        |
| Partner             | Consumer eksternal/B2B                       |
| Scope               | Hak akses terhadap resource/API              |
| Permission          | Hak melakukan action tertentu                |
| Idempotency         | Mekanisme mencegah transaksi ganda           |
| Reconciliation      | Proses mencocokkan transaksi dengan provider |
| Audit Log           | Catatan aktivitas/security                   |
| Raw Response        | Response asli dari provider                  |
| Normalized Response | Response yang sudah distandarkan Crews       |

---

**Dokumen ini harus dianggap sebagai dokumentasi foundation.** Setiap perubahan besar pada model authentication, authorization, credential, API registry, atau transaction harus memperbarui dokumen arsitektur ini.
