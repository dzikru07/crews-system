CREATE TABLE api_access_grants (
    id UUID PRIMARY KEY,

    client_id UUID NOT NULL,
    endpoint_id UUID NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    scopes JSONB NOT NULL DEFAULT '[]'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT api_access_grants_client_fk
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE CASCADE,

    CONSTRAINT api_access_grants_endpoint_fk
        FOREIGN KEY (endpoint_id)
        REFERENCES api_endpoints(id)
        ON DELETE CASCADE,

    CONSTRAINT api_access_grants_unique
        UNIQUE (client_id, endpoint_id),

    CONSTRAINT api_access_grants_status_check
        CHECK (
            status IN (
                'active',
                'disabled'
            )
        )
);

CREATE INDEX idx_api_access_grants_client_id
    ON api_access_grants(client_id);

CREATE INDEX idx_api_access_grants_endpoint_id
    ON api_access_grants(endpoint_id);

CREATE INDEX idx_api_access_grants_status
    ON api_access_grants(status);