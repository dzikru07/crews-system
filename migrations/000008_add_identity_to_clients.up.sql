ALTER TABLE clients
ADD COLUMN identity_id UUID;

ALTER TABLE clients
ADD CONSTRAINT clients_identity_fk
FOREIGN KEY (identity_id)
REFERENCES identities(id)
ON DELETE RESTRICT;

CREATE INDEX idx_clients_identity_id
ON clients(identity_id);