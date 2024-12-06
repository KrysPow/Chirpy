-- +goose UP
ALTER TABLE users 
    ADD COLUMN hashed_password TEXT NOT NULL 
    CONSTRAINT check_hashed_password CHECK (hashed_password <> '')
    DEFAULT 'unset';

-- +goose Down
ALTER TABLE users DROP COLUMN hashed_password;