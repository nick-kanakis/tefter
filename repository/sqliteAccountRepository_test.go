package repository

import (
	"os"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	testRepo := NewAccountRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	err := testRepo.CreateAccount("nick", "pass123")
	if err != nil {
		t.Errorf("Could not save account to DB, error msg: %v", err)
	}
}

func TestCreateAccountShouldFail(t *testing.T) {
	testRepo := NewAccountRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	err := testRepo.CreateAccount("", "")
	if err.Error() != "Username or/and password are empty" {
		t.Error("Expected error with message: 'Username or/and password are empty'")
	}
}

func TestGetAccount(t *testing.T) {
	testRepo := NewAccountRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	testRepo.CreateAccount("nick1", "pass123")
	testRepo.CreateAccount("nick2", "pass1234")
	account, err := testRepo.GetAccount("nick1")
	if err != nil {
		t.Errorf("Could not retrieve account from DB, error msg: %v", err)
	}
	if account.Password != "pass123" {
		t.Error("Could not correctly retrieve account from DB")
	}
}

func TestDeleteAccount(t *testing.T) {
	testRepo := NewAccountRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	testRepo.CreateAccount("nick1", "pass123")
	testRepo.CreateAccount("nick2", "pass1234")
	err := testRepo.DeleteAccount("nick1")
	if err != nil {
		t.Errorf("Could not delete account from DB, error msg: %v", err)
	}

	account, _ := testRepo.GetAccount("nick1")

	if account != nil {
		t.Error("Could not correctly delete account from DB")
	}
}


func TestGetUsernames(t *testing.T) {
	testRepo := NewAccountRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	testRepo.CreateAccount("nick1", "pass123")
	testRepo.CreateAccount("nick2", "pass1234")
	users := testRepo.GetUsernames()
	
	if len(users) !=2{
		t.Error("Could not correctly retrieve users from DB")
	}
}
