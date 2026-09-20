CREATE TABLE user_accounts (
    id UUID PRIMARY KEY,

    identity_id UUID NOT NULL UNIQUE,

    email VARCHAR(255) NOT NULL,

    username VARCHAR(100),

    password_hash TEXT,

    email_verified_at TIMESTAMPTZ,

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT user_accounts_identity_fk
        FOREIGN KEY (identity_id)
        REFERENCES identities(id)
        ON DELETE CASCADE,

    CONSTRAINT user_accounts_email_unique
        UNIQUE (email)
);

CREATE UNIQUE INDEX user_accounts_username_unique
    ON user_accounts(username)
    WHERE username IS NOT NULL;

CREATE INDEX idx_user_accounts_identity_id
    ON user_accounts(identity_id);

CREATE INDEX idx_user_accounts_email
    ON user_accounts(email);