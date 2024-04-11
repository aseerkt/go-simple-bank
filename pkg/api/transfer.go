package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/gin-gonic/gin"
)

type createTransferPayload struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var payload createTransferPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		handleBadRequest(c, err)
		return
	}

	if !s.hasAccountWithCurrency(c, payload.FromAccountID, payload.Currency) {
		return
	}

	if !s.hasAccountWithCurrency(c, payload.ToAccountID, payload.Currency) {
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: payload.FromAccountID,
		ToAccountID:   payload.ToAccountID,
		Amount:        payload.Amount,
	}

	result, err := s.store.TransferTx(c, arg)

	if err != nil {
		handleInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (s *Server) hasAccountWithCurrency(c *gin.Context, accountId int64, currency string) bool {
	account, err := s.store.GetAccount(c, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			handleNotFound(c, err)
			return false
		}

		handleInternalError(c, err)
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s vs %s", accountId, account.Currency, currency)
		handleBadRequest(c, err)
		return false

	}

	return true
}
