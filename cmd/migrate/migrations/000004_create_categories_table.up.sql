CREATE TABLE IF NOT EXISTS categories (
id TEXT PRIMARY KEY,
name TEXT NOT NULL COLLATE NOCASE UNIQUE, -- schützt vor Duplikaten mit unterschiedlicher Groß-/Kleinschreibung
is_default INTEGER NOT NULL DEFAULT 0
CHECK (is_default IN (0, 1))
);

INSERT INTO categories (id, name, is_default) VALUES
('1', 'Meeting', 1),
('2', 'Workshop', 1),
('3', 'Team-Event', 1),
('4', 'Planning', 1),
('5', 'UI&UX', 1),
('6', 'Konferenz', 1),
('7', 'Weiterbildung', 1),
('8', 'Product Launch', 1);
