-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipes (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	description TEXT,
	url VARCHAR(512) CHARACTER SET 'ascii' COLLATE 'ascii_general_ci',
	prep_time TINYTEXT,
	cook_time TINYTEXT,
	total_time TINYTEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipes;
-- +goose StatementEnd
