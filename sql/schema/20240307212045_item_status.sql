-- +goose Up
-- +goose StatementBegin
ALTER TABLE items 
	ADD COLUMN is_complete BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items DROP is_complete;
-- +goose StatementEnd
