package getproperties

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greatmove/backend/models"
	"github.com/greatmove/backend/services"
)

type GetPropertiesHandler interface {
	GetPropertiesHandler(c *gin.Context)
}

type GetProperties struct {
	RightMoveService  services.RightMoveService
	RepositoryService services.MongoRepository[models.Property]
}

func ConstructGetPropertiesHandler(rightMoveService services.RightMoveService, repositoryService services.MongoRepository[models.Property]) *GetProperties {
	return &GetProperties{
		RightMoveService:  rightMoveService,
		RepositoryService: repositoryService,
	}
}

func (rm *GetProperties) GetPropertiesHandler(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location query parameter is required"})
		return
	}
	locationIdentifiers, err := rm.RightMoveService.GetLocationIdentifiers(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch location identifiers"})
		return
	}

	type result struct {
		locationId string
		properties []models.Property
		err        error
	}

	resultChannel := make(chan result, len(locationIdentifiers))
	for _, locationId := range locationIdentifiers {
		go func(locId string) {
			props, err := rm.RightMoveService.SearchProperties(locId)
			resultChannel <- result{locationId: locId, properties: props, err: err}
		}(locationId)
	}

	var properties []models.Property
	for i := 0; i < len(locationIdentifiers); i++ {
		res := <-resultChannel
		if res.err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch properties"})
			return
		}
		rm.SaveProperties(&res.properties, &res.locationId)
		properties = append(properties, res.properties...)
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "1.0.0",
		"data":    properties,
	})
}

func (rm *GetProperties) SaveProperties(properties *[]models.Property, locationId *string) {
	for index := range *properties {
		// this line allows me to work with the orignal slice and mutate it how I want it, I cannot use _,proper : range * properties as this will create a copy of properties
		prop := &(*properties)[index]
		prop.LocationId = locationId
		_, err := rm.RepositoryService.Create(*prop)
		if err != nil {
			//not good because if we fail to save then we lose the property but for now we will we will just log the error and continue.
			fmt.Println("Error saving property:", err)
			continue
		}
	}
}
