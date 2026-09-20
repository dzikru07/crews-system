CREATE TABLE machine_credentials (
    id UUID PRIMARY KEY,

    identity_id UUID NOT NULL,

    client_id VARCHAR(100) NOT NULL,

    client_secret_hash TEXT NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    expires_at TIMESTAMPTZ,

    last_used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT machine_credentials_identity_fk
        FOREIGN KEY (identity_id)
        REFERENCES identities(id)
        ON DELETE CASCADE,

    CONSTRAINT machine_credentials_client_id_unique
        UNIQUE (client_id),

    CONSTRAINT machine_credentials_status_check
        CHECK (
            status IN (
                'active',
                'revoked',
                'expired'
            )
        )
);

CREATE INDEX idx_machine_credentials_identity_id
    ON machine_credentials(identity_id);

CREATE INDEX idx_machine_credentials_status
    ON machine_credentials(status);