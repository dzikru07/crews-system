CREATE TABLE clients (
    id UUID PRIMARY KEY,

    environment_id UUID NOT NULL,

    name VARCHAR(150) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT clients_environment_fk
        FOREIGN KEY (environment_id)
        REFERENCES environments(id)
        ON DELETE RESTRICT,

    CONSTRAINT clients_environment_slug_unique
        UNIQUE (environment_id, slug)
);

CREATE INDEX idx_clients_environment_id
    ON clients(environment_id);

CREATE INDEX idx_clients_status
    ON clients(status);