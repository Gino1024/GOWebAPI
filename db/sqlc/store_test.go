package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTeansferTx(t *testing.T){	
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	// run n concurrent transfer transactions
	n := 5
	ammount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	fmt.Printf("before: acc1.Balance: %v, acc2.Balance: %v \n", acc1.Balance, acc2.Balance)

	for i:=0; i < n; i++{
		go func(){
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID: acc2.ID,
				Ammount: ammount,
			})
			//chan
			errs <-err
			results <-result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i:=0; i < n;i++{
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, ammount, transfer.Ammount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.Equal(t, fromEntry.AccountID, acc1.ID)
		require.Equal(t, fromEntry.Ammount, -ammount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.Equal(t, toEntry.AccountID, acc2.ID)
		require.Equal(t, toEntry.Ammount, ammount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		
		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount) 
		require.Equal(t, acc1.ID, fromAccount.ID)


		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount) 
		require.Equal(t, acc2.ID, toAccount.ID)

		// check accounts balance
	  fmt.Printf("tx: acc1.Balance: %v, acc2.Balance: %v \n", fromAccount.Balance, toAccount.Balance)
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t,diff1 > 0)
		require.True(t,diff2 > 0)
 		require.True(t, diff1%ammount == 0) // the diff1 is a multiple of the difference

		k := int(diff1 / ammount)
		require.True(t, k>=1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAcount1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	updateAcount2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance - (int64(n) * ammount), updateAcount1.Balance)
	require.Equal(t, acc2.Balance + (int64(n) * ammount), updateAcount2.Balance)
	fmt.Printf("after: acc1.Balance: %v, acc2.Balance: %v \n", updateAcount1.Balance, updateAcount2.Balance)
}

//https://wiki.postgresql.org/wiki/Lock_Monitoring
func TestTeansferTxDeadLock(t *testing.T){	
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	// run n concurrent transfer transactions
	n := 10
	// Half of the transactions transfer funds from Account1 to Account2
	// and the other half transactions transfer funds from account2 to account1
	ammount := int64(10)
	errs := make(chan error)
	fmt.Printf("before: acc1.Balance: %v, acc2.Balance: %v \n", acc1.Balance, acc2.Balance)

	for i:=0; i < n; i++{
		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i% 2 == 1 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}

		go func(){
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Ammount: ammount,
			})
			//chan
			errs <-err
		}()
	}

	for i:=0; i < n;i++{
		err := <-errs
		require.NoError(t, err)
	}

	updateAcount1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	updateAcount2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, updateAcount1.Balance)
	require.Equal(t, acc2.Balance, updateAcount2.Balance)
	fmt.Printf("after: acc1.Balance: %v, acc2.Balance: %v \n", updateAcount1.Balance, updateAcount2.Balance)
}