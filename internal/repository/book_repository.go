package repository

import (
	"context"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"general_ledger_golang/models"
)

type bookGormRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookGormRepository{db: db}
}

func (r *bookGormRepository) Upsert(ctx context.Context, tx *gorm.DB, book *models.Book) (*models.Book, string, error) {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}
	operation := "update"
	updateResult := db.Model(&models.Book{}).Where("name = ?", book.Name).Clauses(clause.Returning{}).Updates(book).Scan(&book)
	if updateResult.Error != nil {
		return nil, "", updateResult.Error
	}

	if updateResult.RowsAffected == 0 {
		createResult := db.Create(book)
		if createResult.Error != nil {
			return nil, "", createResult.Error
		}
		operation = "create"
	}
	return book, operation, nil
}

func (r *bookGormRepository) GetByID(ctx context.Context, tx *gorm.DB, bookId string) (*models.Book, error) {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}
	book := models.Book{}
	res := db.Model(&models.Book{}).Where("id = ?", bookId).Select("id", "name", "metadata", "createdAt", "updatedAt").First(&book)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, res.Error
	}
	return &book, nil
}

func (r *bookGormRepository) GetMany(ctx context.Context, tx *gorm.DB, bookIds []string) ([]models.Book, error) {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}
	unique := funk.UniqString(bookIds)
	var books []models.Book
	res := db.Model(&models.Book{}).Where("id IN ?", unique).Select("id", "name", "metadata", "createdAt", "updatedAt").Find(&books)
	if res.Error != nil {
		return nil, res.Error
	}
	return books, nil
}
