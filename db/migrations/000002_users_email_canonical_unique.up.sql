-- Replace case-sensitive email unique constraint with canonical unique index.
-- Canonical form: lower(trim(email)).
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_canonical_unique
ON users ((lower(btrim(email))));
