-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes 
RENAME COLUMN cool_time 
TO total_time;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes 
RENAME COLUMN total_time 
TO cool_time;
-- +goose StatementEnd
