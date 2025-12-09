package domain

type OperationEntry struct {
	BookId  string `json:"bookId" validate:"required,min=1"`
	AssetId string `json:"assetId" validate:"required,min=1"`
	Value   string `json:"value" validate:"required"`
}
