package item

import (
	"chico/takeout/handlers"
	usecase "chico/takeout/usecase/item"

	"github.com/gin-gonic/gin"
)

type ItemKindCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Priority int    `form:"priority" binding:"required"`
}

type ItemKindCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

func (i *ItemKindCreateRequest) toModel() *usecase.ItemKindCreateModel {
	return &usecase.ItemKindCreateModel{Name: i.Name, Priority: i.Priority}
}

type ItemKindUpdateRequest struct {
	Id       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Priority int    `form:"priority" binding:"required"`
}

func (i *ItemKindUpdateRequest) toModel() *usecase.ItemKinddUpdateModel {
	return &usecase.ItemKinddUpdateModel{Id: i.Id, Name: i.Name, Priority: i.Priority}
}

type ItemKindData struct {
	Id       string `json:"id"`
	Name     string `json:"name" binding:"required"`
	Priority int    `form:"priority" binding:"required"`
}

func newItemKindData(item *usecase.ItemKindModel) *ItemKindData {
	return &ItemKindData{
		Id:       item.Id,
		Name:     item.Name,
		Priority: item.Priority,
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
	}
	i.HandleOK(c, newItemKindData(item))
}

func (i *itemKindHandler) Post(c *gin.Context) {
	var req ItemKindCreateRequest
	// validation is executed model
	c.ShouldBind(&req)
	id, err := i.usecase.Create(req.toModel())
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, ItemKindCreateResponse{Id: id})
}

func (i *itemKindHandler) Put(c *gin.Context) {
	var req ItemKindUpdateRequest
	// validation is executed model
	c.ShouldBind(&req)
	err := i.usecase.Update(req.toModel())
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, nil)
}

func (i *itemKindHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.usecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, nil)
}
