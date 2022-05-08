package item

import (
	"chico/takeout/handlers"
	usecase "chico/takeout/usecase/item"

	"github.com/gin-gonic/gin"
)

type FoodItemResponse struct {
	CommonItemResponse
	ScheduleIds    []string `json:"scheduleIds" binding:"required"`
	MaxOrderPerDay int      `json:"maxOrderPerDay" binding:"required"`
}

func newFoodItemData(item *usecase.FoodItemModel) *FoodItemResponse {
	kind := newItemKindData(&item.Kind)
	return &FoodItemResponse{
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
		ScheduleIds:    item.ScheduleIds,
		MaxOrderPerDay: item.MaxOrderPerDay,
	}
}

type FoodItemCreateRequest struct {
	CommonItemCreateRequest
	ScheduleIds    []string `json:"scheduleIds" binding:"required"`
	MaxOrderPerDay int      `json:"maxOrderPerDay" binding:"required"`
}

type FoodItemCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

func (s *FoodItemCreateRequest) toModel() *usecase.FoodItemCreateModel {
	return &usecase.FoodItemCreateModel{
		CommonItemCreateModel: usecase.CommonItemCreateModel{
			KindId: s.KindId,
			CommonItemBaseModel: usecase.CommonItemBaseModel{
				Name: s.Name, Priority: s.Priority, MaxOrder: s.MaxOrder, Price: s.Price, Description: s.Description,
			},
		},
		ScheduleIds: s.ScheduleIds, MaxOrderPerDay: s.MaxOrderPerDay,
	}
}

type FoodItemUpdateRequest struct {
	CommonItemUpdateRequest
	ScheduleIds    []string `json:"scheduleIds" binding:"required"`
	MaxOrderPerDay int      `json:"maxOrderPerDay" binding:"required"`
}

func (s *FoodItemUpdateRequest) toModel(id string) *usecase.FoodItemUpdateModel {
	return &usecase.FoodItemUpdateModel{
		CommonItemUpdateModel: usecase.CommonItemUpdateModel{
			Id:     id,
			KindId: s.KindId,
			CommonItemBaseModel: usecase.CommonItemBaseModel{
				Name: s.Name, Priority: s.Priority, MaxOrder: s.MaxOrder, Price: s.Price, Description: s.Description,
			},
		},
		ScheduleIds: s.ScheduleIds, MaxOrderPerDay: s.MaxOrderPerDay,
	}
}

type foodItemHandler struct {
	*handlers.BaseHandler
	usecase usecase.FoodItemUseCase
}

func NewFoodItemHandler(u usecase.FoodItemUseCase) *foodItemHandler {
	return &foodItemHandler{usecase: u}
}

func (i *foodItemHandler) GetAll(c *gin.Context) {
	items, err := i.usecase.FindAll()
	if err != nil {
		i.HandleError(c, err)
		return
	}

	models := []FoodItemResponse{}
	for _, item := range items {
		models = append(models, *newFoodItemData(&item))
	}
	i.HandleOK(c, models)
}

func (i *foodItemHandler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := i.usecase.Find((id))
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, newFoodItemData(item))
}

func (i *foodItemHandler) Post(c *gin.Context) {
	var req FoodItemCreateRequest
	if !i.ShouldBind(c, &req) {
		return
	}
	id, err := i.usecase.Create(req.toModel())
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, FoodItemCreateResponse{Id: id})
}

func (i *foodItemHandler) Put(c *gin.Context) {
	id := c.Param("id")
	var req FoodItemUpdateRequest
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

func (i *foodItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.usecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, nil)
}
