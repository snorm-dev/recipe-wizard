-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipes (
	id VARCHAR(36) PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	description TEXT,
	url VARCHAR(512) CHARACTER SET 'ascii' COLLATE 'ascii_general_ci',
	prep_time TIME,
	cook_time TIME,
	cool_time TIME
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipes;
-- +goose StatementEnd
