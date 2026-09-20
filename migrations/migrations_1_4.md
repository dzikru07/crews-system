┌─────────────────────┐
│       tenants       │
├─────────────────────┤
│ id PK               │
│ name                │
│ slug                │
│ type                │
│ status              │
│ created_at          │
│ updated_at          │
└──────────┬──────────┘
           │
           │ 1:N
           ▼
┌─────────────────────┐
│    applications     │
├─────────────────────┤
│ id PK               │
│ tenant_id FK        │
│ name                │
│ slug                │
│ description         │
│ status              │
│ created_at          │
│ updated_at          │
└──────────┬──────────┘
           │
           │ 1:N
           ▼
┌─────────────────────┐
│    environments     │
├─────────────────────┤
│ id PK               │
│ application_id FK   │
│ name                │
│ type                │
│ status              │
│ created_at          │
│ updated_at          │
└──────────┬──────────┘
           │
           │ 1:N
           ▼
┌─────────────────────┐
│       clients       │
├─────────────────────┤
│ id PK               │
│ environment_id FK   │
│ name                │
│ slug                │
│ type                │
│ status              │
│ created_at          │
│ updated_at          │
└─────────────────────┘