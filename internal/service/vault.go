package service

import (
	"context"
	"fmt"

	"github.com/axosec/vault/internal/data/db"
	"github.com/axosec/vault/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VaultService struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewVaultService(pool *pgxpool.Pool, q *db.Queries) *VaultService {
	return &VaultService{
		pool: pool,
		q:    q,
	}
}

func (s *VaultService) CreateFolder(ctx context.Context, userID uuid.UUID, req dto.CreateFolderReq) (*dto.FolderResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	folder, err := qtx.CreateFolder(ctx, db.CreateFolderParams{
		OwnerID:  userID,
		ParentID: req.ParentID,
		Nonce:    req.NameNonce,
		EncName:  req.EncName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	err = qtx.CreateFolderKey(ctx, db.CreateFolderKeyParams{
		UserID:   userID,
		FolderID: &folder.ID,

		EncKey:      req.EncKey,
		Nonce:       req.KeyNonce,
		AccessLevel: "OWNER",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create folder key: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return &dto.FolderResponse{
		ID:       folder.ID,
		ParentID: req.ParentID,
	}, nil
}

func (s *VaultService) ListFolders(ctx context.Context, userID uuid.UUID) ([]dto.FolderSummary, error) {
	foldersDb, err := s.q.GetUserFolders(ctx, userID)
	if err != nil {
		return nil, err
	}

	folders := make([]dto.FolderSummary, len(foldersDb))

	for i, folder := range foldersDb {
		folders[i].ID = folder.ID
		folders[i].EncName = folder.EncName
		folders[i].EncNameNonce = folder.Nonce
		folders[i].KeyNonce = folder.KeyNonce
		folders[i].WrappedKey = folder.WrappedKey
		folders[i].ParentID = folder.ParentID
	}

	return folders, nil
}

func (s *VaultService) CreateItem(ctx context.Context, userID uuid.UUID, req dto.CreateItemReq) (*dto.ItemResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	var folderIDPtr *uuid.UUID
	if req.FolderID != uuid.Nil {
		folderIDPtr = &req.FolderID
	}

	item, err := qtx.CreateItem(ctx, db.CreateItemParams{
		OwnerID:     userID,
		FolderID:    folderIDPtr,
		Type:        req.Type,
		Nonce:       req.DataNonce,
		EncData:     req.EncData,
		EncOverview: req.EncOverview,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	err = qtx.CreateItemKey(ctx, db.CreateItemKeyParams{
		UserID: userID,
		ItemID: &item.ID,

		EncKey:      req.EncKey,
		Nonce:       req.KeyNonce,
		AccessLevel: "OWNER",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create item key: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return &dto.ItemResponse{
		ID: item.ID,
	}, nil
}

