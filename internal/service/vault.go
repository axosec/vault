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
		OwnerID:     userID,
		Nonce:       req.NameNonce,
		EncMetadata: req.EncMetadata,
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
		ID: folder.ID,
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
		folders[i].EncMetadata = folder.EncMetadata
		folders[i].Nonce = folder.Nonce
		folders[i].KeyNonce = folder.KeyNonce
		folders[i].WrappedKey = folder.WrappedKey
	}

	return folders, nil
}

func (s *VaultService) UpdateFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req dto.UpdateFolderReq) error {
	rowsAffected, err := s.q.UpdateFolderMetadata(ctx, db.UpdateFolderMetadataParams{
		EncMetadata: req.EncMetadata,
		Nonce:       req.MetadataNonce,
		ID:          folderID,
		OwnerID:     userID,
	})

	if err != nil {
		return fmt.Errorf("failed to update folder: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("folder not found or access denied")
	}

	return nil
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

func (s *VaultService) ListItems(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) ([]dto.ItemSummary, error) {
	var folderIDPtr *uuid.UUID
	if folderID != uuid.Nil {
		folderIDPtr = &folderID
	}

	itemsDb, err := s.q.GetFolderItems(ctx, db.GetFolderItemsParams{
		FolderID: folderIDPtr,
		UserID:   userID,
	})
	if err != nil {
		return nil, err
	}

	items := make([]dto.ItemSummary, len(itemsDb))
	for i, item := range itemsDb {
		items[i] = dto.ItemSummary{
			ID:          item.ID,
			Type:        item.Type,
			EncOverview: item.EncOverview,
			WrappedKey:  item.WrappedKey,
			KeyNonce:    item.KeyNonce,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return items, nil
}

func (s *VaultService) GetItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*dto.ItemDetail, error) {
	item, err := s.q.GetItemData(ctx, db.GetItemDataParams{
		ID:     itemID,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("item not found or access denied: %w", err)
	}

	return &dto.ItemDetail{
		ID:         item.ID,
		FolderID:   item.FolderID,
		Type:       item.Type,
		EncData:    item.EncData,
		DataNonce:  item.ItemNonce,
		WrappedKey: item.WrappedKey,
		KeyNonce:   item.KeyNonce,
		UpdatedAt:  item.UpdatedAt,
	}, nil
}

func (s *VaultService) UpdateItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req dto.UpdateItemReq) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	meta, err := qtx.GetItemData(ctx, db.GetItemDataParams{
		ID:     itemID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("cannot access item: %w", err)
	}

	canWrite := false
	if meta.AccessLevel == "OWNER" || meta.AccessLevel == "WRITE" {
		canWrite = true
	}

	if !canWrite {
		return fmt.Errorf("access denied: requires WRITE or OWNER permission")
	}

	err = qtx.UpdateItemBlob(ctx, db.UpdateItemBlobParams{
		EncData:     req.EncData,
		EncOverview: req.EncOverview,
		Nonce:       req.Nonce,
		ID:          itemID,
	})
	if err != nil {
		return fmt.Errorf("failed to update item blob: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (s *VaultService) DeleteResource(ctx context.Context, userID uuid.UUID, resourceID uuid.UUID, resourceType dto.ResourceType) error {
	var rowsAffected int64
	var err error

	switch resourceType {
	case dto.TypeFolder:
		rowsAffected, err = s.q.SoftDeleteFolder(ctx, db.SoftDeleteFolderParams{
			ID:      resourceID,
			OwnerID: userID,
		})

	case dto.TypeItem:
		rowsAffected, err = s.q.SoftDeleteItem(ctx, db.SoftDeleteItemParams{
			ID:      resourceID,
			OwnerID: userID,
		})

	default:
		return fmt.Errorf("invalid resource type")
	}

	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("resource not found or access denied (not owner)")
	}

	return nil
}

func (s *VaultService) ShareResource(ctx context.Context, ownerID uuid.UUID, req dto.ShareParams) error {
	var err error
	if req.ResourceType == dto.TypeFolder {
		_, err = s.q.IsFolderOwner(ctx, db.IsFolderOwnerParams{
			ID:      req.ResourceID,
			OwnerID: ownerID,
		})
	} else {
		_, err = s.q.IsItemOwner(ctx, db.IsItemOwnerParams{
			ID:      req.ResourceID,
			OwnerID: ownerID,
		})
	}

	if err != nil {
		return fmt.Errorf("resource not found or access denied")
	}

	switch req.ResourceType {
	case dto.TypeFolder:
		err := s.q.CreateFolderKey(ctx, db.CreateFolderKeyParams{
			UserID:      req.TargetUserID,
			FolderID:    &req.ResourceID,
			EncKey:      req.EncKey,
			Nonce:       req.KeyNonce,
			AccessLevel: req.AccessLevel,
		})
		if err != nil {
			return fmt.Errorf("failed to share folder: %w", err)
		}

	case dto.TypeItem:
		err := s.q.CreateItemKey(ctx, db.CreateItemKeyParams{
			UserID:      req.TargetUserID,
			ItemID:      &req.ResourceID,
			EncKey:      req.EncKey,
			Nonce:       req.KeyNonce,
			AccessLevel: req.AccessLevel,
		})
		if err != nil {
			return fmt.Errorf("failed to share item: %w", err)
		}

	default:
		return fmt.Errorf("invalid resource type")
	}

	return nil
}

func (s *VaultService) RevokeAccess(ctx context.Context, ownerID uuid.UUID, targetUserID uuid.UUID, resourceID uuid.UUID) error {
	isOwner := false

	_, err := s.q.IsFolderOwner(ctx, db.IsFolderOwnerParams{
		ID:      resourceID,
		OwnerID: ownerID,
	})
	if err == nil {
		isOwner = true
	} else {
		_, err := s.q.IsItemOwner(ctx, db.IsItemOwnerParams{
			ID:      resourceID,
			OwnerID: ownerID,
		})
		if err == nil {
			isOwner = true
		}
	}

	if !isOwner {
		return fmt.Errorf("access denied: you are not the owner of this resource")
	}

	err = s.q.RevokeUserAccess(ctx, db.RevokeUserAccessParams{
		UserID:   targetUserID,
		FolderID: &resourceID,
	})

	if err != nil {
		return fmt.Errorf("failed to revoke access: %w", err)
	}

	return nil
}
