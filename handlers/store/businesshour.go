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
	Name     string             `json:"name" binding:"required"`
	Start    string             `json:"start" binding:"required"`
	End      string             `json:"end" binding:"required"`
	Weekdays []usecases.Weekday `json:"weekdays" binding:"required"`
}

func (b *BusinessHoursUpdateData) toModel(id string) *usecases.BusinessHoursUpdateModel {
	return &usecases.BusinessHoursUpdateModel{
		Id:       id,
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
	model, err := b.usecase.GetAll()
	if err != nil {
		b.HandleError(c, err)
		return
	}
	b.HandleOK(c, newBusinessHoursData(*model))
}

func (b *businessHoursHandler) Put(c *gin.Context) {
	id := c.Param("id")
	var req BusinessHoursUpdateData
	if !b.ShouldBind(c, &req) {
		return
	}
	err := b.usecase.Update(req.toModel(id))
	if err != nil {
		b.HandleError(c, err)
		return
	}
	b.HandleOK(c, nil)
}
