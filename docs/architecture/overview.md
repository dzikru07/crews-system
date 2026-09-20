1. Apa itu Tenant di Crews?

Dalam sistem kita, Tenant adalah boundary kepemilikan data dan akses.

Contoh:

CREWS
│
├── Tenant: Internal Crews
│
├── Tenant: Loket Digital
│
├── Tenant: Atensi
│
└── Tenant: Partner ABC

Tetapi ada satu hal penting:

Tidak semua benda harus menjadi Tenant.

Misalnya provider seperti:

Citilink
Duitku
Provider PPOB
Provider Flight

tidak otomatis menjadi tenant.

Provider adalah resource yang dikelola Crews.

2. Kita bedakan Owner dan Tenant

Saya ingin architecture kita menggunakan dua konsep:

Tenant
Owner

Misalnya:

Tenant
└── Loket Digital
      │
      ├── Application
      ├── Client
      ├── Credential
      └── API Access

Sedangkan:

Crews Platform
└── Provider
      ├── Provider A
      ├── Provider B
      └── Provider C

Provider dimiliki/dikelola platform, bukan oleh Loket.

3. Tenant Types

Kita bisa memiliki:

tenant_type

dengan kemungkinan:

platform
organization
partner

Contoh:

Tenant	Type
Crews	platform
Loket Digital	organization
Atensi	organization
Partner ABC	partner

Namun jangan terlalu mengikat implementasi database ke enum yang sulit diubah.

Untuk awal, kita bisa menggunakan:

type

dengan validation di application layer.

4. Struktur Tenant

Secara konseptual:

Tenant
│
├── Users
│
├── Applications
│
├── Clients
│
├── Credentials
│
├── API Access
│
├── Transactions
│
├── Webhooks
│
└── Reconciliation

Tetapi tidak semua tabel harus langsung memiliki tenant_id.

Misalnya:

providers
api_definitions
permissions

adalah resource platform/global.

5. Application

Sekarang kita masuk ke Application.

Contoh:

Tenant: Loket Digital
│
└── Application: Loket Digital

Kenapa masih perlu Application kalau sudah ada Tenant?

Karena satu organisasi bisa memiliki beberapa aplikasi.

Contoh:

PT Abyakta Teknologi Indonesia
│
├── Loket Digital
├── Atensi
├── Ayo Gas
└── Internal Dashboard

Jadi:

Tenant
   │
   ├── Application A
   ├── Application B
   └── Application C
6. Environment

Application kemudian mempunyai environment.

Application
│
├── Development
├── Staging
└── Production

Saya menyarankan environment menjadi resource eksplisit.

Jadi bukan hanya string:

environment = "production"

di setiap tabel.

Tetapi:

applications
      │
      ▼
environments

Contoh:

Loket Digital
│
├── development
├── staging
└── production
7. Kenapa Environment perlu tabel?

Karena nantinya environment dapat mempunyai configuration.

Misalnya:

production
├── API base URL
├── provider configuration
├── rate limit
├── webhook configuration
└── security policy

Sedangkan:

development
├── sandbox provider
└── relaxed rate limit

Namun secret tetap disimpan di credential system, bukan configuration biasa.

8. Client

Sekarang:

Tenant
  ↓
Application
  ↓
Environment
  ↓
Client

Contoh:

Loket Digital
│
├── Development
│   └── Mobile Backend
│
├── Staging
│   └── Mobile Backend
│
└── Production
    ├── Mobile Backend
    └── Reconciliation Worker
9. Kenapa Client bukan Application?

Karena satu application dapat mempunyai banyak caller.

Misalnya:

Loket Digital
│
└── Production
    │
    ├── Backend API
    ├── Payment Worker
    ├── Reconciliation Worker
    └── Admin Service

Masing-masing dapat diberikan permission berbeda.

10. Credential

Client kemudian memiliki credential:

Client
  │
  ├── Credential A
  ├── Credential B
  └── Credential C

Misalnya:

loket-production-backend
│
└── API Key

atau:

partner-production
│
└── HMAC Credential
11. ERD awal

Sekarang bentuk utamanya:

┌──────────────┐
│   tenants    │
│              │
│ id           │
│ name         │
│ type         │
│ status       │
└──────┬───────┘
       │
       │ 1:N
       ▼
┌──────────────────┐
│   applications   │
│                  │
│ id               │
│ tenant_id        │
│ name             │
│ slug             │
│ status            │
└────────┬─────────┘
         │
         │ 1:N
         ▼
┌──────────────────┐
│  environments    │
│                  │
│ id               │
│ application_id   │
│ name             │
│ type             │
│ status            │
└────────┬─────────┘
         │
         │ 1:N
         ▼
┌──────────────────┐
│     clients      │
│                  │
│ id               │
│ environment_id   │
│ name             │
│ type             │
│ status            │
└────────┬─────────┘
         │
         │ 1:N
         ▼
┌──────────────────┐
│   credentials    │
│                  │
│ id               │
│ client_id        │
│ type             │
│ secret_ciphertext│
│ status            │
└──────────────────┘

Ini menjadi salah satu backbone utama Crews.

12. Tambahkan Partner

Sekarang partner.

Partner sebenarnya adalah business entity, bukan sekadar client.

Contoh:

Partner ABC
│
├── Users
├── Applications
└── Clients

Tetapi untuk awal kita tidak perlu memaksa partner menjadi Application.

Kita bisa modelkan:

partners

sebagai resource tersendiri.

13. Partner ERD
┌──────────────┐
│   tenants    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   partners   │
│              │
│ id           │
│ tenant_id    │
│ name         │
│ code         │
│ status       │
└──────┬───────┘
       │
       ├──────────────┐
       ▼              ▼
┌─────────────┐  ┌──────────────┐
│ partner_    │  │ partner_     │
│ users       │  │ clients      │
└─────────────┘  └──────┬───────┘
                        │
                        ▼
                   credentials

Tetapi ada pertanyaan penting:

Apakah Partner ABC menjadi Tenant?

Untuk architecture kita, saya sarankan:

Partner tetap menjadi entity bisnis yang terpisah dari Tenant.

Alasannya:

Crews bisa memiliki partner:

Crews
│
├── Partner ABC
├── Partner XYZ
└── Partner DEF

tetapi partner tidak berarti mereka memiliki seluruh isolation boundary seperti organisasi internal.

14. Maka struktur finalnya
                         CREWS PLATFORM
                               │
              ┌────────────────┴────────────────┐
              │                                 │
              ▼                                 ▼
           TENANTS                           PARTNERS
              │                                 │
       ┌──────┼──────┐                    ┌─────┴─────┐
       │      │      │                    │           │
       ▼      ▼      ▼                    ▼           ▼
     Apps   Users   Config             Users       Clients
       │
       ▼
 Environments
       │
       ▼
    Clients
       │
       ▼
  Credentials
       │
       ▼
 Authorization
15. Global vs Tenant Resource

Ini sangat penting untuk database nanti.

Global

Tidak memiliki tenant_id.

Contoh:

permissions
api_definitions
providers
provider_types
Tenant-owned

Memiliki tenant_id.

Contoh:

applications
Application-owned

Memiliki:

application_id

Contoh:

environments
Environment-owned

Memiliki:

environment_id

Contoh:

clients
Client-owned

Memiliki:

client_id

Contoh:

credentials
16. Jangan melakukan ini

Jangan membuat semua tabel:

tenant_id

secara otomatis.

Misalnya:

permissions

tidak perlu:

tenant_id

karena permission merupakan katalog global.

Begitu juga:

api_definitions

bisa global.

17. API Registry

Nanti akan ada:

api_registry

misalnya:

Flight API
Payment API
PPOB API
Booking API
Deposit API

Kemudian endpoint:

POST /v1/booking
GET  /v1/booking/:id
POST /v1/payment

Struktur:

API
│
└── Versions
      │
      └── Endpoints

Misalnya:

Flight API
│
└── v1
    ├── POST /search
    ├── POST /booking
    └── GET /booking/:id

Access diberikan kepada client.

Client
   │
   ▼
API Access
   │
   ▼
API / Endpoint
18. Security Boundary

