-- +goose Up
-- +goose StatementBegin
ALTER TABLE items 
	ADD COLUMN name TEXT NOT NULL,
	ADD COLUMN description TEXT,
	ADD COLUMN amount DOUBLE NOT NULL, 
	ADD COLUMN units VARCHAR(32) NOT NULL, 
	ADD COLUMN standard_amount DOUBLE NOT NULL, 
	ADD COLUMN standard_units VARCHAR(32) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items DROP name, description, amount, units, standard_amount, standard_units;
-- +goose StatementEnd
