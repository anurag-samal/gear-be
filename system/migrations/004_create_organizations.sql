CREATE TYPE business_type AS ENUM (
    'B2B',
    'B2C',
    'BOTH'
);

CREATE TABLE organizations (
    id UUID PRIMARY KEY,

    name TEXT NOT NULL,

    business_type business_type NOT NULL,

    industry TEXT,

    country TEXT,

    owner_user_id UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);