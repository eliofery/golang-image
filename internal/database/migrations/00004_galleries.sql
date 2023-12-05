-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS galleries(
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL,
    title VARCHAR(255) UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS galleries;
-- +goose StatementEnd
