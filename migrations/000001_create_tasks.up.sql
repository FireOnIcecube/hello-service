CREATE TABLE tasks (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title text NOT NULL,
    completed boolean NOT NULL DEFAULT false
);