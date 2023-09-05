-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipe_ingredients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    ingredient_id VARCHAR(36) NOT NULL,
    KEY ingredient_id_idx (ingredient_id),
    recipe_id VARCHAR(36) NOT NULL,
    KEY recipe_id_idx (recipe_id),
    quantity INT NOT NULL,
    units VARCHAR(32) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipe_ingredients;
-- +goose StatementEnd
