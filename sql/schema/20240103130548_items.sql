-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	grocery_list_id INTEGER NOT NULL,
	recipe_instance_id INTEGER,
	ingredient_id INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd
