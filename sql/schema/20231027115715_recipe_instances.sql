-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipe_instances (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	grocery_list_id BIGINT NOT NULL,
	recipe_id BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipe_instances;
-- +goose StatementEnd
