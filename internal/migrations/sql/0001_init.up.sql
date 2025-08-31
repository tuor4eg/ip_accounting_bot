BEGIN;

-- === users ===
CREATE TABLE users (
    id                      BIGSERIAL PRIMARY KEY,
    tax_scheme              TEXT NOT NULL DEFAULT 'usn_6' CHECK (tax_scheme IN ('usn_6','usn_dr')),
    -- PII consent (opt-in): сохраняем факт согласия и версию текста согласия
    pii_consent_version     TEXT,
    pii_consent_at          TIMESTAMPTZ,
    pii_consent_revoked_at  TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- === user identities (transport mapping) ===
CREATE TABLE user_identities (
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transport   TEXT        NOT NULL,
    external_id TEXT        NOT NULL,
    PRIMARY KEY (transport, external_id)
);
CREATE INDEX user_identities_user_id_idx ON user_identities(user_id);

-- === counterparties (per-user; archived, not deleted) ===
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

-- === categories (per-user; scope: income|expense|both) ===
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

-- === incomes (gross receipts) ===
CREATE TABLE incomes (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    at              DATE        NOT NULL,                         -- UTC day
    amount          BIGINT      NOT NULL CHECK (amount > 0),      -- kopecks
    note            TEXT,
    counterparty_id BIGINT      NULL REFERENCES counterparties(id) ON DELETE SET NULL,
    category_id     BIGINT      NULL REFERENCES categories(id)     ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    voided_at       TIMESTAMPTZ                                   -- soft-delete
);
-- быстрые агрегаты по активным строкам
CREATE INDEX incomes_user_at_active_idx
  ON incomes (user_id, at)
  WHERE voided_at IS NULL;

-- === payments (contrib/advance) ===
CREATE TABLE payments (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    at          DATE        NOT NULL,                         -- UTC day
    type        TEXT        NOT NULL CHECK (type IN ('contrib','advance')),
    amount      BIGINT      NOT NULL CHECK (amount > 0),      -- kopecks
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

COMMIT;
