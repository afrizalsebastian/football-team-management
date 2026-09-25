CREATE TABLE IF NOT EXISTS admin_account (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username      VARCHAR(255) NOT NULL,
  password      TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  is_deleted    BOOL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_admin_accoutn_active_username
ON admin_account (username)
WHERE is_deleted = false;