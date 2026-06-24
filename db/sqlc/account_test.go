package db

import (
	"context"
	"testing"

	"github.com/nikolascatur/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Accounts {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	oneAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, oneAccount.Owner, account.Owner)
	require.Equal(t, oneAccount.ID, account.ID)
	require.Equal(t, oneAccount.Balance, account.Balance)
	require.Equal(t, oneAccount.Currency, account.Currency)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	updateAccount, errUpdate := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{ID: account.ID, Balance: 20})
	require.NoError(t, errUpdate)
	require.Equal(t, updateAccount.Balance, int64(20))
	require.Equal(t, updateAccount.Owner, account.Owner)
	require.Equal(t, updateAccount.ID, account.ID)
	require.Equal(t, updateAccount.Currency, account.Currency)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	errDel := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, errDel)
	account, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  5,
		Offset: 5,
	})

	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
