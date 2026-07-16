CREATE TABLE auth_accounts (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL,

    provider TEXT NOT NULL,

    provider_user_id TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(provider, provider_user_id),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);