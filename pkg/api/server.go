package api

import (
	"net/http"

	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.SQLStore
	router *gin.Engine
}

func NewServer(store *db.SQLStore) *Server {
	return &Server{store: store, router: gin.Default()}
}

func (s *Server) LoadRoutes() {
	s.router.POST("/accounts", s.createAccount)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.GET("/accounts", s.listAccounts)
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

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
