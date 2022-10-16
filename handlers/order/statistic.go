package order

import (
	"chico/takeout/handlers"
	queryUseCases "chico/takeout/usecase/order/query"

	"github.com/gin-gonic/gin"
)

type MonthlyStatisticResponse struct {
	Data []MonthlyData `json:"data" binding:"required"`
}

type MonthlyData struct {
	Month         string `json:"month" binding:"required"`
	OrderTotal    int    `json:"order_total" binding:"required"`
	QuantityTotal int    `json:"quantity_total" binding:"required"`
	MoneyTotal    int    `json:"money_total" binding:"required"`
}

func newMonthlyStatisticResponse(m queryUseCases.MonthlyStatisticData) *MonthlyStatisticResponse {
	data := []MonthlyData{}

	for _, d := range m.Data {
		data = append(data, *newMonthlyData(d))
	}

	return &MonthlyStatisticResponse{
		Data: data,
	}
}

func newMonthlyData(d queryUseCases.MonthlyData) *MonthlyData {
	return &MonthlyData{
		Month:         d.Month,
		OrderTotal:    d.OrderTotal,
		QuantityTotal: d.QuantityTotal,
		MoneyTotal:    d.MoneyTotal,
	}
}

type statisticInfoHandler struct {
	*handlers.BaseHandler
	queryUseCase *queryUseCases.OrderStatisticUseCase
}

func NewStatisticInfoHandler(queryUseCase *queryUseCases.OrderStatisticUseCase) *statisticInfoHandler {
	return &statisticInfoHandler{
		queryUseCase: queryUseCase,
	}
}

func (s *statisticInfoHandler) GetMonthly(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	req := queryUseCases.MonthlyStatisticRequestModel{
		Start: start,
		End:   end,
	}
	model, err := s.queryUseCase.FetchMonthlyData(req)
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, newMonthlyStatisticResponse(*model))
}
