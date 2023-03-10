package message

import (
	"chico/takeout/common"
	"chico/takeout/handlers"
	usecase "chico/takeout/usecase/message"

	"github.com/gin-gonic/gin"
)

type StoreMessageData struct {
	Id      string `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Edited  string `json:"edited" binding:"required"`
	Created string `json:"created" binding:"required"`
}

func newStoreMessageData(item *usecase.StoreMessageModel) *StoreMessageData {
	return &StoreMessageData{
		Id:      item.Id,
		Content: item.Content,
		Edited:  common.ConvertTimeToDateTimeStr(item.Edited),
		Created: common.ConvertTimeToDateTimeStr(item.Created),
	}
}

type storeMessageHandler struct {
	*handlers.BaseHandler
	usecase usecase.StoreMessageUseCase
}

func NewStoreMessageHandler(u usecase.StoreMessageUseCase) *storeMessageHandler {
	return &storeMessageHandler{usecase: u}
}

type StoreMessageCreateRequest struct {
	Id      string `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type StoreMessageCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

func (s *StoreMessageCreateRequest) toModel() *usecase.StoreMessageCreateModel {
	return &usecase.StoreMessageCreateModel{Id: s.Id, Content: s.Content}
}

type StoreMessageUpdateRequest struct {
	Content string `json:"content" binding:"required"`
}

func (s *StoreMessageUpdateRequest) toModel(id string) *usecase.StoreMessageUpdateModel {
	return &usecase.StoreMessageUpdateModel{Id: id, Content: s.Content}
}

func (s *storeMessageHandler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := s.usecase.Find((id))
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, newStoreMessageData(item))
}

func (s *storeMessageHandler) Post(c *gin.Context) {
	var req StoreMessageCreateRequest
	if !s.ShouldBind(c, &req) {
		return
	}
	id, err := s.usecase.Create(req.toModel())
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, StoreMessageCreateResponse{Id: id})
}

func (s *storeMessageHandler) Put(c *gin.Context) {
	id := c.Param("id")
	var req StoreMessageUpdateRequest
	if !s.ShouldBind(c, &req) {
		return
	}
	err := s.usecase.Update(req.toModel(id))
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, nil)
}
