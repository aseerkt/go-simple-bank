package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type userResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func getUserResponse(user db.User) *userResponse {
	return &userResponse{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}
}

type createUserPayload struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *Server) createUser(c *gin.Context) {
	var payload createUserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		handleBadRequest(c, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		handleBadRequest(c, fmt.Errorf("failed to hash password: %s", err))
		return
	}

	arg := db.CreateUserParams{
		Username:       payload.Username,
		FullName:       payload.FullName,
		Email:          payload.Email,
		HashedPassword: string(hashedPassword),
	}

	user, err := s.store.CreateUser(c, arg)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				handleBadRequest(c, fmt.Errorf("user already exists: %s", err))
				return
			}
		}
		handleInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, getUserResponse(user))
}

type loginUserPayload struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) loginUser(c *gin.Context) {
	var payload loginUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handleBadRequest(c, err)
		return
	}

	user, err := s.store.GetUser(c, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			handleNotFound(c, err)
			return
		}
		handleInternalError(c, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(payload.Password)); err != nil {
		handleUnauthorized(c, err)
		return
	}

	token, err := s.tokenMaker.CreateToken(user.Username, 12*time.Hour)

	if err != nil {
		handleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  getUserResponse(user),
	})
}
