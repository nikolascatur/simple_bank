package db

import (
	"context"
	"testing"

	"github.com/nikolascatur/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount Accounts, toAccount Accounts) Transfers {
	arg, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, arg)
	return arg
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromAccount, toAccount)
	arg, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, arg)
	require.Equal(t, transfer.ID, arg.ID)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
}

func TestListTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount, toAccount)
	}
	arg, err := testQueries.GetListTransfer(context.Background(), GetListTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	})
	require.NoError(t, err)
	require.Len(t, arg, 5)
	for _, transfer := range arg {
		require.NotEmpty(t, transfer)
	}
}
