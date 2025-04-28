package gamer

import "github.com/gin-gonic/gin"
import "github.com/RamelHenderson/ggs/utilites"

type Handler struct {
	// Service is the service layer for user operations
	Service *Service
}

func (handler *Handler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, &utilites.JsonResponse{
			Message: "",
		})
	}
}
