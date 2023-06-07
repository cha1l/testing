package handler

import (
	"constester-go/internal/repository"
	"github.com/gin-gonic/gin"
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

	h.SuccessfulResponse(c, "the task was successfully inserted")
}
