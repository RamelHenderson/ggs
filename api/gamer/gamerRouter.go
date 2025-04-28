package gamer

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the vehicle-related routes with the provided router and database connection.
// r: The Gin engine instance
// db: The GORM database connection
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	gamerGroup := r.Group("/gamer")

	handler := Handler{
		Service: &Service{
			DB: db,
		},
	}

	// Create a new group for vehicle routes
	{
		gamerGroup.GET("/getGamerByID/:id", handler.GetById())
	}
}
