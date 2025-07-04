package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *TaskHandler) *gin.Engine {
	router := gin.Default()
	
	v1 := router.Group("/api/v1")

	v1.POST("/tasks", handler.CreateTask)
	v1.GET("/tasks/:id", handler.GetTaskStatus)
	v1.DELETE("/tasks/:id", handler.DeleteTask)

	return router
}
