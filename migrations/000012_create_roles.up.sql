CREATE TABLE roles (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    name VARCHAR(150) NOT NULL,

    description TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT roles_status_check
        CHECK (status IN ('active', 'disabled')),

    CONSTRAINT roles_unique_name
        UNIQUE (tenant_id, name)
);

CREATE INDEX idx_roles_tenant_id
    ON roles(tenant_id);