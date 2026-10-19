CREATE TABLE tenants (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    name_key TEXT NOT NULL UNIQUE,
    slug     TEXT NOT NULL UNIQUE,
    status     TEXT NOT NULL DEFAULT 'active'
               CHECK (status IN ('active', 'suspended')),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT tenants_name_key_unique UNIQUE (name_key),
    CONSTRAINT tenants_slug_unique UNIQUE (slug)
);
