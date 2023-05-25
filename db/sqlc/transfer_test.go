package db

import (
	"context"
	"testing"
	"time"

	"github.com/kvnyijia/bank-app/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, randomAccount1 Account, randomAccount2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: randomAccount1.ID,
		ToAccountID:   randomAccount2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	randomAccount1 := createRandomAccount(t)
	randomAccount2 := createRandomAccount(t)
	createRandomTransfer(t, randomAccount1, randomAccount2)
}

func TestGetTransfer(t *testing.T) {
	randomAccount1 := createRandomAccount(t)
	randomAccount2 := createRandomAccount(t)
	randomTransfer := createRandomTransfer(t, randomAccount1, randomAccount2)

	transfer, err := testQueries.GetTransfer(context.Background(), randomTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randomTransfer.ID, transfer.ID)
	require.Equal(t, randomTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, randomTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, randomTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, randomTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestListTarnsfer(t *testing.T) {
	lenOfTransfers := 5
	randomAccount1 := createRandomAccount(t)
	randomAccount2 := createRandomAccount(t)
	for i := 0; i < lenOfTransfers; i++ {
		createRandomTransfer(t, randomAccount1, randomAccount2)
		createRandomTransfer(t, randomAccount2, randomAccount1)
	}

	arg := ListTarnsferParams{
		FromAccountID: randomAccount1.ID,
		ToAccountID:   randomAccount1.ID,
		Limit:         int32(lenOfTransfers),
		Offset:        int32(lenOfTransfers),
	}
	listOfTransfers, err := testQueries.ListTarnsfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, listOfTransfers, lenOfTransfers)

	for _, transfer := range listOfTransfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == randomAccount1.ID || transfer.ToAccountID == randomAccount1.ID)
	}
}
