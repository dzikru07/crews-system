DROP INDEX IF EXISTS idx_clients_identity_id;

ALTER TABLE clients
DROP CONSTRAINT IF EXISTS clients_identity_fk;

ALTER TABLE clients
DROP COLUMN IF EXISTS identity_id;