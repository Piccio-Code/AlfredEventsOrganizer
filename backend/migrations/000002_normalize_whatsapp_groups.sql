-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM poll_templates)
        OR EXISTS (SELECT 1 FROM polls)
        OR EXISTS (SELECT 1 FROM votes)
    THEN
        RAISE EXCEPTION
            'migration 000002 requires empty poll_templates, polls, and votes tables';
    END IF;
END
$$;
-- +goose StatementEnd

CREATE TABLE groups (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title            TEXT NOT NULL,
    whatsapp_chat_id TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT groups_title_not_blank
        CHECK (btrim(title) <> ''),
    CONSTRAINT groups_whatsapp_chat_id_not_blank
        CHECK (btrim(whatsapp_chat_id) <> ''),
    CONSTRAINT groups_whatsapp_chat_id_key
        UNIQUE (whatsapp_chat_id)
);

CREATE TABLE group_members (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id         UUID NOT NULL,
    name             TEXT NOT NULL,
    whatsapp_user_id TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT group_members_group_fkey
        FOREIGN KEY (group_id)
        REFERENCES groups (id)
        ON DELETE CASCADE,
    CONSTRAINT group_members_name_not_blank
        CHECK (btrim(name) <> ''),
    CONSTRAINT group_members_whatsapp_user_id_not_blank
        CHECK (btrim(whatsapp_user_id) <> ''),
    CONSTRAINT group_members_group_id_id_key
        UNIQUE (group_id, id),
    CONSTRAINT group_members_group_whatsapp_user_id_key
        UNIQUE (group_id, whatsapp_user_id)
);

ALTER TABLE poll_template_options
    ADD COLUMN motivation_needed BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN congratulation_needed BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE poll_options
    ADD COLUMN motivation_needed BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN congratulation_needed BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE poll_templates
    ADD COLUMN group_id UUID NOT NULL,
    ADD CONSTRAINT poll_templates_group_fkey
        FOREIGN KEY (group_id)
        REFERENCES groups (id);

ALTER TABLE poll_templates
    DROP CONSTRAINT poll_templates_whatsapp_chat_id_not_blank,
    DROP COLUMN whatsapp_chat_id;

ALTER TABLE polls
    ADD COLUMN group_id UUID NOT NULL,
    ADD CONSTRAINT polls_group_fkey
        FOREIGN KEY (group_id)
        REFERENCES groups (id),
    ADD CONSTRAINT polls_id_group_id_key
        UNIQUE (id, group_id);

ALTER TABLE polls
    DROP CONSTRAINT polls_whatsapp_chat_id_not_blank,
    DROP COLUMN whatsapp_chat_id;

ALTER TABLE votes
    DROP CONSTRAINT votes_pkey,
    DROP CONSTRAINT votes_whatsapp_user_id_not_blank,
    ADD COLUMN group_id UUID NOT NULL,
    ADD COLUMN group_member_id UUID NOT NULL;

ALTER TABLE votes
    DROP COLUMN whatsapp_user_id,
    ADD CONSTRAINT votes_pkey
        PRIMARY KEY (poll_id, group_member_id, option_id),
    ADD CONSTRAINT votes_poll_group_fkey
        FOREIGN KEY (poll_id, group_id)
        REFERENCES polls (id, group_id)
        ON DELETE CASCADE,
    ADD CONSTRAINT votes_group_member_fkey
        FOREIGN KEY (group_id, group_member_id)
        REFERENCES group_members (group_id, id)
        ON DELETE CASCADE;

-- +goose Down
ALTER TABLE votes
    ADD COLUMN whatsapp_user_id TEXT;

UPDATE votes
SET whatsapp_user_id = group_members.whatsapp_user_id
FROM group_members
WHERE group_members.group_id = votes.group_id
  AND group_members.id = votes.group_member_id;

ALTER TABLE votes
    ALTER COLUMN whatsapp_user_id SET NOT NULL,
    DROP CONSTRAINT votes_pkey,
    DROP CONSTRAINT votes_poll_group_fkey,
    DROP CONSTRAINT votes_group_member_fkey,
    DROP COLUMN group_member_id,
    DROP COLUMN group_id,
    ADD CONSTRAINT votes_pkey
        PRIMARY KEY (poll_id, whatsapp_user_id, option_id),
    ADD CONSTRAINT votes_whatsapp_user_id_not_blank
        CHECK (btrim(whatsapp_user_id) <> '');

ALTER TABLE polls
    ADD COLUMN whatsapp_chat_id TEXT;

UPDATE polls
SET whatsapp_chat_id = groups.whatsapp_chat_id
FROM groups
WHERE groups.id = polls.group_id;

ALTER TABLE polls
    ALTER COLUMN whatsapp_chat_id SET NOT NULL,
    DROP CONSTRAINT polls_id_group_id_key,
    DROP CONSTRAINT polls_group_fkey,
    DROP COLUMN group_id,
    ADD CONSTRAINT polls_whatsapp_chat_id_not_blank
        CHECK (btrim(whatsapp_chat_id) <> '');

ALTER TABLE poll_templates
    ADD COLUMN whatsapp_chat_id TEXT;

UPDATE poll_templates
SET whatsapp_chat_id = groups.whatsapp_chat_id
FROM groups
WHERE groups.id = poll_templates.group_id;

ALTER TABLE poll_templates
    ALTER COLUMN whatsapp_chat_id SET NOT NULL,
    DROP CONSTRAINT poll_templates_group_fkey,
    DROP COLUMN group_id,
    ADD CONSTRAINT poll_templates_whatsapp_chat_id_not_blank
        CHECK (btrim(whatsapp_chat_id) <> '');

ALTER TABLE poll_options
    DROP COLUMN motivation_needed,
    DROP COLUMN congratulation_needed;

ALTER TABLE poll_template_options
    DROP COLUMN motivation_needed,
    DROP COLUMN congratulation_needed;

DROP TABLE group_members;
DROP TABLE groups;
