package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
	"errors"
)

type TransactionService struct {
	repo      repository.TransactionRepository
	assetRepo repository.AssetRepository
}

func NewTransactionService(repo repository.TransactionRepository, assetRepo repository.AssetRepository) *TransactionService {
	return &TransactionService{
		repo: repo, assetRepo: assetRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, transaction *model.Transaction) error {
	if transaction.UserID <= 0 {
		return errors.New("user id must be greater than 0")
	}

	if transaction.AssetID <= 0 {
		return errors.New("asset id must be greater than 0")
	}

	if transaction.Type != model.TransactionBuy && transaction.Type != model.TransactionSell {
		return errors.New("invalid transaction type")
	}

	if transaction.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if transaction.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if transaction.Commission < 0 {
		return errors.New("commission cannot be negative")
	}

	_, err := s.assetRepo.FindByID(ctx, transaction.AssetID)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, transaction)
}

func (s *TransactionService) FindByID(ctx context.Context, id int64) (*model.Transaction, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TransactionService) FindByIDAndUserID(ctx context.Context, transactionID int64, userID int64) (*model.Transaction, error) {
	return s.repo.FindByIDAndUserID(ctx, transactionID, userID)
}

func (s *TransactionService) FindByUserID(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	return s.repo.FindByUserID(ctx, userID)
}
