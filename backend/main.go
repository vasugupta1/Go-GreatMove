package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	getproperties "github.com/greatmove/backend/internal/features/getproperties"
	health "github.com/greatmove/backend/internal/features/health"
	"github.com/greatmove/backend/internal/features/transformproperties"
	models "github.com/greatmove/backend/models"
	services "github.com/greatmove/backend/services"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	log.Println("Starting GreatMove backend...")
	config := loadConfiguration()
	httpClient := services.ConstructHttpClient()
	rightMoveService := services.ConstructRightMove(httpClient, config.RightMoveAPIURL)
	repository := getPropertyRepository(config)
	router := gin.Default()
	router.GET("/healthCheck", health.HealthHandler)
	router.GET("/properties", getPropertiesHandler(repository, rightMoveService).GetPropertiesHandler)
	router.POST("/transform", transformPropertiesHandler(repository, rightMoveService).TransformPropertiesHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  300 * time.Second, // 5 minutes
		WriteTimeout: 300 * time.Second, // 5 minutes
		IdleTimeout:  300 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getPropertiesHandler(repository services.MongoRepository[models.Property], rightMoveService services.RightMoveService) getproperties.GetPropertiesHandler {

	return getproperties.ConstructGetPropertiesHandler(rightMoveService, repository)
}

func transformPropertiesHandler(repository services.MongoRepository[models.Property], rightMoveService services.RightMoveService) transformproperties.TransformPropertiesHandler {
	return transformproperties.ConstructGetPropertiesHandler(repository, rightMoveService)
}

func getPropertyRepository(configuration models.Configuration) services.MongoRepository[models.Property] {
	clientOptions := options.Client().ApplyURI(configuration.MongoURI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	database := client.Database(configuration.MongoDB)
	return services.ConstructMongoRepository[models.Property](database, configuration.MongoCollection)
}

func loadConfiguration() models.Configuration {
	file, err := os.Open("configuration.json")
	if err != nil {
		log.Fatal("Failed to open configuration file:", err)
	}
	defer file.Close()
	var config models.Configuration
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("Failed to decode configuration file:", err)
	}
	return config
}
