BEGIN;

-- users: платформа-агностичные пользователи
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- user_identities: соответствие транспорта и внешнего ID -> внутренний user_id
CREATE TABLE user_identities (
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transport   TEXT        NOT NULL,
    external_id TEXT        NOT NULL,
    PRIMARY KEY (transport, external_id)
);
CREATE INDEX user_identities_user_id_idx ON user_identities(user_id);

-- incomes: доходы в копейках, дата учёта как DATE (UTC полночь), soft-delete через voided_at
CREATE TABLE incomes (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    at          DATE        NOT NULL,                         -- дата учёта (UTC день)
    amount      BIGINT      NOT NULL CHECK (amount > 0),     -- деньги в копейках
    note        TEXT,                                         -- опциональный комментарий
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),           -- когда запись создана
    voided_at   TIMESTAMPTZ                                   -- soft-delete (NULL = активна)
);

-- Индекс для быстрых агрегатов по "живым" записям внутри периода
CREATE INDEX incomes_user_at_active_idx
  ON incomes (user_id, at)
  WHERE voided_at IS NULL;

COMMIT;
