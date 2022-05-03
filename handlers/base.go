package handlers

import (
	"errors"
	"net/http"

	"chico/takeout/common"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (b *BaseHandler) HandleError(c *gin.Context, e error) {
	var vErr *common.ValidationError
	if errors.As(e, &vErr) {
		c.String(http.StatusBadRequest, vErr.Error())
		return
	}
	var uErr *common.UpdateTargetNotFoundError
	if errors.As(e, &uErr) {
		c.String(http.StatusBadRequest, uErr.Error())
		return
	}
	var nErr *common.NotFoundError
	if errors.As(e, &nErr) {
		c.String(http.StatusNotFound, nErr.Error())
		return
	}
	b.HandleServerError(c)
}

func (b *BaseHandler) HandleServerError(c *gin.Context) {
	c.String(http.StatusInternalServerError, "Server Error")
}

func (b *BaseHandler) HandleOK(c *gin.Context, jsonData interface{}) {
	if jsonData != nil {
		c.JSON(http.StatusOK, jsonData)
		return
	}
	c.JSON(http.StatusOK, "")
}
