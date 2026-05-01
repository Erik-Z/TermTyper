PRAGMA foreign_keys = ON;

ALTER TABLE users ADD COLUMN display_name TEXT DEFAULT '' CHECK(length(display_name) <= 20);
