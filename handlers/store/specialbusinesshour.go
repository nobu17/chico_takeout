package store

import (
	"chico/takeout/handlers"
	usecases "chico/takeout/usecase/store"

	"github.com/gin-gonic/gin"
)

type SpecialBusinessHourData struct {
	Id             string `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Start          string `json:"start" binding:"required"`
	End            string `json:"end" binding:"required"`
	BusinessHourId string `json:"businessHourId" binding:"required"`
}

func newSpecialBusinessHourData(model usecases.SpecialBusinessHourModel) *SpecialBusinessHourData {
	return &SpecialBusinessHourData{
		Id:             model.Id,
		Name:           model.Name,
		Date:           model.Date,
		Start:          model.Start,
		End:            model.End,
		BusinessHourId: model.BusinessHourId,
	}
}

type SpecialBusinessHourCreateRequest struct {
	Name           string `json:"name" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Start          string `json:"start" binding:"required"`
	End            string `json:"end" binding:"required"`
	BusinessHourId string `json:"businessHourId" binding:"required"`
}

func (b *SpecialBusinessHourCreateRequest) toModel() *usecases.SpecialBusinessHourCreateModel {
	return &usecases.SpecialBusinessHourCreateModel{
		Name:           b.Name,
		Date:           b.Date,
		Start:          b.Start,
		End:            b.End,
		BusinessHourId: b.BusinessHourId,
	}
}

type SpecialBusinessHourCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

type SpecialBusinessHourUpdateRequest struct {
	Id             string `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Start          string `json:"start" binding:"required"`
	End            string `json:"end" binding:"required"`
	BusinessHourId string `json:"businessHourId" binding:"required"`
}

func (b *SpecialBusinessHourUpdateRequest) toModel() *usecases.SpecialBusinessHourUpdateModel {
	return &usecases.SpecialBusinessHourUpdateModel{
		Id:             b.Id,
		Name:           b.Name,
		Date:           b.Date,
		Start:          b.Start,
		End:            b.End,
		BusinessHourId: b.BusinessHourId,
	}
}

type specialBusinessHourHandler struct {
	*handlers.BaseHandler
	usecase usecases.SpecialBusinessHoursUseCase
}

func NewSpecialBusinessHourHandler(usecase usecases.SpecialBusinessHoursUseCase) *specialBusinessHourHandler {
	return &specialBusinessHourHandler{
		usecase: usecase,
	}
}

func (s *specialBusinessHourHandler) Get(c *gin.Context) {
	id := c.Param("id")
	model, err := s.usecase.Find(id)
	if err != nil {
		s.HandleError(c, err)
	}
	s.HandleOK(c, newSpecialBusinessHourData(*model))
}

func (s *specialBusinessHourHandler) GetAll(c *gin.Context) {
	alls, err := s.usecase.FindAll()
	if err != nil {
		s.HandleError(c, err)
	}
	allData := []SpecialBusinessHourData{}
	for _, item := range alls {
		allData = append(allData, *newSpecialBusinessHourData(item))
	}
	s.HandleOK(c, allData)
}

func (s *specialBusinessHourHandler) Post(c *gin.Context) {
	var req SpecialBusinessHourCreateRequest
	// validation is executed model
	c.ShouldBind(&req)
	id, err := s.usecase.Create(*req.toModel())
	if err != nil {
		s.HandleError(c, err)
	}
	s.HandleOK(c, SpecialHolidayCreateResponse{Id: id})
}

func (s *specialBusinessHourHandler) Put(c *gin.Context) {
	var req SpecialBusinessHourUpdateRequest
	// validation is executed model
	c.ShouldBind(&req)
	err := s.usecase.Update(*req.toModel())
	if err != nil {
		s.HandleError(c, err)
	}
	s.HandleOK(c, nil)
}

func (i *specialBusinessHourHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := i.usecase.Delete(id)
	if err != nil {
		i.HandleError(c, err)
	}
	i.HandleOK(c, nil)
}
