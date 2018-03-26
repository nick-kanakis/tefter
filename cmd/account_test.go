package cmd

import (
	"testing"
	"github.com/nicolasmanic/tefter/repository"
)

func TestGetAccounts(t *testing.T){
	originalAccountDB := AccountDB 
	AccountDB = mockAccountDB{}
	defer func(){
		AccountDB = originalAccountDB
	}()
	getAccounts(nil, []string{})
}

func TestGetAccountsEmptyResultWillNotFail(t *testing.T){
	originalAccountDB := AccountDB 
	AccountDB = mockAccountDBReturnEmpty{}
	defer func(){
		AccountDB = originalAccountDB
	}()
	getAccounts(nil, []string{})
}

type mockAccountDB struct{
	repository.AccountRepository
} 

func (mDB mockAccountDB) GetUsernames() []string{
	return []string{"username1", "username2"}
}

type mockAccountDBReturnEmpty struct{
	repository.AccountRepository
} 

func (mDB mockAccountDBReturnEmpty) GetUsernames() []string{
	return []string{}
}