-- +goose Up
-- +goose StatementBegin
ALTER TABLE items
RENAME COLUMN recipe_instance_id
TO meal_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items
RENAME COLUMN meal_id
TO recipe_instance_id;
-- +goose StatementEnd
