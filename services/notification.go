package services

import (
	"context"

	"github.com/akrck02/valhalla-core-sdk/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Alert struct {
	Title   string
	Message string
}

func AlertTeam(conn context.Context, client *mongo.Client, team models.Team) models.Error {

	return models.Error{
		Status:  200,
		Message: "Ok.",
	}
}
