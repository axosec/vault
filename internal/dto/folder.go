package dto

import (
	"github.com/google/uuid"
)

type CreateFolderReq struct {
	EncMetadata []byte `json:"enc_metadata" binding:"required"`
	NameNonce   []byte `json:"nonce" binding:"required"`

	EncKey   []byte `json:"enc_key" binding:"required"`
	KeyNonce []byte `json:"key_nonce" binding:"required"`
}

type FolderResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateFolderReq struct {
	EncMetadata   []byte `json:"enc_metadata" binding:"required"`
	MetadataNonce []byte `json:"nonce" binding:"required"`
}

type FolderSummary struct {
	ID          uuid.UUID `json:"id"`
	EncMetadata []byte    `json:"enc_metadata"`
	Nonce       []byte    `json:"nonce"`

	WrappedKey []byte `json:"wrapped_key"`
	KeyNonce   []byte `json:"key_nonce"`
}
