package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	s := NewStore(testDB)

	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        10,
	}

	n := 5

	results := make(chan TransferTxResult)
	errs := make(chan error)

	for range n {
		go func() {
			result, err := s.TransferTx(context.Background(), arg)

			results <- result
			errs <- err
		}()
	}

	for i := range n {
		result, err := <-results, <-errs

		require.NoError(t, err)
		require.NotEmpty(t, result)

		require.NotEmpty(t, result.Transfer.ID)
		require.NotEmpty(t, result.Transfer.CreatedAt)
		require.Equal(t, result.Transfer.Amount, arg.Amount)

		_, err = s.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.FromEntry)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreateAt)
		require.Equal(t, result.FromEntry.AccountID, fromAccount.ID)
		require.Equal(t, result.FromEntry.Amount, -arg.Amount)

		_, err = s.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.ToEntry)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreateAt)
		require.Equal(t, result.ToEntry.AccountID, toAccount.ID)
		require.Equal(t, result.ToEntry.Amount, arg.Amount)

		_, err = s.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// from account
		require.NotEmpty(t, result.FromAccount)
		require.NotZero(t, result.FromAccount.ID)
		require.NotZero(t, result.FromAccount.CreatedAt)
		require.Equal(t, fromAccount.ID, result.FromAccount.ID)
		require.Equal(t, fromAccount.Balance-int64(i+1)*arg.Amount, result.FromAccount.Balance)

		// from account
		require.NotEmpty(t, result.ToAccount)
		require.NotZero(t, result.ToAccount.ID)
		require.NotZero(t, result.ToAccount.CreatedAt)
		require.Equal(t, toAccount.ID, result.ToAccount.ID)
		require.Equal(t, toAccount.Balance+int64(i+1)*arg.Amount, result.ToAccount.Balance)

	}

}
