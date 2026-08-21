-- +goose Up

CREATE TYPE event_status AS ENUM (
    'pending',
    'created',
    'failed'
);

ALTER TABLE polls
    ADD COLUMN status event_status NOT NULL DEFAULT 'created';

-- +goose Down

DROP TYPE event_status;

ALTER TABLE polls
    DROP COLUMN status;
