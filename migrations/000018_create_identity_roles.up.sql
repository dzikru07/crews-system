CREATE TABLE identity_roles (
    identity_id UUID NOT NULL,

    role_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (
        identity_id,
        role_id
    ),

    CONSTRAINT identity_roles_identity_fk
        FOREIGN KEY (identity_id)
        REFERENCES identities(id)
        ON DELETE CASCADE,

    CONSTRAINT identity_roles_role_fk
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_identity_roles_role_id
    ON identity_roles(role_id);