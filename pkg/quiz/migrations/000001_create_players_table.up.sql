-- DROP TABLE IF EXISTS schema_migrations;

CREATE TABLE IF NOT EXISTS players (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    joined timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    last_update timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    score integer NOT NULL DEFAULT 0
);


CREATE TABLE IF NOT EXISTS quizes (
    id bigserial PRIMARY KEY,
	category text DEFAULT 'general',
	reward integer DEFAULT 0,
	questions text[] DEFAULT '{}',
    answers text[] DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS games (
    id bigserial PRIMARY KEY,
    finished timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    player bigserial,
    quiz bigserial,
    FOREIGN KEY (player) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (quiz) REFERENCES quizes(id) ON DELETE CASCADE
);
