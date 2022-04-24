package models

import (
	"errors"

	"github.com/thoas/go-funk"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Book struct {
	Model
	Name     string         `gorm:"index;unique" json:"name"`
	Metadata datatypes.JSON `json:"metadata"`
}

func (b *Book) CreateOrUpdateBook(book *Book) (*gorm.DB, string) {
	var updateResult *gorm.DB
	if updateResult = db.Model(&book).Where("name = ?", book.Name).Updates(&book); updateResult.RowsAffected == 0 {
		return db.Create(&book), "create"
	}
	return updateResult, "update"
}

func (b *Book) GetBook(bookId string) (*Book, error) {
	book := Book{}
	q := db.Model(&b).Where("id = ?", bookId)

	res := q.Select("id", "name", "metadata", `createdAt`, `updatedAt`).Find(&book)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &book, nil
}

func (b *Book) GetBooks(bookIds []string, tx *gorm.DB) (*[]Book, error) {
	var books []Book
	var d *gorm.DB

	if tx != nil {
		d = tx
	} else {
		d = db
	}
	bookIdsUnique := funk.UniqString(bookIds)
	q := d.Model(&b).Where("id IN ?", bookIdsUnique)

	res := q.Select("id", "name", "metadata", `createdAt`, `updatedAt`).Find(&books)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &books, nil
}
