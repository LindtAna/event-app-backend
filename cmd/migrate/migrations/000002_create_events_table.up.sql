CREATE TABLE IF NOT EXISTS events(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT DEFAULT '',
    location TEXT NOT NULL,
    start_date_time TEXT NOT NULL,
    end_date_time TEXT NOT NULL,
    category_id TEXT DEFAULT '',
    url TEXT DEFAULT '',
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);