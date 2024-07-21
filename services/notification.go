package services

import (
	"net/http"

	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
	teammodels "github.com/akrck02/valhalla-core-sdk/models/team"
)

type Alert struct {
	Title   string
	Message string
}

func AlertTeam(team teammodels.Team) *apimodels.Error {

	return &apimodels.Error{
		Status:  http.StatusOK,
		Message: "Ok.",
	}
}
