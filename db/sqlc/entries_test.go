package db

import (
	"context"
	"testing"

	"github.com/nikolascatur/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntries(t *testing.T, account Accounts) Entries {
	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntries(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntries(t, account)
	oneEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, oneEntry)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntries(t, account)
	}
	args, err := testQueries.ListEntries(context.Background(), ListEntriesParams{AccountID: account.ID, Limit: 5, Offset: 5})
	require.NoError(t, err)
	require.Len(t, args, 5)
	for _, entry := range args {
		require.NotEmpty(t, entry)
	}
}
