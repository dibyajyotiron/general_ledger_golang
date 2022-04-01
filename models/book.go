package models

import (
	"errors"
	"general_ledger_golang/pkg/util"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Book struct {
	Model
	Name         string         `gorm:"index;unique" json:"name"`
	Metadata     datatypes.JSON `json:"metadata"`
	Restrictions datatypes.JSON `json:"restrictions"`
}

func (b *Book) CreateOrUpdateBook(book *Book) (*gorm.DB, string) {
	var updateResult *gorm.DB
	if updateResult = db.Model(&book).Where("name = ?", book.Name).Updates(&book); updateResult.RowsAffected == 0 {
		return db.Create(&book), "create"
	}
	return updateResult, "update"
}

func (b *Book) GetBook(bookId string) *map[string]interface{} {
	book := Book{}
	res := db.Model(&b).Where("id = ?", bookId).Last(&book)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	result := util.StructToJSON(book)
	return &result
}
