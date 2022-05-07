package store

import (
	"chico/takeout/handlers"
	usecases "chico/takeout/usecase/store"

	"github.com/gin-gonic/gin"
)

type SpecialHolidayData struct {
	Id    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

func newSpecialHolidayData(model usecases.SpecialHolidayModel) *SpecialHolidayData {
	return &SpecialHolidayData{
		Id:    model.Id,
		Name:  model.Name,
		Start: model.Start,
		End:   model.End,
	}
}

type SpecialHolidayCreateData struct {
	Name  string `json:"name" binding:"required"`
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

func (b *SpecialHolidayCreateData) toModel() *usecases.SpecialHolidayCreateModel {
	return &usecases.SpecialHolidayCreateModel{
		Name:  b.Name,
		Start: b.Start,
		End:   b.End,
	}
}

type SpecialHolidayCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

type SpecialHolidayUpdateData struct {
	Id    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

func (b *SpecialHolidayUpdateData) toModel() *usecases.SpecialHolidayUpdateModel {
	return &usecases.SpecialHolidayUpdateModel{
		Id:    b.Id,
		Name:  b.Name,
		Start: b.Start,
		End:   b.End,
	}
}

type specialHolidayHandler struct {
	*handlers.BaseHandler
	usecase usecases.SpecialHolidayUseCase
}

func NewSpecialHolidayHandler(usecase usecases.SpecialHolidayUseCase) *specialHolidayHandler {
	return &specialHolidayHandler{
		usecase: usecase,
	}
}

func (s *specialHolidayHandler) Get(c *gin.Context) {
	id := c.Param("id")
	model, err := s.usecase.Find(id)
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, newSpecialHolidayData(*model))
}

func (s *specialHolidayHandler) GetAll(c *gin.Context) {
	alls, err := s.usecase.FindAll()
	if err != nil {
		s.HandleError(c, err)
		return
	}
	allData := []SpecialHolidayData{}
	for _, item := range alls {
		allData = append(allData, *newSpecialHolidayData(item))
	}
	s.HandleOK(c, allData)
}

func (s *specialHolidayHandler) Post(c *gin.Context) {
	var req SpecialHolidayCreateData
	// validation is executed model
	c.ShouldBind(&req)
	id, err := s.usecase.Create(req.toModel())
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, SpecialHolidayCreateResponse{Id: id})
}

func (s *specialHolidayHandler) Put(c *gin.Context) {
	var req SpecialHolidayUpdateData
	// validation is executed model
	c.ShouldBind(&req)
	err := s.usecase.Update(req.toModel())
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, nil)
}

func (i *specialHolidayHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.usecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
		return
	}
	i.HandleOK(c, nil)
}
