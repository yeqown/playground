package controllers

import (
	"log"
	"net/http"

	"swagger-demo/models"
	"swagger-demo/services"

	"github.com/gin-gonic/gin"
)

// swagger:parameters opid-get
type getForm struct {
	// in: query
	ID int `form:"id" binding:"required"`
}

// swagger:model
type getmodels []*models.GetModel

// getResponse response demo of get controller
// swagger:response getResponse
type getResponse struct {
	// code int
	// in: body
	Code int `json:"code"`
	// modesl getmodel
	// in: body
	Models getmodels `json:"models"`
}

// Get ...
func Get(c *gin.Context) {
	// swagger:route GET /base/get 范例 opid-get
	//
	// swagger Get范例
	//
	//
	//     Consumes:
	//     - application/json
	//     - application/x-protobuf
	//
	//     Produces:
	//     - application/json
	//     - application/x-protobuf
	//
	//     Schemes: http, https, ws, wss
	//
	//     Responses:
	//       default: getResponse
	var (
		form = new(getForm)
		resp = new(getResponse)
	)

	resp.Code = 0
	resp.Models = services.ListGetModels()

	log.Println("form is :", *form)
	c.JSON(http.StatusOK, resp)
}
