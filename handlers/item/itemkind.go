package item

import (
	"chico/takeout/handlers"
	usecase "chico/takeout/usecase/item"

	"github.com/gin-gonic/gin"
)

type ItemKindCreateRequest struct {
	Name          string   `json:"name" binding:"required"`
	Priority      int      `json:"priority" binding:"required"`
	OptionItemIds []string `json:"optionItemIds" binding:"required"`
}

type ItemKindCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

func (i *ItemKindCreateRequest) toModel() *usecase.ItemKindCreateModel {
	return &usecase.ItemKindCreateModel{Name: i.Name, Priority: i.Priority, OptionItemIds: i.OptionItemIds}
}

type ItemKindUpdateRequest struct {
	Name          string   `json:"name" binding:"required"`
	Priority      int      `json:"priority" binding:"required"`
	OptionItemIds []string `json:"optionItemIds" binding:"required"`
}

func (i *ItemKindUpdateRequest) toModel(id string) *usecase.ItemKindUpdateModel {
	return &usecase.ItemKindUpdateModel{Id: id, Name: i.Name, Priority: i.Priority, OptionItemIds: i.OptionItemIds}
}

type ItemKindData struct {
	Id            string   `json:"id"`
	Name          string   `json:"name" binding:"required"`
	Priority      int      `json:"priority" binding:"required"`
	OptionItemIds []string `json:"optionItemIds" binding:"required"`
}

func newItemKindData(item *usecase.ItemKindModel) *ItemKindData {
	return &ItemKindData{
		Id:            item.Id,
		Name:          item.Name,
		Priority:      item.Priority,
		OptionItemIds: item.OptionItemIds,
	}
}

type itemKindHandler struct {
	*handlers.BaseHandler
	usecase usecase.ItemKindUseCase
}

func NewItemKindHandler(u usecase.ItemKindUseCase) *itemKindHandler {
	return &itemKindHandler{usecase: u}
}

func (i *itemKindHandler) GetAll(c *gin.Context) {
	items, err := i.usecase.FindAll()
	if err != nil {
		i.HandleError(c, err)
		return
	}

	models := []ItemKindData{}
	for _, item := range items {
		models = append(models, *newItemKindData(&item))
	}
	i.HandleOK(c, models)
}

func (i *itemKindHandler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := i.usecase.Find((id))
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, newItemKindData(item))
}

func (i *itemKindHandler) Post(c *gin.Context) {
	var req ItemKindCreateRequest
	if !i.ShouldBind(c, &req) {
		return
	}
	id, err := i.usecase.Create(req.toModel())
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, ItemKindCreateResponse{Id: id})
}

func (i *itemKindHandler) Put(c *gin.Context) {
	id := c.Param("id")
	var req ItemKindUpdateRequest
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

func (i *itemKindHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.usecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, nil)
}
