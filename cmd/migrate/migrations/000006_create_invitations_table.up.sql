CREATE TABLE invitations (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    tenant_id  INTEGER NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    email      TEXT COLLATE NOCASE,
    role       TEXT NOT NULL DEFAULT 'member'
               CHECK (role = 'member'),
    expires_at TEXT NOT NULL,
    used_at    TEXT,
    revoked_at TEXT,
    created_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT invitations_creator_same_tenant_fk
        FOREIGN KEY (tenant_id, created_by)
        REFERENCES users (tenant_id, id) ON DELETE CASCADE
);

CREATE INDEX invitations_tenant_created_idx
    ON invitations (tenant_id, created_at DESC);
CREATE INDEX invitations_tenant_email_idx
    ON invitations (tenant_id, email);
