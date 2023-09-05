-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes MODIFY prep_time TINYTEXT;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE recipes MODIFY cook_time TINYTEXT;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE recipes MODIFY total_time TINYTEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes MODIFY prep_time TIME;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE recipes MODIFY cook_time TIME;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE recipes MODIFY total_time TIME;
-- +goose StatementEnd
