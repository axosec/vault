package dto

import (
	"github.com/google/uuid"
)

type CreateFolderReq struct {
	ParentID *uuid.UUID `json:"parent_id"`

	EncName   []byte `json:"enc_name" binding:"required"`
	NameNonce []byte `json:"name_nonce" binding:"required"`

	EncKey   []byte `json:"enc_key" binding:"required"`
	KeyNonce []byte `json:"key_nonce" binding:"required"`
}

type FolderResponse struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id"`
}

type FolderSummary struct {
	ID           uuid.UUID  `json:"id"`
	ParentID     *uuid.UUID `json:"parent_id"`
	EncName      []byte     `json:"enc_name"`
	EncNameNonce []byte     `json:"enc_name_nonce"`

	WrappedKey []byte `json:"wrapped_key"`
	KeyNonce   []byte `json:"key_nonce"`
}
