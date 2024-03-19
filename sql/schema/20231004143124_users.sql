-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	username VARCHAR(32) NOT NULL UNIQUE,
	hashed_password TEXT NOT NULL,
	first_name TEXT,
	last_name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
