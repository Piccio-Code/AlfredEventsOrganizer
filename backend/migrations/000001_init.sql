-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE poll_templates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name             TEXT NOT NULL,
    title            TEXT NOT NULL,
    whatsapp_chat_id TEXT NOT NULL,
    multiple_choice  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT poll_templates_name_not_blank
        CHECK (btrim(name) <> ''),
    CONSTRAINT poll_templates_title_not_blank
        CHECK (btrim(title) <> ''),
    CONSTRAINT poll_templates_whatsapp_chat_id_not_blank
        CHECK (btrim(whatsapp_chat_id) <> ''),
    CONSTRAINT poll_templates_name_key
        UNIQUE (name)
);

CREATE TABLE poll_template_options (
    id          UUID NOT NULL DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL,
    label       TEXT NOT NULL,
    position    INTEGER NOT NULL,

    CONSTRAINT poll_template_options_pkey
        PRIMARY KEY (template_id, id),
    CONSTRAINT poll_template_options_template_fkey
        FOREIGN KEY (template_id)
        REFERENCES poll_templates (id)
        ON DELETE CASCADE,
    CONSTRAINT poll_template_options_label_not_blank
        CHECK (btrim(label) <> ''),
    CONSTRAINT poll_template_options_position_nonnegative
        CHECK (position >= 0),
    CONSTRAINT poll_template_options_position_key
        UNIQUE (template_id, position),
    CONSTRAINT poll_template_options_label_key
        UNIQUE (template_id, label)
);

CREATE TABLE polls (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id      UUID,
    whatsapp_poll_id TEXT,
    whatsapp_chat_id TEXT NOT NULL,
    title            TEXT NOT NULL,
    multiple_choice  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at       TIMESTAMPTZ,

    CONSTRAINT polls_template_fkey
        FOREIGN KEY (template_id)
        REFERENCES poll_templates (id)
        ON DELETE SET NULL,
    CONSTRAINT polls_whatsapp_poll_id_key
        UNIQUE (whatsapp_poll_id),
    CONSTRAINT polls_whatsapp_chat_id_not_blank
        CHECK (btrim(whatsapp_chat_id) <> ''),
    CONSTRAINT polls_whatsapp_poll_id_not_blank
        CHECK (
            whatsapp_poll_id IS NULL
            OR btrim(whatsapp_poll_id) <> ''
        ),
    CONSTRAINT polls_title_not_blank
        CHECK (btrim(title) <> '')
);

CREATE TABLE poll_options (
    id       UUID NOT NULL DEFAULT gen_random_uuid(),
    poll_id  UUID NOT NULL,
    label    TEXT NOT NULL,
    position INTEGER NOT NULL,

    CONSTRAINT poll_options_pkey
        PRIMARY KEY (poll_id, id),
    CONSTRAINT poll_options_poll_fkey
        FOREIGN KEY (poll_id)
        REFERENCES polls (id)
        ON DELETE CASCADE,
    CONSTRAINT poll_options_label_not_blank
        CHECK (btrim(label) <> ''),
    CONSTRAINT poll_options_position_nonnegative
        CHECK (position >= 0),
    CONSTRAINT poll_options_position_key
        UNIQUE (poll_id, position),
    CONSTRAINT poll_options_label_key
        UNIQUE (poll_id, label)
);

CREATE TABLE votes (
    poll_id          UUID NOT NULL,
    option_id        UUID NOT NULL,
    whatsapp_user_id TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT votes_pkey
        PRIMARY KEY (poll_id, whatsapp_user_id, option_id),
    CONSTRAINT votes_poll_option_fkey
        FOREIGN KEY (poll_id, option_id)
        REFERENCES poll_options (poll_id, id)
        ON DELETE CASCADE,
    CONSTRAINT votes_whatsapp_user_id_not_blank
        CHECK (btrim(whatsapp_user_id) <> '')
);

CREATE INDEX polls_template_id_idx
    ON polls (template_id)
    WHERE template_id IS NOT NULL;

CREATE INDEX polls_expires_at_idx
    ON polls (expires_at)
    WHERE expires_at IS NOT NULL;

CREATE INDEX votes_poll_option_idx
    ON votes (poll_id, option_id);

-- +goose Down
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS poll_options;
DROP TABLE IF EXISTS polls;
DROP TABLE IF EXISTS poll_template_options;
DROP TABLE IF EXISTS poll_templates;
