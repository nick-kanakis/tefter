package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nicolasmanic/tefter/model"
)

type sqliteAccountRepository struct {
	dbPath string
	*sqlx.DB
}

//NewAccountRepository returns a AccountRepository interface
func NewAccountRepository(dbPath string) AccountRepository {
	db := connect2DB(dbPath)
	return &sqliteAccountRepository{dbPath, db}
}

func (accountRepo *sqliteAccountRepository) CreateAccount(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("Username or/and password are empty")
	}

	tx, err := accountRepo.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()
			err = panicErr
		}
	}()

	tx.MustExec(`INSERT INTO account (username, password) VALUES(?,?)`, username, password)
	checkError(err)

	err = tx.Commit()
	checkError(err)
	return err
}

func (accountRepo *sqliteAccountRepository) GetAccount(username string) (*model.Account, error) {

	selectNotebook := "SELECT username, password FROM account WHERE username = ?"
	accounts := []*model.Account{}
	err := accountRepo.Select(&accounts, selectNotebook, []interface{}{username}...)
	checkError(err)

	if len(accounts) == 0 {
		return nil, fmt.Errorf("No account found for username: %v", username)
	}

	return accounts[0], err
}
func (accountRepo *sqliteAccountRepository) DeleteAccount(username string) error {
	deleteAccount := "DELETE FROM account WHERE username = ?"

	tx, err := accountRepo.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()
			err = panicErr
		}
	}()

	tx.MustExec(deleteAccount, []interface{}{username}...)
	err = tx.Commit()
	checkError(err)

	return err
}

func (accountRepo *sqliteAccountRepository) GetUsernames() []string{
	getUsernames := "SELECT username FROM account"
	usernames := []string{}
	err := accountRepo.Select(&usernames, getUsernames, []interface{}{}...)
	checkError(err)
	return usernames
}

func (accountRepo *sqliteAccountRepository) CloseDB() error {
	return accountRepo.Close()
}
