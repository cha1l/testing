package handler

import (
	"constester-go/internal/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) AddTaskEndpoint(c *gin.Context) {
	var task repository.Task

	if err := c.BindJSON(&task); err != nil {
		h.ApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AddTask(task); err != nil {
		h.ApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusOK)
	log.Infof("Task %s was created\n", task.Name)
}

func (h *Handler) GetAllTasksEndpoint(c *gin.Context) {
	tasks, err := h.service.Tasks.GetAllTasks()
	if err != nil {
		h.ApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tasks)
}
