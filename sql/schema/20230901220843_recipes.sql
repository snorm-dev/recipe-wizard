-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipes (
	id VARCHAR(36) PRIMARY KEY,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipes;
-- +goose StatementEnd
