CREATE TABLE endpoint_permissions (
    endpoint_id UUID NOT NULL
        REFERENCES api_endpoints(id)
        ON DELETE CASCADE,

    permission_id UUID NOT NULL
        REFERENCES permissions(id)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (endpoint_id, permission_id)
);

CREATE INDEX idx_endpoint_permissions_permission_id
    ON endpoint_permissions(permission_id);