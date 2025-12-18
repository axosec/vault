-- name: CreateFolder :one
INSERT INTO folders (owner_id, parent_id, nonce, enc_name)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetUserFolders :many
SELECT
    f.id,
    f.parent_id,
    f.nonce,
    f.enc_name,
    f.updated_at,
    k.enc_key AS wrapped_key,
    k.nonce AS key_nonce,
    k.access_level
FROM folders f
JOIN keys k ON f.id = k.folder_id
WHERE k.user_id = $1
  AND f.deleted_at IS NULL
ORDER BY f.created_at ASC;

-- name: SoftDeleteFolder :exec
UPDATE folders
SET deleted_at = NOW()
WHERE id = $1 AND owner_id = $2;


-- name: CreateItem :one
INSERT INTO items (owner_id, folder_id, type, nonce, enc_data, enc_overview)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at;

-- name: GetFolderItems :many
SELECT
    i.id,
    i.type,
    i.nonce AS item_nonce,
    i.enc_overview,
    i.updated_at,
    k.enc_key AS wrapped_key,
    k.nonce AS key_nonce
FROM items i
JOIN keys k ON i.id = k.item_id
WHERE i.folder_id = $1
  AND k.user_id = $2
  AND i.deleted_at IS NULL
ORDER BY i.created_at DESC;

-- name: GetItemData :one
SELECT
    i.id,
    i.folder_id,
    i.type,
    i.nonce AS item_nonce,
    i.enc_data,
    i.enc_overview,
    i.created_at,
    i.updated_at,
    k.enc_key AS wrapped_key,
    k.nonce AS key_nonce,
    k.access_level
FROM items i
JOIN keys k ON i.id = k.item_id
WHERE i.id = $1
  AND k.user_id = $2
  AND i.deleted_at IS NULL;

-- name: SoftDeleteItem :exec
UPDATE items
SET deleted_at = NOW()
WHERE id = $1 AND owner_id = $2;

-- name: CreateFolderKey :exec
INSERT INTO keys (user_id, folder_id, enc_key, nonce, access_level)
VALUES ($1, $2, $3, $4, $5);

-- name: CreateItemKey :exec
INSERT INTO keys (user_id, item_id, enc_key, nonce, access_level)
VALUES ($1, $2, $3, $4, $5);

-- name: RevokeUserAccess :exec
DELETE FROM keys
WHERE user_id = $1
  AND (folder_id = $2 OR item_id = $2);
