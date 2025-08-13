package transformproperties

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greatmove/backend/models"
	"github.com/greatmove/backend/services"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

		//work with pagenation here at somepoint but for mvp just use top 10
		options := options.Find().SetLimit(10)
		result, err := tp.RepositoryService.FindByFilter(filter, options)
		if err != nil {
			fmt.Printf("Failed to fetch properties for location %s: %v", locationId, err)
			continue
		}

		if len(result) == 0 || result == nil {
			fmt.Printf("No properties found for location %s\n", locationId)
			continue
		}

		tp.TransformPicture(&result)

		properties = append(properties, result...)
	}

	c.JSON(200, gin.H{"properties": properties})
}

func (rm *TransformProperties) TransformPicture(properties *[]models.Property) {

	for i := range *properties {
		property := &(*properties)[i]
		images := property.PropertyImages.Images
		for j := range images {
			url := images[j].SrcUrl
			images[j].Summary, _ = rm.EncodePropertyPicture(url)
		}
	}
}

func (rm *TransformProperties) EncodePropertyPicture(url *string) (*string, error) {
	if url == nil || *url == "" {
		return nil, fmt.Errorf("invalid image URL")
	}

	resp, err := http.Get(*url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image body: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(body)
	return &encoded, nil
}
