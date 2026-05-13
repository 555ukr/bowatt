BEGIN;

CREATE TABLE IF NOT EXIST photo (
	path text NOT NULL,
    tags TEXT[],
	created_at timestamptz NOT NULL,
);

COMMIT;