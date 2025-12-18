package service

import (
	"github.com/axosec/vault/internal/data/db"
)

type VaultService struct {
	q *db.Queries
}

func NewVaultService(q *db.Queries) *VaultService {
	return &VaultService{
		q: q,
	}
}
