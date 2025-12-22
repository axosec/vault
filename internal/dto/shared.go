package dto

import "github.com/google/uuid"

type ResourceType string

const (
	TypeFolder ResourceType = "FOLDER"
	TypeItem   ResourceType = "ITEM"
)

type ShareParams struct {
	TargetUserID uuid.UUID    `json:"target_user_id" binding:"required"`
	ResourceID   uuid.UUID    `json:"resource_id" binding:"required"`
	ResourceType ResourceType `json:"resource_type" binding:"required,oneof=FOLDER ITEM"`

	EncKey      []byte `json:"enc_key" binding:"required"`
	KeyNonce    []byte `json:"key_nonce" binding:"required"`
	AccessLevel string `json:"access_level" binding:"required,oneof=READ WRITE OWNER"`
}
