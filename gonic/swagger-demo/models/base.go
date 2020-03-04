package models

// GetModel ...
// swagger:model GetModel
type GetModel struct {
	ModelName string `json:"name"`
	Get       bool   `json:"is_get"`
}

// PostModel ...
// swagger:model PostModel
type PostModel struct {
	ModelName string `json:"name"`
	Post      bool   `json:"is_post"`
}
