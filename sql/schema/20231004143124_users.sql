-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	username VARCHAR(32) NOT NULL,
	UNIQUE (username),
	hashed_password CHAR(60) NOT NULL,
	first_name TEXT,
	last_name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +
