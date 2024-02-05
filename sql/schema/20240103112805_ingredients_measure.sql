-- +goose Up
-- +goose StatementBegin
ALTER TABLE ingredients 
	ADD COLUMN amount DOUBLE NOT NULL, 
	ADD COLUMN units VARCHAR(32) NOT NULL, 
	ADD COLUMN standard_amount DOUBLE NOT NULL, 
	ADD COLUMN standard_units VARCHAR(32) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ingredients DROP amount, units, standard_amount, standard_units;
-- +goose StatementEnd
