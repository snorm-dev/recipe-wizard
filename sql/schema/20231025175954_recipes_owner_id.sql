-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes ADD COLUMN owner_id INTEGER NOT NULL DEFAULT -1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes DROP COLUMN owner_id;
-- +goose StatementEnd
