package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all functions to excute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx, nil) 
	if err != nil{
		return err
	}
	//取得tx(dbtx)再new出一個包含transaction的Quries
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err: %v", err, rbErr)
		}
		return err
	}
	
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Ammount int64 `json:"ammount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry`
	ToEntry Entry `json:"to_entry`
}


func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	var result TransferTxResult
	var err error


	query := func(q *Queries) error {
		//此時result 為callback外的變數

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Ammount: arg.Ammount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Ammount: -arg.Ammount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Ammount: arg.Ammount,
		})
		if err != nil {
			return err
		}

		// TODO: update account's balance.
		if(arg.FromAccountID < arg.ToAccountID){
			result.FromAccount, result.ToAccount, err = 
				addMoney(ctx, q, arg.FromAccountID, -arg.Ammount, arg.ToAccountID, arg.Ammount)

			if err != nil{
				return err
			}
		}else{
			result.ToAccount, result.FromAccount, err = 
				addMoney(ctx, q, arg.ToAccountID, arg.Ammount, arg.FromAccountID, -arg.Ammount)

			if err != nil{
				return err
			}
		}



		return nil
	}

	err = store.execTx(ctx,query)

	return result, err
}

func addMoney(ctx context.Context,
	q *Queries,
	accountID1 int64, 
	ammount1 int64,
	accountID2 int64, 
	ammount2 int64) (account1 Account, account2 Account, err error){
		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: accountID1,
			Amount: ammount1,
		})
		if err != nil {
			return
		}
		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: accountID2,
			Amount: ammount2,
		})
		if err != nil {
			return
		}

		return account1, account2, nil
}