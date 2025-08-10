package transformproperties

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/greatmove/backend/models"
	"github.com/greatmove/backend/services"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransformPropertiesHandler interface {
	TransformPropertiesHandler(c *gin.Context)
}

type TransformProperties struct {
	RepositoryService services.MongoRepository[models.Property]
	RightMoveService  services.RightMoveService
}

func ConstructGetPropertiesHandler(repositoryService services.MongoRepository[models.Property], rightMoveService services.RightMoveService) *TransformProperties {
	return &TransformProperties{
		RepositoryService: repositoryService,
		RightMoveService:  rightMoveService,
	}
}

func (tp *TransformProperties) TransformPropertiesHandler(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		c.JSON(400, gin.H{"error": "Location query parameter is required"})
		return
	}

	locationIds, err := tp.RightMoveService.GetLocationIdentifiers(location)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to fetch location identifiers"})
		return
	}

	var properties []models.Property

	for _, locationId := range locationIds {
		filter := bson.M{
			"locationId": locationId,
		}

		result, err := tp.RepositoryService.FindByFilter(filter)
		if err != nil {
			fmt.Printf("Failed to fetch properties for location %s: %v", locationId, err)
			continue
		}

		if len(result) == 0 || result == nil {
			fmt.Printf("No properties found for location %s\n", locationId)
			continue
		}

		properties = append(properties, result...)
	}
	c.JSON(200, gin.H{"properties": properties})
}
