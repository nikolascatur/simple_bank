package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferResult, error) {
	var result TransferResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Println(txName, "get Acoount for update 1")
		// fromAccount, err := q.GetAccountForUpdate(ctx, args.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(txName, "update account 1")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      args.FromAccountID,
		// 	Balance: fromAccount.Balance - args.Amount,
		// })

		if args.FromAccountID < args.ToAccountID {
			// fmt.Println(txName, "get Acoount for update 2")
			// toAccount, err := q.GetAccountForUpdate(ctx, args.ToAccountID)
			// if err != nil {
			// 	return err
			// }

			// fmt.Println(txName, "update account 2")
			// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			// 	ID:      args.ToAccountID,
			// 	Balance: toAccount.Balance + args.Amount,
			// })

			// fmt.Println(txName, "update account 1")
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     args.FromAccountID,
			// 	Amount: -args.Amount,
			// })

			// if err != nil {
			// 	return err
			// }

			// fmt.Println(txName, "update account 2")
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     args.ToAccountID,
			// 	Amount: args.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
		} else {
			// fmt.Println(txName, "update account 2")
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     args.ToAccountID,
			// 	Amount: args.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			// fmt.Println(txName, "update account 1")
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     args.FromAccountID,
			// 	Amount: -args.Amount,
			// })

			// if err != nil {
			// 	return err
			// }
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amountID1 int64,
	accountID2 int64,
	amountID2 int64,
) (account1 Accounts, account2 Accounts, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amountID2,
	})
	return
}
