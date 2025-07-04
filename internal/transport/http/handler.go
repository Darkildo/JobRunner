package http

import (
	"JobRunner/internal/transport"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service transport.Service
}

func NewTaskHandler(service transport.Service) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	task := h.service.CreateTask()
	h.service.Enqueue(task.ID)

	c.JSON(http.StatusAccepted, gin.H{"id": task.ID})
}

func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("id")

	task, err := h.service.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	response := gin.H{
		"id":         task.ID,
		"status":     string(task.State),
		"created_at": task.CreatedAt.Format(time.RFC3339),
		"duration":   task.Duration(),
	}

	if task.StartedAt != nil {
		response["started_at"] = task.StartedAt.Format(time.RFC3339)
	}

	if task.FinishedAt != nil {
		response["finished_at"] = task.FinishedAt.Format(time.RFC3339)
	}

	if task.Result != "" {
		response["result"] = task.Result
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	if err := h.service.CancelTask(taskID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := h.service.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
