CREATE TABLE applications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,

    name VARCHAR(150) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT applications_tenant_fk
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE RESTRICT,

    CONSTRAINT applications_tenant_slug_unique
        UNIQUE (tenant_id, slug)
);

CREATE INDEX idx_applications_tenant_id
    ON applications(tenant_id);

CREATE INDEX idx_applications_status
    ON applications(status);