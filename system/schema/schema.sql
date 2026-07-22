-- ============================================================
-- Gear Database Schema
-- Auto-generated from system/migrations/ — do not edit directly
-- Generated: 2026-07-21T19:06:34Z
-- ============================================================

-- ============================================================
-- Create users
-- ============================================================

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

-- ============================================================
-- Create user credentials
-- ============================================================

CREATE TABLE user_credentials (
    user_id UUID PRIMARY KEY,

    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- ============================================================
-- Create auth accounts
-- ============================================================

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

-- ============================================================
-- Create organizations
-- ============================================================

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

-- ============================================================
-- Create teams
-- ============================================================

CREATE TABLE teams (
    id UUID PRIMARY KEY,

    organization_id UUID NOT NULL,

    name TEXT NOT NULL,

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE
);

-- ============================================================
-- Create team members
-- ============================================================

CREATE TYPE team_role AS ENUM (
    'OWNER',
    'ADMIN',
    'MEMBER',
    'VIEWER'
);

CREATE TABLE team_members (
    team_id UUID NOT NULL,

    user_id UUID NOT NULL,

    role team_role NOT NULL,

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY(team_id, user_id),

    FOREIGN KEY (team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- ============================================================
-- Create refresh tokens
-- ============================================================

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL,

    token_hash BYTEA NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,

    revoked BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- ============================================================
-- Create invitations
-- ============================================================

CREATE TYPE invitation_status AS ENUM (
    'PENDING',
    'ACCEPTED',
    'EXPIRED',
    'REVOKED'
);

CREATE TABLE invitations (
    id UUID PRIMARY KEY,

    organization_id UUID NOT NULL,

    team_id UUID NOT NULL,

    email TEXT NOT NULL,

    invited_by UUID NOT NULL,

    status invitation_status NOT NULL DEFAULT 'PENDING',

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,

    FOREIGN KEY (team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE,

    FOREIGN KEY (invited_by)
        REFERENCES users(id)
);

-- ============================================================
-- Create subscriptions
-- ============================================================

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

