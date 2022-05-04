package store

import (
	"chico/takeout/handlers"
	usecases "chico/takeout/usecase/store"

	"github.com/gin-gonic/gin"
)

type BusinessHoursData struct {
	Schedules []BusinessHourData `json:"schedules"`
}

func newBusinessHoursData(model usecases.BusinessHoursModel) *BusinessHoursData {
	schedules := []BusinessHourData{}
	for _, schedule := range model.Schedules {
		schedules = append(schedules, *newBusinessHourData(schedule))
	}
	return &BusinessHoursData{
		Schedules: schedules,
	}
}

type BusinessHoursUpdateData struct {
	Id       string             `json:"id" binding:"required"`
	Name     string             `json:"name" binding:"required"`
	Start    string             `json:"start" binding:"required"`
	End      string             `json:"end" binding:"required"`
	Weekdays []usecases.Weekday `json:"weekdays" binding:"required"`
}

func (b *BusinessHoursUpdateData) toModel() *usecases.BusinessHoursUpdateModel {
	return &usecases.BusinessHoursUpdateModel{
		Id:       b.Id,
		Name:     b.Name,
		Start:    b.Start,
		End:      b.End,
		Weekdays: b.Weekdays,
	}
}

type BusinessHourData struct {
	Id       string             `json:"id" binding:"required"`
	Name     string             `json:"name" binding:"required"`
	Start    string             `json:"start" binding:"required"`
	End      string             `json:"end" binding:"required"`
	Weekdays []usecases.Weekday `json:"weekdays" binding:"required"`
}

func newBusinessHourData(model usecases.BusinessHourModel) *BusinessHourData {
	return &BusinessHourData{
		Id:       model.Id,
		Name:     model.Name,
		Start:    model.Start,
		End:      model.End,
		Weekdays: model.Weekdays,
	}
}

type businessHoursHandler struct {
	*handlers.BaseHandler
	usecase usecases.BusinessHoursUseCase
}

func NewbusinessHoursHandler(usecase usecases.BusinessHoursUseCase) *businessHoursHandler {
	return &businessHoursHandler{
		usecase: usecase,
	}
}

func (b *businessHoursHandler) Get(c *gin.Context) {
	model, err := b.usecase.Fetch()
	if err != nil {
		b.HandleError(c, err)
	}
	b.HandleOK(c, newBusinessHoursData(*model))
}

func (b *businessHoursHandler) Put(c *gin.Context) {
	var req BusinessHoursUpdateData
	// validation is executed model
	c.ShouldBind(&req)
	err := b.usecase.Update(*req.toModel())
	if err != nil {
		b.HandleError(c, err)
	}
	b.HandleOK(c, nil)
}
