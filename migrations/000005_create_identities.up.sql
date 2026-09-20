CREATE TABLE identities (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL,

    type VARCHAR(30) NOT NULL,

    name VARCHAR(150) NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    metadata JSONB,

    last_authenticated_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT identities_tenant_fk
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE RESTRICT,

    CONSTRAINT identities_type_check
        CHECK (
            type IN (
                'human',
                'machine'
            )
        ),

    CONSTRAINT identities_status_check
        CHECK (
            status IN (
                'active',
                'disabled',
                'suspended'
            )
        )
);

CREATE INDEX idx_identities_tenant_id
    ON identities(tenant_id);

CREATE INDEX idx_identities_type
    ON identities(type);

CREATE INDEX idx_identities_status
    ON identities(status);