CREATE TABLE api_endpoints (
    id UUID PRIMARY KEY,

    service_id UUID NOT NULL,

    name VARCHAR(150) NOT NULL,

    method VARCHAR(10) NOT NULL,

    path VARCHAR(255) NOT NULL,

    description TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    requires_auth BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT api_endpoints_service_fk
        FOREIGN KEY (service_id)
        REFERENCES api_services(id)
        ON DELETE RESTRICT,

    CONSTRAINT api_endpoints_method_check
        CHECK (
            method IN (
                'GET',
                'POST',
                'PUT',
                'PATCH',
                'DELETE'
            )
        ),

    CONSTRAINT api_endpoints_status_check
        CHECK (
            status IN (
                'active',
                'disabled'
            )
        ),

    CONSTRAINT api_endpoints_service_method_path_unique
        UNIQUE (
            service_id,
            method,
            path
        )
);

CREATE INDEX idx_api_endpoints_service_id
    ON api_endpoints(service_id);

CREATE INDEX idx_api_endpoints_status
    ON api_endpoints(status);

CREATE INDEX idx_api_endpoints_method
    ON api_endpoints(method);