-- +goose Up
-- +goose StatementBegin
ALTER TABLE ingredients 
	ADD COLUMN amount DOUBLE NOT NULL;
ALTER TABLE ingredients 
	ADD COLUMN units VARCHAR(32) NOT NULL;
ALTER TABLE ingredients 
	ADD COLUMN standard_amount DOUBLE NOT NULL;
ALTER TABLE ingredients 
	ADD COLUMN standard_units VARCHAR(32) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ingredients DROP amount;
ALTER TABLE ingredients DROP units;
ALTER TABLE ingredients DROP standard_amount;
ALTER TABLE ingredients DROP standard_units;
-- +goose StatementEnd
