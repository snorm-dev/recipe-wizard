-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipe_instances
RENAME TO meals;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meals
RENAME TO recipe_instances;
-- +goose StatementEnd
