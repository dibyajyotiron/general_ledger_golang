package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Book struct {
	Model
	Name         string         `json:"name"`
	Metadata     datatypes.JSON `json:"metadata"`
	Restrictions datatypes.JSON `json:"restrictions"`
}

func (b *Book) CreateBook(data *Book) *gorm.DB {
	book := Book{Name: data.Name, Metadata: data.Metadata, Restrictions: data.Restrictions}
	result := db.Create(&book)
	return result
}
