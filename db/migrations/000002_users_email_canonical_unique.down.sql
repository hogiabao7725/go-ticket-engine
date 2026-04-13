DROP INDEX IF EXISTS idx_users_email_canonical_unique;

ALTER TABLE users
ADD CONSTRAINT users_email_key UNIQUE (email);
