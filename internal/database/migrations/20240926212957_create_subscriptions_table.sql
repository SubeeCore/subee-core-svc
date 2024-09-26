-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
    id              VARCHAR(40)         PRIMARY KEY NOT NULL,
    user_id         VARCHAR(40)         NOT NULL,
    platform        VARCHAR(255)        NOT NULL,
    reccurence      INTEGER             NOT NULL,
    price           DOUBLE PRECISION    NOT NULL,
    started_at      TIMESTAMP(6)        NOT NULL,
    created_at      TIMESTAMP(6)        NOT NULL,
    finished_at     TIMESTAMP(6),
    updated_at      TIMESTAMP(6)
);

CREATE INDEX uidx_subscriptions_user_id ON subscriptions (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uidx_subscriptions_user_id;
DROP TABLE subscriptions;
-- +goose StatementEnd
