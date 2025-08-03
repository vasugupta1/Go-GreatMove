package main

import (
	"log"

	"github.com/gin-gonic/gin"
	getproperties "github.com/greatmove/backend/internal/features/getproperties"
	health "github.com/greatmove/backend/internal/features/health"
	"github.com/greatmove/backend/models"
	services "github.com/greatmove/backend/services"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	router := gin.Default()
	router.GET("/healthCheck", health.HealthHandler)
	router.GET("/properties", getPropertiesHandler().GetPropertiesHandler)
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getPropertiesHandler() getproperties.GetPropertiesHandler {
	httpClient := services.ConstructHttpClient()
	database := getMonogoDatabase()
	repository := services.ConstructMongoRepository[models.Property](database, "properties")
	rightMoveService := services.ConstructRightMove(httpClient)
	return getproperties.ConstructGetPropertiesHandler(rightMoveService, *repository)
}

func getMonogoDatabase() *mongo.Database {
	uri := "mongodb://root:example@localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	return client.Database("greatmove")
}
