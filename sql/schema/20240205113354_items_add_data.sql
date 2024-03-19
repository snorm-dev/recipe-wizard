-- +goose Up
-- +goose StatementBegin
ALTER TABLE items 
	ADD COLUMN name TEXT NOT NULL;
ALTER TABLE items 
	ADD COLUMN description TEXT;
ALTER TABLE items 
	ADD COLUMN amount DOUBLE NOT NULL;
ALTER TABLE items 
	ADD COLUMN units VARCHAR(32) NOT NULL;
ALTER TABLE items 
	ADD COLUMN standard_amount DOUBLE NOT NULL;
ALTER TABLE items 
	ADD COLUMN standard_units VARCHAR(32) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items DROP name;
ALTER TABLE items DROP description;
ALTER TABLE items DROP amount;
ALTER TABLE items DROP units;
ALTER TABLE items DROP standard_amount;
ALTER TABLE items DROP standard_units;
-- +goose StatementEnd
