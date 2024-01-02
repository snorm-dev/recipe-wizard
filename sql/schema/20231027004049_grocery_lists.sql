-- +goose Up
-- +goose StatementBegin
CREATE TABLE grocery_lists (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	owner_id BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grocery_lists;
-- +goose StatementEnd
