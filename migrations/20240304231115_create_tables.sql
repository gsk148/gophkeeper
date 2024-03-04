-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    password VARCHAR(255),
    UNIQUE(name),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS storage(
    id UUID DEFAULT gen_random_uuid(),
    uid UUID,
    data BYTEA,
    type INT,
    PRIMARY KEY(id),
    CONSTRAINT fk_user
    FOREIGN KEY (uid)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions(
    cid VARCHAR(50),
    token VARCHAR(165),
    PRIMARY KEY (cid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS storage;
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
