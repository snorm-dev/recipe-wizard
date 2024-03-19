-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	description TEXT,
	url VARCHAR(512),
	prep_time TEXT,
	cook_time TEXT,
	total_time TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipes;
-- +goose StatementEnd
