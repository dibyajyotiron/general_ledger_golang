package dto

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"

	"general_ledger_golang/domain"
	"general_ledger_golang/models"
)

type OperationPayload struct {
	Type     string                  `json:"type" validate:"required,min=3,max=20"`
	Memo     string                  `json:"memo" validate:"required,min=3"`
	Entries  []domain.OperationEntry `json:"entries" validate:"required,dive"`
	Metadata map[string]interface{}  `json:"metadata" validate:"required"`
}

type OperationDTO struct {
	Id              uint64                  `json:"id"`
	Type            string                  `json:"type"`
	Memo            string                  `json:"memo"`
	Entries         []domain.OperationEntry `json:"entries"`
	Status          string                  `json:"status"`
	RejectionReason string                  `json:"rejectionReason"`
	Metadata        map[string]interface{}  `json:"metadata"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
}

func (p OperationPayload) Validate() error {
	v := validator.New()
	return v.Struct(p)
}

func OperationPayloadToModel(p OperationPayload) (*models.Operation, error) {
	entriesBytes, err := json.Marshal(p.Entries)
	if err != nil {
		return nil, err
	}
	metadataBytes, err := json.Marshal(p.Metadata)
	if err != nil {
		return nil, err
	}
	return &models.Operation{
		Type:     p.Type,
		Memo:     p.Memo,
		Entries:  datatypes.JSON(entriesBytes),
		Metadata: datatypes.JSON(metadataBytes),
	}, nil
}

func OperationToDTO(o *models.Operation) (*OperationDTO, error) {
	var entries []domain.OperationEntry
	if len(o.Entries) > 0 {
		if err := json.Unmarshal(o.Entries, &entries); err != nil {
			return nil, err
		}
	}
	metadata := map[string]interface{}{}
	if len(o.Metadata) > 0 {
		if err := json.Unmarshal(o.Metadata, &metadata); err != nil {
			return nil, err
		}
	}
	return &OperationDTO{
		Id:              o.Id,
		Type:            o.Type,
		Memo:            o.Memo,
		Entries:         entries,
		Status:          o.Status,
		RejectionReason: o.RejectionReason,
		Metadata:        metadata,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}, nil
}
