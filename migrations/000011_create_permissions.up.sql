CREATE TABLE permissions (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    name VARCHAR(150) NOT NULL,

    description TEXT,

    resource VARCHAR(100) NOT NULL,

    action VARCHAR(100) NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT permissions_status_check
        CHECK (status IN ('active', 'disabled')),

    CONSTRAINT permissions_unique_name
        UNIQUE (tenant_id, name)
);

CREATE INDEX idx_permissions_tenant_id
    ON permissions(tenant_id);

CREATE INDEX idx_permissions_resource_action
    ON permissions(tenant_id, resource, action);