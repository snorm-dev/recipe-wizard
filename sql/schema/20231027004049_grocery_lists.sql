-- +goose Up
-- +goose StatementBegin
CREATE TABLE grocery_lists (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	owner_id INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grocery_lists;
-- +goose StatementEnd
