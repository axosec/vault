CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,

    nonce BYTEA NOT NULL,
    enc_metadata BYTEA NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    folder_id UUID REFERENCES folders(id),

    type VARCHAR(50) NOT NULL,

    nonce BYTEA NOT NULL,
    enc_data BYTEA NOT NULL,

    overview_nonce BYTEA,
    enc_overview BYTEA,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,

    folder_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    item_id   UUID REFERENCES items(id)   ON DELETE CASCADE,

    CONSTRAINT check_resource_target
        CHECK (
            (folder_id IS NOT NULL AND item_id IS NULL) OR
            (folder_id IS NULL AND item_id IS NOT NULL)
        ),

    enc_key BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    access_level VARCHAR(20) NOT NULL DEFAULT 'READ',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_folders_owner ON folders(owner_id);
CREATE INDEX idx_items_folder ON items(folder_id);
CREATE INDEX idx_keys_user ON keys(user_id);
CREATE INDEX idx_keys_folder ON keys(folder_id);
CREATE INDEX idx_keys_item ON keys(item_id);
