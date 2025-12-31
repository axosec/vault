package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateItemReq struct {
	FolderID uuid.UUID `json:"folder_id" binding:"required"`
	Type     string    `json:"type" binding:"required"`

	EncData       []byte `json:"enc_data" binding:"required"`
	EncOverview   []byte `json:"enc_overview"`
	DataNonce     []byte `json:"data_nonce" binding:"required"`
	OverviewNonce []byte `json:"overview_nonce" binding:"required"`

	EncKey   []byte `json:"enc_key" binding:"required"`
	KeyNonce []byte `json:"key_nonce" binding:"required"`
}

type UpdateItemReq struct {
	EncData       []byte `json:"enc_data" binding:"required"`
	EncOverview   []byte `json:"enc_overview"`
	Nonce         []byte `json:"data_nonce" binding:"required"`
	OverviewNonce []byte `json:"overview_nonce" binding:"required"`
}

type ItemResponse struct {
	ID uuid.UUID `json:"id"`
}

type ItemSummary struct {
	ID            uuid.UUID `json:"id"`
	Type          string    `json:"type"`
	EncOverview   []byte    `json:"enc_overview"`
	OverviewNonce []byte    `json:"overview_nonce"`
	WrappedKey    []byte    `json:"wrapped_key"`
	KeyNonce      []byte    `json:"key_nonce"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ItemDetail struct {
	ID         uuid.UUID  `json:"id"`
	FolderID   *uuid.UUID `json:"folder_id"`
	Type       string     `json:"type"`
	EncData    []byte     `json:"enc_data"`
	DataNonce  []byte     `json:"data_nonce"`
	WrappedKey []byte     `json:"wrapped_key"`
	KeyNonce   []byte     `json:"key_nonce"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
