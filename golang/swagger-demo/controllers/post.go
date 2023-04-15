package controllers

import (
	"log"
	"net/http"

	"swagger-demo/models"
	"swagger-demo/services"

	"github.com/gin-gonic/gin"
)

// swagger:parameters opid-post
type postForm struct {
	// in: body
	ID int `form:"id" binding:"required"`
}

// swagger:model
type postmodels []*models.PostModel

// postResponse response demo of post controller
// swagger:response postResponse
type postResponse struct {
	// code int
	// in: body
	Code int `json:"code"`
	// models postmodels
	// in: body
	Models postmodels `json:"models"`
}

// Post ...
func Post(c *gin.Context) {
	// swagger:route POST /base/post 范例 opid-post
	//
	// swagger Get范例
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
	//       default: postResponse
	var (
		form = new(postForm)
		resp = new(postResponse)
	)

	resp.Code = 0
	resp.Models = services.ListPostModels()

	log.Println("form is :", *form)
	c.JSON(http.StatusOK, resp)
}
