package item

import (
	"chico/takeout/handlers"
	usecase "chico/takeout/usecase/item"

	"github.com/gin-gonic/gin"
)

type StockItemResponse struct {
	CommonItemResponse
	Remain int `json:"remain" binding:"required"`
}

func newStockItemData(item *usecase.StockItemModel) *StockItemResponse {
	kind := newItemKindData(&item.Kind)
	enabled := item.Enabled
	imageUrl := item.ImageUrl
	return &StockItemResponse{
		CommonItemResponse: CommonItemResponse{
			Id:   item.Id,
			Kind: *kind,
			CommonItemBaseData: CommonItemBaseData{
				Name:        item.Name,
				Priority:    item.Priority,
				MaxOrder:    item.MaxOrder,
				Price:       item.Price,
				Description: item.Description,
				Enabled:     &enabled,
				ImageUrl:    &imageUrl,
			},
		},
		Remain: item.Remain,
	}
}

type StockItemCreateRequest struct {
	CommonItemCreateRequest
}

type StockItemCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

func (s *StockItemCreateRequest) toModel() *usecase.StockItemCreateModel {
	return &usecase.StockItemCreateModel{
		CommonItemCreateModel: usecase.CommonItemCreateModel{
			KindId: s.KindId,
			CommonItemBaseModel: usecase.CommonItemBaseModel{
				Name: s.Name, Priority: s.Priority, MaxOrder: s.MaxOrder, Price: *s.Price, Description: s.Description, Enabled: *s.Enabled, ImageUrl: *s.ImageUrl,
			},
		},
	}
}

type StockItemUpdateRequest struct {
	CommonItemUpdateRequest
}

func (s *StockItemUpdateRequest) toModel(id string) *usecase.StockItemUpdateModel {
	return &usecase.StockItemUpdateModel{
		CommonItemUpdateModel: usecase.CommonItemUpdateModel{
			Id:     id,
			KindId: s.KindId,
			CommonItemBaseModel: usecase.CommonItemBaseModel{
				Name: s.Name, Priority: s.Priority, MaxOrder: s.MaxOrder, Price: *s.Price, Description: s.Description, Enabled: *s.Enabled, ImageUrl: *s.ImageUrl,
			},
		},
	}
}

type StockItemRemainUpdateRequest struct {
	Remain int    `json:"remain" binding:"required"`
}

func (s *StockItemRemainUpdateRequest) toModel(id string) *usecase.StockItemRemainUpdateModel {
	return &usecase.StockItemRemainUpdateModel{Id: id, Remain: s.Remain}
}

type stockItemHandler struct {
	*handlers.BaseHandler
	usecase usecase.StockItemUseCase
}

func NewStockItemHandler(u usecase.StockItemUseCase) *stockItemHandler {
	return &stockItemHandler{usecase: u}
}

func (i *stockItemHandler) GetAll(c *gin.Context) {
	items, err := i.usecase.FindAll()
	if err != nil {
		i.HandleError(c, err)
		return
	}

	models := []StockItemResponse{}
	for _, item := range items {
		models = append(models, *newStockItemData(&item))
	}
	i.HandleOK(c, models)
}

func (i *stockItemHandler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := i.usecase.Find((id))
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, newStockItemData(item))
}

func (i *stockItemHandler) Post(c *gin.Context) {
	var req StockItemCreateRequest
	if !i.ShouldBind(c, &req) {
		return
	}
	id, err := i.usecase.Create(req.toModel())
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, StockItemCreateResponse{Id: id})
}

func (i *stockItemHandler) Put(c *gin.Context) {
	id := c.Param("id")
	var req StockItemUpdateRequest
	if !i.ShouldBind(c, &req) {
		return
	}
	err := i.usecase.Update(req.toModel(id))
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, nil)
}

func (i *stockItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.usecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, nil)
}

func (i *stockItemHandler) PutRemain(c *gin.Context) {
	id := c.Param("id")
	var req StockItemRemainUpdateRequest
	if !i.ShouldBind(c, &req) {
		return
	}
	err := i.usecase.UpdateRemain(req.toModel(id))
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, nil)
}
