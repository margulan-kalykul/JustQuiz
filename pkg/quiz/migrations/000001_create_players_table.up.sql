DROP TABLE IF EXISTS schema_migrations;

CREATE TABLE IF NOT EXISTS players (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    joined timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    last_update timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    score integer NOT NULL DEFAULT 0
);

-- CREATE TABLE IF NOT EXISTS quizes (
--     id bigserial PRIMARY KEY,
-- 	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
-- 	updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
-- 	category  text,
-- 	dificulty integer DEFAULT 0,
-- 	q1 text NOT NULL,
-- 	a1 text NOT NULL,
-- 	q2 text NOT NULL,
-- 	a2 text NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS games (
--     id bigserial PRIMARY KEY,
--     time timestamp(0) with time zone NOT NULL DEFAULT NOW(),
--     place text,
--     FOREIGN KEY (quiz) REFERENCES quizes(id),
--     FOREIGN KEY (player) REFERENCES players(id)
-- );
