package db

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T) Entry {
	account := createTestAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    int64(gofakeit.IntRange(20, 100)),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreateAt)
	require.Equal(t, entry.Amount, arg.Amount)
	require.Equal(t, entry.AccountID, arg.AccountID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createTestEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createTestEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.ID, entry1.ID)
	require.Equal(t, entry2.Amount, entry1.Amount)
	require.Equal(t, entry2.AccountID, entry1.AccountID)
	require.WithinDuration(t, entry2.CreateAt, entry1.CreateAt, time.Second)
}
