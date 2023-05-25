package db

import (
	"context"
	"testing"
	"time"

	"github.com/kvnyijia/bank-app/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, randomAccount Account) Entry {
	arg := CreateEntryParams{
		AccountID: randomAccount.ID,
		Amount:    util.RandomInt(-100, 100),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	createRandomEntry(t, randomAccount)
}

func TestGetEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	randomEntry := createRandomEntry(t, randomAccount)
	entry, err := testQueries.GetEntry(context.Background(), randomEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randomEntry.ID, entry.ID)
	require.Equal(t, randomEntry.Amount, entry.Amount)
	require.WithinDuration(t, randomEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	randomAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, randomAccount)
	}

	arg := ListEntriesParams{
		AccountID: randomAccount.ID,
		Limit:     5,
		Offset:    5,
	}

	listOfEntries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, listOfEntries, 5)
	for _, entry := range listOfEntries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
