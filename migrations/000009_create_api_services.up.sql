CREATE TABLE api_services (
    id UUID PRIMARY KEY,

    application_id UUID NOT NULL,

    name VARCHAR(150) NOT NULL,

    slug VARCHAR(100) NOT NULL,

    description TEXT,

    base_path VARCHAR(255),

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT api_services_application_fk
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE RESTRICT,

    CONSTRAINT api_services_application_slug_unique
        UNIQUE (application_id, slug),

    CONSTRAINT api_services_status_check
        CHECK (
            status IN (
                'active',
                'disabled'
            )
        )
);

CREATE INDEX idx_api_services_application_id
    ON api_services(application_id);

CREATE INDEX idx_api_services_status
    ON api_services(status);