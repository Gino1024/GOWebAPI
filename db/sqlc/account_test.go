package db

import (
	"context"
	"database/sql"
	"simple_blank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
			Owner: utils.RandomOwner(),
			Balance: utils.RandomMoney(),
			Currency: utils.RandomCurrency(),
	}
	account, err := testQuries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt   )
	return account
}

func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	acc1 := createRandomAccount(t)
	acc2, err := testQuries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t,acc2)

	require.Equal(t,acc1.ID, acc2.ID)
	require.Equal(t,acc1.Owner, acc2.Owner)
	require.Equal(t,acc1.Balance, acc2.Balance)
	require.Equal(t,acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestUpdate(t *testing.T){
	acc1 := createRandomAccount(t)
	param := UpdateAccountParams{ ID: acc1.ID, Balance: utils.RandomMoney() }
	acc2,err := testQuries.UpdateAccount(context.Background(), param)
	require.NoError(t, err)
	require.NotEmpty(t,acc2)

	require.Equal(t,acc1.ID, acc2.ID)
	require.Equal(t,acc1.Owner, acc2.Owner)
	require.Equal(t,param.Balance, acc2.Balance)
	require.Equal(t,acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDelete(t *testing.T){
	acc1 := createRandomAccount(t)
	err := testQuries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQuries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T){
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQuries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _,acc := range accounts{
		require.NotEmpty(t, acc)
	}
}