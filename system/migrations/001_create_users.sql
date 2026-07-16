CREATE TABLE users (
    id UUID PRIMARY KEY,

    organization_id UUID,

    email TEXT NOT NULL UNIQUE,

    full_name TEXT NOT NULL,

    avatar_url TEXT,

    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);