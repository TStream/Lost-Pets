-- +goose Up
-- +goose StatementBegin
INSERT INTO types(name)
	VALUES ('dog'),('cat');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM types WHERE name = 'dog' OR name = 'cat'
-- +goose StatementEnd
