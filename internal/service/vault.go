package service

import (
	"context"

	"github.com/axosec/vault/internal/data/db"
	"github.com/google/uuid"
)

type VaultService struct {
	q *db.Queries
}

func NewVaultService(q *db.Queries) *VaultService {
	return &VaultService{
		q: q,
	}
}
