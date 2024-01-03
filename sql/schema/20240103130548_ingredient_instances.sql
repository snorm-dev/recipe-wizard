-- +goose Up
-- +goose StatementBegin
CREATE TABLE ingredient_instances (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	grocery_list_id BIGINT NOT NULL,
	recipe_instance_id BIGINT,
	ingredient_id BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ingredient_instances;
-- +goose StatementEnd
