package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountPayload struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(c *gin.Context) {
	var payload createAccountPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		handleBadRequest(c, err)
		return
	}

	authPayload := getAuthCtx(c)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: payload.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(c, arg)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			fmt.Println(err)
			fmt.Println(pqError.Code.Name())
			switch pqError.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				handleForbidden(c, pqError)
				return
			}
		}
		handleInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, account)
}

type getAccountUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(c *gin.Context) {
	var uri getAccountUri

	if err := c.ShouldBindUri(&uri); err != nil {
		handleBadRequest(c, err)
		return
	}

	account, err := s.store.GetAccount(c, uri.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			handleNotFound(c, err)
		} else {
			handleInternalError(c, err)
		}
		return
	}

	authPayload := getAuthCtx(c)

	if authPayload.Username != account.Owner {
		err := errors.New("account doesn't belong to current user")
		handleUnauthorized(c, err)
		return
	}

	c.JSON(http.StatusOK, account)
}

type listAccountsQuery struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=20"`
}

func (s *Server) listAccounts(c *gin.Context) {
	var query listAccountsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		handleBadRequest(c, err)
		return
	}

	authPayload := getAuthCtx(c)

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Offset: int32(query.PageID - 1),
		Limit:  int32(query.PageSize),
	}

	accounts, err := s.store.ListAccounts(c, arg)

	if err != nil {
		handleInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, accounts)
}
