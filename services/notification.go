package services

import (
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/models"
)

type Alert struct {
	Title   string
	Message string
}

func AlertTeam(team models.Team) models.Error {

	return models.Error{
		Status:  http.HTTP_STATUS_OK,
		Message: "Ok.",
	}
}