Sekarang kita bisa mulai menentukan security boundary:

                    REQUEST
                       │
                       ▼
                Request ID
                       │
                       ▼
                TLS / HTTPS
                       │
                       ▼
                 IP Policy
                       │
                       ▼
               Authentication
                       │
                       ▼
                  Principal
                       │
                       ▼
                    Tenant
                       │
                       ▼
                  Client
                       │
                       ▼
                Authorization
                       │
                       ▼
                  API Access
                       │
                       ▼
                 Rate Limit
                       │
                       ▼
                   Handler

Ini akan menjadi dasar middleware kita nanti.

19. Contoh nyata: Loket Digital

Sekarang kita tes desain ini.

Tenant
└── PT Abyakta / Loket Digital

Application:

Application
└── Loket Digital

Environment:

Loket Digital
├── development
├── staging
└── production

Client:

production
├── backend
├── payment-worker
└── reconciliation-worker

Credential:

backend
└── production-api-key

payment-worker
└── production-worker-key

Mobile users:

Loket Digital DB
├── User A
├── User B
└── User C

User A/B/C tidak masuk ke database Crews.

Ini sesuai requirement yang Anda jelaskan sebelumnya.

20. Contoh B2B

Misalnya Partner ABC ingin memakai API Crews.

Partner ABC
│
├── User
│   └── admin@abc
│
└── Client
    └── abc-production
          │
          ▼
       Credential
          │
          ▼
      API Access
          │
          ▼
      Booking API

Security:

ABC Request
   │
   ├── HTTPS
   ├── IP whitelist
   ├── Authentication
   ├── Client validation
   ├── API permission
   ├── Rate limit
   └── Request validation
21. Contoh Crews memanggil Provider

Ini sisi lainnya.

Tenant/Application
        │
        ▼
Crews Transaction
        │
        ▼
Provider Interface
        │
        ▼
Provider
        │
        ▼
Provider Credential

Credential provider tidak pernah diberikan ke partner.

Contoh:

Partner ABC
    ↓
Crews
    ↓
Provider A

Partner hanya mengetahui:

Crews API

bukan:

Provider A username
Provider A password
Provider A API key
22. ERD backbone final sementara

Ini yang saya ingin kita jadikan acuan:

                         ┌──────────────┐
                         │    TENANT    │
                         └───────┬──────┘
                                 │
                                 │
                         ┌───────▼────────┐
                         │  APPLICATION   │
                         └───────┬────────┘
                                 │
                         ┌───────▼────────┐
                         │  ENVIRONMENT   │
                         └───────┬────────┘
                                 │
                         ┌───────▼────────┐
                         │     CLIENT     │
                         └───────┬────────┘
                                 │
                         ┌───────▼────────┐
                         │   CREDENTIAL   │
                         └───────┬────────┘
                                 │
                                 ▼
                          AUTHORIZATION
                                 │
                                 ▼
                            API ACCESS
                                 │
                                 ▼
                              API
                                 │
                                 ▼
                            ENDPOINT


       ┌──────────────┐
       │   PARTNER    │
       └──────┬───────┘
              │
        ┌─────┴──────┐
        ▼            ▼
    PARTNER       PARTNER
     USERS        CLIENTS
                     │
                     ▼
                CREDENTIALS


       ┌──────────────┐
       │   PROVIDER   │
       └──────┬───────┘
              │
              ▼
       PROVIDER CREDENTIAL
23. Keputusan architecture yang kita kunci

Untuk sementara kita tetapkan:

Konsep	Keputusan
Crews	API Platform
Tenant	Boundary kepemilikan/isolation
Application	Aplikasi yang menggunakan Crews
Environment	Dev/Staging/Production
Client	Machine caller
Credential	Authentication material
Platform User	Manusia pengelola Crews
Mobile User	Milik aplikasi, bukan Crews
Partner	Business entity
Partner User	User dashboard partner
Partner Client	Machine caller partner
Provider	External API/provider
Provider Credential	Credential untuk provider
API	API capability
Endpoint	Operasi API
Permission	Hak akses granular
API Access	Hubungan client dengan API