package handler

import (
	"constester-go/internal/docker"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) TestingCodeEndpoint(c *gin.Context) {
	var code docker.Code

	if err := c.BindJSON(&code); err != nil {
		h.ApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	res, err := h.service.RunTests(ctx, code)
	if err != nil {
		h.ApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
