package item

import (
	"chico/takeout/handlers"
	usecase "chico/takeout/usecase/item"

	"github.com/gin-gonic/gin"
)

type OptionItemCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Priority    int    `json:"priority" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Enabled     *bool   `json:"enabled" binding:"required"`
}

type OptionItemCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

func (i *OptionItemCreateRequest) toModel() *usecase.OptionItemCreateModel {
	return &usecase.OptionItemCreateModel{Name: i.Name, Description: i.Description, Priority: i.Priority, Price: i.Price, Enabled: *i.Enabled}
}

type OptionItemUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Priority    int    `json:"priority" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Enabled     *bool   `json:"enabled" binding:"required"`
}

func (i *OptionItemUpdateRequest) toModel(id string) *usecase.OptionItemUpdateModel {
	return &usecase.OptionItemUpdateModel{Id: id, Name: i.Name, Description: i.Description, Priority: i.Priority, Price: i.Price, Enabled: *i.Enabled}
}

type OptionItemData struct {
	Id          string `json:"id"`
	Name        string `json:"name" binding:"required"`
	Priority    int    `json:"priority" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Enabled     bool   `json:"enabled" binding:"required"`
}

func newOptionItemData(item *usecase.OptionItemModel) *OptionItemData {
	return &OptionItemData{
		Id:          item.Id,
		Name:        item.Name,
		Description: item.Description,
		Priority:    item.Priority,
		Price:       item.Price,
		Enabled:     item.Enabled,
	}
}

type optionItemHandler struct {
	*handlers.BaseHandler
	usecase usecase.OptionItemUseCase
}

func NewOptionItemHandler(u usecase.OptionItemUseCase) *optionItemHandler {
	return &optionItemHandler{usecase: u}
}

func (i *optionItemHandler) GetAll(c *gin.Context) {
	items, err := i.usecase.FindAll()
	if err != nil {
		i.HandleError(c, err)
		return
	}

	models := []OptionItemData{}
	for _, item := range items {
		models = append(models, *newOptionItemData(&item))
	}
	i.HandleOK(c, models)
}

func (i *optionItemHandler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := i.usecase.Find((id))
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, newOptionItemData(item))
}

func (i *optionItemHandler) Post(c *gin.Context) {
	var req OptionItemCreateRequest
	if !i.ShouldBind(c, &req) {
		return
	}
	id, err := i.usecase.Create(req.toModel())
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, OptionItemCreateResponse{Id: id})
}

func (i *optionItemHandler) Put(c *gin.Context) {
	id := c.Param("id")
	var req OptionItemUpdateRequest
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

func (i *optionItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.usecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, nil)
}
