package api

import (
	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/aseerkt/go-simple-bank/pkg/token"
	"github.com/aseerkt/go-simple-bank/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     *utils.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(store db.Store, config *utils.Config) *Server {

	tm := token.NewPasetoMaker()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return &Server{config: config, tokenMaker: tm, store: store, router: gin.Default()}
}

func (s *Server) LoadRoutes() {
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)

	authRoutes := s.router.Group("/").Use(auth(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccounts)

	authRoutes.POST("/transfers", s.createTransfer)

}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
