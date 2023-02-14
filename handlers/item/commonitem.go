package item

type CommonItemResponse struct {
	Id   string       `json:"id" binding:"required"`
	Kind ItemKindData `json:"kind" binding:"required"`
	CommonItemBaseData
}

type CommonItemCreateRequest struct {
	KindId string `json:"kindId" binding:"required"`
	CommonItemUpdateData
}

type CommonItemUpdateRequest struct {
	KindId string `json:"kindId" binding:"required"`
	CommonItemUpdateData
}

type CommonItemBaseData struct {
	Name        string `json:"name" binding:"required"`
	Priority    int    `json:"priority" binding:"required"`
	MaxOrder    int    `json:"maxOrder" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Enabled     *bool  `json:"enabled" binding:"required"`
	ImageUrl    *string `json:"imageUrl" binding:"required"`
}
type CommonItemUpdateData struct {
	Name        string `json:"name" binding:"required"`
	Priority    int    `json:"priority" binding:"required"`
	MaxOrder    int    `json:"maxOrder" binding:"required"`
	Price       *int    `json:"price" binding:"required,number"`
	Description string `json:"description" binding:"required"`
	Enabled     *bool  `json:"enabled" binding:"required"`
	ImageUrl    *string `json:"imageUrl" binding:"required"`
}