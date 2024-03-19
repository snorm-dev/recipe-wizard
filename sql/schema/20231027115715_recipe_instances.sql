-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipe_instances (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	grocery_list_id INTEGER NOT NULL,
	recipe_id INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipe_instances;
-- +goose StatementEnd
