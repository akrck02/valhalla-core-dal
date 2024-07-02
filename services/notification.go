package services

import (
	"github.com/akrck02/valhalla-core-sdk/http"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	teammodels "github.com/akrck02/valhalla-core-sdk/models/team"
)

type Alert struct {
	Title   string
	Message string
}

func AlertTeam(team teammodels.Team) systemmodels.Error {

	return systemmodels.Error{
		Status:  http.HTTP_STATUS_OK,
		Message: "Ok.",
	}
}
