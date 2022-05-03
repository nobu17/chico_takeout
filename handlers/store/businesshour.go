package store

import (
	"chico/takeout/handlers"
	usecases "chico/takeout/usecase/store"

	"github.com/gin-gonic/gin"
)

type BusinessHoursData struct {
	Morning BusinessHourData `json:"morning"`
	Lunch   BusinessHourData `json:"lunch"`
	Dinner  BusinessHourData `json:"dinner"`
}

func newBusinessHoursData(model usecases.BusinessHoursModel) *BusinessHoursData {
	return &BusinessHoursData{
		Morning: *newBusinessHourData(model.Morning),
		Lunch:   *newBusinessHourData(model.Lunch),
		Dinner:  *newBusinessHourData(model.Dinner),
	}
}

type BusinessHoursUpdateData struct {
	Morning *BusinessHourData `json:"morning"`
	Lunch   *BusinessHourData `json:"lunch"`
	Dinner  *BusinessHourData `json:"dinner"`
}

func (b *BusinessHoursUpdateData) toModel() *usecases.BusinessHoursUpdateModel {
	var morning *usecases.BusinessHourModel
	var lunch *usecases.BusinessHourModel
	var dinner *usecases.BusinessHourModel
	if b.Morning != nil {
		morning = b.Morning.toModel()
	}
	if b.Lunch != nil {
		lunch = b.Lunch.toModel()
	}
	if b.Dinner != nil {
		dinner = b.Dinner.toModel()
	}
	return &usecases.BusinessHoursUpdateModel{
		Morning: morning,
		Lunch:   lunch,
		Dinner:  dinner,
	}
}

type BusinessHourData struct {
	Start    string             `json:"start" binding:"required"`
	End      string             `json:"end" binding:"required"`
	Weekdays []usecases.Weekday `json:"weekdays" binding:"required"`
}

func (b *BusinessHourData) toModel() *usecases.BusinessHourModel {
	return &usecases.BusinessHourModel{
		Start:    b.Start,
		End:      b.End,
		Weekdays: b.Weekdays,
	}
}

func newBusinessHourData(model usecases.BusinessHourModel) *BusinessHourData {
	return &BusinessHourData{
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
