ALTER TABLE parents ADD COLUMN IF NOT EXISTS is_blocked BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE children ADD COLUMN IF NOT EXISTS is_blocked BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE balance_logs (
    log_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id   UUID NOT NULL REFERENCES children(child_id) ON DELETE CASCADE,
    delta      INTEGER NOT NULL,
    reason     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE complaints (
    complaint_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id    UUID NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
    subject      VARCHAR(255) NOT NULL,
    body         TEXT NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'open'
                     CHECK (status IN ('open', 'resolved')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE chat_messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id  UUID NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
    from_admin BOOLEAN NOT NULL DEFAULT false,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_balance_logs_child_id ON balance_logs(child_id);
CREATE INDEX idx_complaints_parent_id  ON complaints(parent_id);
CREATE INDEX idx_chat_messages_parent  ON chat_messages(parent_id);