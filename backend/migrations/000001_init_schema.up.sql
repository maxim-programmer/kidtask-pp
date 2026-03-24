CREATE TABLE parents (
    parent_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(100) NOT NULL
);

CREATE TABLE children (
    child_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id     UUID NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
    username      VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(100) NOT NULL,
    balance       INTEGER NOT NULL DEFAULT 0 CHECK (balance >= 0),
    age_group     VARCHAR(10) NOT NULL DEFAULT 'junior'
                      CHECK (age_group IN ('junior', 'senior')),
    UNIQUE (parent_id, username)
);

CREATE TABLE tasks (
    task_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id         UUID NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
    child_id          UUID NOT NULL REFERENCES children(child_id) ON DELETE CASCADE,
    title             VARCHAR(255) NOT NULL,
    description       TEXT,
    reward            INTEGER NOT NULL CHECK (reward > 0),
    status            VARCHAR(50) NOT NULL DEFAULT 'active'
                          CHECK (status IN ('active', 'pending_review', 'needs_rework', 'completed')),
    rejection_comment TEXT
);

CREATE TABLE wishes (
    wish_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id    UUID NOT NULL REFERENCES children(child_id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    price       INTEGER CHECK (price > 0),
    status      VARCHAR(50) NOT NULL DEFAULT 'awaiting_price'
                    CHECK (status IN ('awaiting_price', 'available', 'purchased', 'delivered'))
);

CREATE INDEX idx_children_parent_id ON children(parent_id);
CREATE INDEX idx_tasks_child_id     ON tasks(child_id);
CREATE INDEX idx_tasks_parent_id    ON tasks(parent_id);
CREATE INDEX idx_tasks_status       ON tasks(status);
CREATE INDEX idx_wishes_child_id    ON wishes(child_id);
CREATE INDEX idx_wishes_status      ON wishes(status);