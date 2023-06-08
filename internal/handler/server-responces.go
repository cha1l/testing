package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// TestEndpoint for testing server working
func (h *Handler) TestEndpoint(c *gin.Context) {
	log.Println("This is the test endpoint")
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) ApiErrorResponse(c *gin.Context, code int, error string) {
	data := map[string]interface{}{
		"error": error,
	}
	log.Error(error)
	c.JSON(code, data)
}

func (h *Handler) SuccessfulResponse(c *gin.Context, msg string) {
	data := map[string]interface{}{
		"status": "ok",
	}
	c.JSON(http.StatusOK, data)
	log.Info(msg)
}
