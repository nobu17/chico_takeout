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
				Name: s.Name, Priority: s.Priority, MaxOrder: s.MaxOrder, Price: s.Price, Description: s.Description,
			},
		},
	}
}

type StockItemUpdateRequest struct {
	CommonItemUpdateRequest
}

func (s *StockItemUpdateRequest) toModel() *usecase.StockItemUpdateModel {
	return &usecase.StockItemUpdateModel{
		CommonItemUpdateModel: usecase.CommonItemUpdateModel{
			Id:     s.Id,
			KindId: s.KindId,
			CommonItemBaseModel: usecase.CommonItemBaseModel{
				Name: s.Name, Priority: s.Priority, MaxOrder: s.MaxOrder, Price: s.Price, Description: s.Description,
			},
		},
	}
}

type StockItemRemainUpdateRequest struct {
	Id     string `json:"id" binding:"required"`
	Remain int    `json:"remain" binding:"required"`
}

func (s *StockItemRemainUpdateRequest) toModel() *usecase.StockItemRemainUpdateModel {
	return &usecase.StockItemRemainUpdateModel{Id: s.Id, Remain: s.Remain}
}

type stockItemHandler struct {
	*handlers.BaseHandler
	useecase usecase.StockItemUseCase
}

func NewStockItemHandler(u usecase.StockItemUseCase) *stockItemHandler {
	return &stockItemHandler{useecase: u}
}

func (i *stockItemHandler) GetAll(c *gin.Context) {
	items, err := i.useecase.FindAll()
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
	item, err := i.useecase.Find((id))
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, newStockItemData(item))
}

func (i *stockItemHandler) Post(c *gin.Context) {
	var req StockItemCreateRequest
	// validation is executed model
	c.ShouldBind(&req)
	id, err := i.useecase.Create(*req.toModel())
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, StockItemCreateResponse{Id: id})
}

func (i *stockItemHandler) Put(c *gin.Context) {
	var req StockItemUpdateRequest
	// validation is executed model
	c.ShouldBind(&req)
	err := i.useecase.Update(*req.toModel())
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, nil)
}

func (i *stockItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.useecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, nil)
}

func (i *stockItemHandler) PutRemain(c *gin.Context) {
	var req StockItemRemainUpdateRequest
	// validation is executed model
	c.ShouldBind(&req)
	err := i.useecase.UpdateRemain(*req.toModel())
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, nil)
}
