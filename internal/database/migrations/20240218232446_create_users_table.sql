-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id              VARCHAR(40)     PRIMARY KEY NOT NULL,
    external_id     VARCHAR(64)     NOT NULL,
    username        VARCHAR(255)    NOT NULL,
    email           VARCHAR(255)    NOT NULL,
    created_at      TIMESTAMP(6)    NOT NULL,
    updated_at      TIMESTAMP(6)
);

CREATE UNIQUE INDEX uidx_users_username ON users (username);
CREATE UNIQUE INDEX uidx_users_email ON users (email);
CREATE UNIQUE INDEX uidx_users_external_id ON users (external_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uidx_users_username;
DROP INDEX uidx_users_email;
DROP INDEX uidx_users_external_id;
DROP TABLE users;
-- +goose StatementEnd
