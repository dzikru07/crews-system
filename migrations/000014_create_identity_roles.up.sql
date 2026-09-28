CREATE TABLE identity_roles (
    identity_id UUID NOT NULL
        REFERENCES identities(id)
        ON DELETE CASCADE,

    role_id UUID NOT NULL
        REFERENCES roles(id)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (identity_id, role_id)
);

CREATE INDEX idx_identity_roles_role_id
    ON identity_roles(role_id);