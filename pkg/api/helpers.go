package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func handleBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, errorResponse(err))
}

func handleInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, errorResponse(err))
}

func handleNotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, errorResponse(err))
}

func handleForbidden(c *gin.Context, err error) {
	c.JSON(http.StatusForbidden, errorResponse(err))
}

func handleUnauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, errorResponse(err))
}
