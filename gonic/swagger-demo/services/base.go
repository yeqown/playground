package services

import (
	"swagger-demo/models"
)

// ListGetModels ...
func ListGetModels() []*models.GetModel {
	return []*models.GetModel{
		{ModelName: "GET", Get: true},
		{ModelName: "GET", Get: true},
		{ModelName: "GET", Get: true},
	}
}

// ListPostModels ...
func ListPostModels() []*models.PostModel {
	return []*models.PostModel{
		{ModelName: "GET", Post: true},
		{ModelName: "GET", Post: true},
		{ModelName: "GET", Post: true},
	}
}
