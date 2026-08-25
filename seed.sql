CREATE TABLE IF NOT EXISTS tasks (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title text NOT NULL,
    completed boolean NOT NULL DEFAULT false
);

INSERT INTO tasks (title,completed)
VALUES
    ('Learn PostgreSQL', false),
    ('Learn Docker', false),
    ('Build API', true);