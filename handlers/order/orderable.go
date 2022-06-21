package order

import (
	"chico/takeout/handlers"
	queryUseCases "chico/takeout/usecase/order/query"

	"github.com/gin-gonic/gin"
)

type orderableInfoHandler struct {
	*handlers.BaseHandler
	queryUseCase *queryUseCases.OrderQueryUseCase
}

func NewOrderableInfoHandler(queryUseCase *queryUseCases.OrderQueryUseCase) *orderableInfoHandler {
	return &orderableInfoHandler{
		queryUseCase: queryUseCase,
	}
}

type OrderableInfoRequestResponse struct {
	StartDate  string                        `json:"startDate" binding:"required"`
	EndDate    string                        `json:"endDate" binding:"required"`
	PerDayInfo []PerDayOrderableInfoResponse `json:"perDayInfo" binding:"required"`
}

type PerDayOrderableInfoResponse struct {
	Date       string                      `json:"date" binding:"required"`
	HourTypeId string                      `json:"hourTypeId" binding:"required"`
	StartTime  string                      `json:"startTime" binding:"required"`
	EndTime    string                      `json:"endTime" binding:"required"`
	Items      []OrderableItemInfoResponse `json:"items" binding:"required"`
}

type OrderableItemInfoResponse struct {
	Id       string `json:"id" binding:"required"`
	ItemType string `json:"itemType" binding:"required"`
	Remain   int    `json:"remain" binding:"required"`
}

func newOrderableInfoRequestResponse(o queryUseCases.OrderableInfo) *OrderableInfoRequestResponse {
	infoList := []PerDayOrderableInfoResponse{}
	for _, info := range o.PerDayInfo {
		resp := newPerDayOrderableInfoResponse(info)
		infoList = append(infoList, *resp)
	}
	return &OrderableInfoRequestResponse{
		StartDate:  o.StartDate,
		EndDate:    o.EndDate,
		PerDayInfo: infoList,
	}

}

func newPerDayOrderableInfoResponse(p queryUseCases.PerDayOrderableInfo) *PerDayOrderableInfoResponse {
	items := []OrderableItemInfoResponse{}
	for _, item := range p.Items {
		items = append(items, *newOrderableItemInfoResponse(item))
	}
	return &PerDayOrderableInfoResponse{
		Date:       p.Date,
		HourTypeId: p.HourTypeId,
		StartTime:  p.StartTime,
		EndTime:    p.EndTime,
		Items:      items,
	}
}

func newOrderableItemInfoResponse(o queryUseCases.OrderableItemInfo) *OrderableItemInfoResponse {
	return &OrderableItemInfoResponse{
		Id:       o.Id,
		ItemType: o.ItemType,
		Remain:   o.Remain,
	}
}

func (s *orderableInfoHandler) Get(c *gin.Context) {
	model, err := s.queryUseCase.FetchOrderableInfo()
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, newOrderableInfoRequestResponse(*model))
}
