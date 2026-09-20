CREATE TABLE environments (
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL,

    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT environments_application_fk
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE RESTRICT,

    CONSTRAINT environments_application_type_unique
        UNIQUE (application_id, type)
);

CREATE INDEX idx_environments_application_id
    ON environments(application_id);

CREATE INDEX idx_environments_status
    ON environments(status);