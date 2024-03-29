package api

import (
	"database/sql"
	"net/http"

	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/gin-gonic/gin"
)

type createAccountPayload struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=INR USD EUR"`
}

func (s *Server) createAccount(c *gin.Context) {
	var payload createAccountPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		handleBadRequest(c, err)
		return
	}

	arg := db.CreateAccountParams{
		Owner:    payload.Owner,
		Currency: payload.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(c, arg)

	if err != nil {

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

	arg := db.ListAccountsParams{
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
