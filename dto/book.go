package dto

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"

	"general_ledger_golang/models"
)

type BookPayload struct {
	Name     string         `json:"name" validate:"required,min=1"`
	Metadata map[string]any `json:"metadata"`
}

type BookDTO struct {
	Id        uint64         `json:"id"`
	Name      string         `json:"name"`
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (p BookPayload) Validate() error {
	v := validator.New()
	return v.Struct(p)
}

func BookPayloadToModel(p BookPayload) (*models.Book, error) {
	metadataBytes, err := json.Marshal(p.Metadata)
	if err != nil {
		return nil, err
	}
	return &models.Book{Name: p.Name, Metadata: datatypes.JSON(metadataBytes)}, nil
}

func BookToDTO(b *models.Book) (*BookDTO, error) {
	metadata := map[string]any{}
	if len(b.Metadata) > 0 {
		if err := json.Unmarshal(b.Metadata, &metadata); err != nil {
			return nil, err
		}
	}
	return &BookDTO{
		Id:        b.Id,
		Name:      b.Name,
		Metadata:  metadata,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}, nil
}
