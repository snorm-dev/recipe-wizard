-- +goose Up
-- +goose StatementBegin
CREATE TABLE ingredients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    recipe_id BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ingredients;
-- +goose StatementEnd
