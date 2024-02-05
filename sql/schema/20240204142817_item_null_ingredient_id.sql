-- +goose Up
-- +goose StatementBegin
ALTER TABLE items MODIFY COLUMN ingredient_id BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items MODIFY COLUMN ingredient_id BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd
