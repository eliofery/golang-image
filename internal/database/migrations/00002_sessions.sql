-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions
(
    id         SERIAL PRIMARY KEY,
    user_id    INT UNIQUE          NOT NULL,
    token_hash VARCHAR(256) UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
