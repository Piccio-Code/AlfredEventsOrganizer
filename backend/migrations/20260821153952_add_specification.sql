-- +goose Up
ALTER TABLE poll_template_options
    ADD COLUMN specification_needed BOOLEAN DEFAULT false;

ALTER TABLE poll_options
    ADD COLUMN specification_needed BOOLEAN DEFAULT false;

-- +goose Down

ALTER TABLE poll_template_options
    DROP COLUMN specification_needed;

ALTER TABLE poll_options
    DROP COLUMN specification_needed;
