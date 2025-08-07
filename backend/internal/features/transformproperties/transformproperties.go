package transformproperties

import (
	"github.com/gin-gonic/gin"
	"github.com/greatmove/backend/models"
	"github.com/greatmove/backend/services"
)

type TransformPropertiesHandler interface {
	TransformPropertiesHandler(c *gin.Context)
}

type TransformProperties struct {
	RepositoryService services.MongoRepository[models.Property]
}

func ConstructGetPropertiesHandler(repositoryService services.MongoRepository[models.Property]) *TransformProperties {
	return &TransformProperties{
		RepositoryService: repositoryService,
	}
}

func (tp *TransformProperties) TransformPropertiesHandler(c *gin.Context) {
	testProperty, err := tp.RepositoryService.FindByID("165076340")
	if err != nil {
		c.JSON(404, gin.H{"error": "Property not found"})
		return
	}
	c.JSON(200, gin.H{"message": "TransformPropertiesHandler called", "property": testProperty})
}
