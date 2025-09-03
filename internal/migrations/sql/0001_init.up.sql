-- 0001_init.sql
-- IP Accounting Bot — init schema

BEGIN;

-- ====== SCHEMAS ======
CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS pii;    -- sensitive data (restricted access)

-- ====== users ======
CREATE TABLE users (
    id                      BIGSERIAL PRIMARY KEY,
    tax_scheme              TEXT NOT NULL DEFAULT 'usn_6' CHECK (tax_scheme IN ('usn_6','usn_dr')),
    -- PII consent (opt-in): store version of consent text and timestamps
    pii_consent_version     TEXT,
    pii_consent_at          TIMESTAMPTZ,
    pii_consent_revoked_at  TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ====== user identities (transport mapping, HMAC instead of raw external_id) ======
-- external_hash = HMAC-SHA256(secret, transport || '|' || external_id)
CREATE TABLE user_identities (
    user_id       BIGINT     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transport     TEXT       NOT NULL,                         -- e.g. 'telegram'
    external_hash BYTEA      NOT NULL CHECK (octet_length(external_hash) = 32),
    hmac_kid      SMALLINT   NOT NULL DEFAULT 1,               -- HMAC key version (for rotation)
    PRIMARY KEY (transport, external_hash, hmac_kid)
);
CREATE INDEX user_identities_user_id_idx ON user_identities(user_id);

-- ====== user profile (non-sensitive data) ======
CREATE TABLE user_profile (
    user_id       BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    ip_reg_year   INT     NOT NULL,
    ip_reg_month  SMALLINT NOT NULL CHECK (ip_reg_month BETWEEN 1 AND 12),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ====== counterparties (per-user; archived, not deleted) ======
CREATE TABLE counterparties (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    name_norm   TEXT        GENERATED ALWAYS AS (lower(btrim(name))) STORED,
    type        TEXT, -- optional: client|vendor|other
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    archived_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX counterparties_user_norm_uq
    ON counterparties(user_id, name_norm);

-- ====== categories (per-user; scope: income|expense|both) ======
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    name_norm   TEXT        GENERATED ALWAYS AS (lower(btrim(name))) STORED,
    scope       TEXT        NOT NULL CHECK (scope IN ('income','expense','both')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    archived_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX categories_user_norm_scope_uq
    ON categories(user_id, name_norm, scope);

-- ====== incomes (gross receipts) ======
CREATE TABLE incomes (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    at              DATE        NOT NULL,                         -- stored in UTC day (DATE only)
    amount          BIGINT      NOT NULL CHECK (amount > 0),      -- stored in kopecks
    note            TEXT,
    counterparty_id BIGINT      NULL REFERENCES counterparties(id) ON DELETE SET NULL,
    category_id     BIGINT      NULL REFERENCES categories(id)     ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    voided_at       TIMESTAMPTZ                                   -- soft-delete
);
-- fast aggregates on active rows
CREATE INDEX incomes_user_at_active_idx
  ON incomes (user_id, at)
  WHERE voided_at IS NULL;

-- ====== payments (contributions / advances) ======
CREATE TABLE payments (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    at          DATE        NOT NULL,                         -- stored in UTC day
    type        TEXT        NOT NULL CHECK (type IN ('contrib','advance')),
    amount      BIGINT      NOT NULL CHECK (amount > 0),      -- stored in kopecks
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    voided_at   TIMESTAMPTZ                                   -- soft-delete
);
CREATE INDEX payments_user_at_active_idx
  ON payments (user_id, at)
  WHERE voided_at IS NULL;
CREATE INDEX payments_user_type_at_active_idx
  ON payments (user_id, type, at)
  WHERE voided_at IS NULL;

-- ====== PII (separate schema) ======
-- telegram chat_id is stored encrypted (AES-GCM box of int64 packed in 8 bytes BigEndian)
CREATE TABLE pii.telegram (
    user_id    BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    chat_enc   BYTEA      NOT NULL,
    enc_kid    SMALLINT   NOT NULL DEFAULT 1,     -- encryption key version (for rotation)
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;
