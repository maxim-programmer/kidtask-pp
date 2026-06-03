CREATE TABLE family_chat_messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id   UUID NOT NULL REFERENCES children(child_id) ON DELETE CASCADE,
    parent_id  UUID NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
    from_child BOOLEAN NOT NULL DEFAULT false,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_family_chat_child_id ON family_chat_messages(child_id);