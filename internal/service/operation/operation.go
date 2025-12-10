package operation

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	proto "general_ledger_golang/api/proto/code/go"
	"general_ledger_golang/domain"
	"general_ledger_golang/dto"
	"general_ledger_golang/internal/repository"
	"general_ledger_golang/models"
)

type OperationService struct {
	repo     repository.OperationRepository
	books    repository.BookRepository
	balances repository.BookBalanceRepository
	postings repository.PostingRepository
	db       *gorm.DB
}

func NewOperationService(db *gorm.DB, repo repository.OperationRepository, books repository.BookRepository, balances repository.BookBalanceRepository, postings repository.PostingRepository) *OperationService {
	if db == nil {
		db, _ = models.GetDB()
	}
	if repo == nil {
		repo = repository.NewOperationRepository(db)
	}
	if books == nil {
		books = repository.NewBookRepository(db)
	}
	if balances == nil {
		balances = repository.NewBookBalanceRepository(db)
	}
	if postings == nil {
		postings = repository.NewPostingRepository(db)
	}
	return &OperationService{
		repo:     repo,
		books:    books,
		balances: balances,
		postings: postings,
		db:       db,
	}
}

func (o *OperationService) GetOperation(ctx context.Context, memo string) (*dto.OperationDTO, error) {
	if memo == "" {
		return nil, errors.New("memo is empty")
	}
	foundOp, err := o.repo.GetByMemo(ctx, nil, memo)
	if err != nil || foundOp == nil {
		return nil, err
	}
	return dto.OperationToDTO(foundOp)
}

// PostOperation op-> {type: ”, memo:”, entries: [{bookId: ”, assetId: ”, value: ”}], metadata: {}}
func (o *OperationService) PostOperation(ctx context.Context, op dto.OperationPayload) (*dto.OperationDTO, error) {
	if err := op.Validate(); err != nil {
		return nil, err
	}
	return o.ApplyOperation(ctx, op)
}

func (o *OperationService) ApplyOperation(ctx context.Context, op dto.OperationPayload) (*dto.OperationDTO, error) {
	existingOp, err := o.GetOperation(ctx, op.Memo)
	if err != nil {
		return nil, err
	}
	if existingOp != nil {
		return existingOp, nil
	}

	entries := op.Entries
	var outDTO *dto.OperationDTO

	err = o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		opModel, err := dto.OperationPayloadToModel(op)
		if err != nil {
			return err
		}
		opModel.Status = string(models.OperationInit)

		newOp, err := o.applyOperationWithRetries(ctx, opModel, tx, 0)
		if err != nil {
			return err
		}

		bookIds := uniqueBookIds(entries)
		ok, bookErr := o.checkBooksExist(ctx, bookIds, tx)

		if !ok {
			reason := ""
			if bookErr != nil {
				reason = bookErr.Error()
			}
			if err := o.repo.UpdateStatus(ctx, tx, op.Memo, models.OperationRejected, reason); err != nil {
				return err
			}
			newOp.Status = string(models.OperationRejected)
			newOp.RejectionReason = reason
			outDTO, err = dto.OperationToDTO(newOp)
			return err
		}

		metadataBytes, err := json.Marshal(op.Metadata)
		if err != nil {
			return err
		}
		if err := o.postings.BulkCreate(ctx, tx, entries, newOp.Id, datatypes.JSON(metadataBytes)); err != nil {
			return err
		}

		if err := o.balances.ModifyBalance(ctx, tx, entries, op.Metadata); err != nil {
			return err
		}

		if err := o.repo.UpdateStatus(ctx, tx, op.Memo, models.OperationApplied, ""); err != nil {
			return err
		}
		newOp.Status = string(models.OperationApplied)
		outDTO, err = dto.OperationToDTO(newOp)
		return err
	})

	if err != nil {
		return nil, err
	}

	return outDTO, nil
}

func (o *OperationService) applyOperationWithRetries(ctx context.Context, op *models.Operation, tx *gorm.DB, rC int) (*models.Operation, error) {
	newOp, err := o.repo.Create(ctx, tx, op)
	totalRetryCount, currRetryCount := rC, 0

	for currRetryCount < totalRetryCount && err != nil {
		newOp, err = o.repo.Create(ctx, tx, op)
		currRetryCount++
	}

	if err != nil {
		return nil, err
	}
	return newOp, nil
}

func (o *OperationService) EntriesToProto(entries []domain.OperationEntry) []*proto.Entries {
	if len(entries) == 0 {
		return nil
	}
	result := make([]*proto.Entries, 0, len(entries))
	for _, entry := range entries {
		result = append(result, &proto.Entries{
			Value:   entry.Value,
			BookId:  entry.BookId,
			AssetId: entry.AssetId,
		})
	}
	return result
}

func (o *OperationService) ProtoEntriesToEntries(entries []*proto.Entries) []domain.OperationEntry {
	if len(entries) == 0 {
		return nil
	}
	result := make([]domain.OperationEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, domain.OperationEntry{
			BookId:  entry.BookId,
			AssetId: entry.AssetId,
			Value:   entry.Value,
		})
	}
	return result
}

func uniqueBookIds(entries []domain.OperationEntry) []string {
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(entries))
	for _, entry := range entries {
		if _, ok := seen[entry.BookId]; ok {
			continue
		}
		seen[entry.BookId] = struct{}{}
		unique = append(unique, entry.BookId)
	}
	return unique
}

func (o *OperationService) checkBooksExist(ctx context.Context, bookIds []string, tx *gorm.DB) (bool, error) {
	books, err := o.books.GetMany(ctx, tx, bookIds)
	if err != nil {
		return false, err
	}
	if len(books) != len(bookIds) {
		return false, errors.New("some bookIds couldn't be found")
	}
	return true, nil
}
