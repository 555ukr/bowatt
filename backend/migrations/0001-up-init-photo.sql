BEGIN;

CREATE TABLE IF NOT EXISTS photo (
	path text NOT NULL,
	tags TEXT[],
	created_at timestamptz NOT NULL
);

COMMIT;
