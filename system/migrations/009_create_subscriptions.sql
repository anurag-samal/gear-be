CREATE TYPE subscription_plan AS ENUM (
    'TRIAL',
    'BASE',
    'PRO'
);

CREATE TYPE subscription_status AS ENUM (
    'ACTIVE',
    'EXPIRED',
    'CANCELLED'
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,

    organization_id UUID NOT NULL UNIQUE,

    plan subscription_plan NOT NULL,

    status subscription_status NOT NULL,

    trial_ends_at TIMESTAMPTZ,

    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE
);