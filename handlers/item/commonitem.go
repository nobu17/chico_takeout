package item

type CommonItemResponse struct {
	Id   string       `json:"id" binding:"required"`
	Kind ItemKindData `json:"kind" binding:"required"`
	CommonItemBaseData
}

type CommonItemCreateRequest struct {
	KindId string `json:"kindId" binding:"required"`
	CommonItemBaseData
}

type CommonItemUpdateRequest struct {
	Id     string `json:"id" binding:"required"`
	KindId string `json:"kindId" binding:"required"`
	CommonItemBaseData
}

type CommonItemBaseData struct {
	Name        string `json:"name" binding:"required"`
	Priority    int    `json:"priority" binding:"required"`
	MaxOrder    int    `json:"maxOrder" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Enabled     bool   `json:"enabled" binding:"required"`
}
