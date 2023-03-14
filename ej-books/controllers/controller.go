package controllers

import (
	"books-api/dtos"
	service "books-api/services"
	e "books-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	service service.Service
}

func NewController(service service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (ctrl *Controller) Get(c *gin.Context) {
	book, apiErr := ctrl.service.Get(c.Param("id"))
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusOK, book)
}

func (ctrl *Controller) Insert(c *gin.Context) {
	var book dtos.BookDTO
	if err := c.BindJSON(&book); err != nil {
		apiErr := e.NewBadRequestApiError(err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	book, apiErr := ctrl.service.Insert(book)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, book)
}
