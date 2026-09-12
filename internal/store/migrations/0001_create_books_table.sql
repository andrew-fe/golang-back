CREATE TABLE IF NOT EXISTS books (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    author      TEXT NOT NULL
);

-- ALTER ... IF NOT EXISTS hace esta migración segura también sobre una
-- tabla "books" preexistente creada por una versión anterior del proyecto
-- (que solo tenía id, title y author).
ALTER TABLE books ADD COLUMN IF NOT EXISTS isbn TEXT NOT NULL DEFAULT '';
ALTER TABLE books ADD COLUMN IF NOT EXISTS year INTEGER NOT NULL DEFAULT 0;
ALTER TABLE books ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE books ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_books_title_author ON books (title, author);
