package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/repository"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	originalAccountDB := AccountDB
	AccountDB = mockAccountDB{}
	defer func() {
		AccountDB = originalAccountDB
	}()
	getAccounts(nil, []string{})
}

func TestGetAccountsEmptyResultWillNotFail(t *testing.T) {
	originalAccountDB := AccountDB
	AccountDB = mockAccountDBReturnEmpty{}
	defer func() {
		AccountDB = originalAccountDB
	}()
	getAccounts(nil, []string{})
}

func TestGetCredentials(t *testing.T) {
	cases := []struct {
		fpr   FakePasswordReader
		input io.Reader
		cred  *credentials
		err   error
	}{
		{
			fpr: FakePasswordReader{
				[]byte("mypassword"),
				nil,
			},
			cred: &credentials{
				"username1",
				[]byte("mypassword"),
			},
			input: strings.NewReader("username1"),
			err:   nil,
		},
		{
			fpr: FakePasswordReader{
				[]byte("mypassword"),
				nil,
			},
			cred:  &credentials{},
			input: strings.NewReader(""),
			err:   errors.New("Empty username"),
		},
		{
			fpr: FakePasswordReader{
				[]byte("123"),
				nil,
			},
			cred:  &credentials{},
			input: strings.NewReader("username1"),
			err:   errors.New("Password must be at least 5 chars"),
		}, {
			fpr: FakePasswordReader{
				[]byte("123"),
				errors.New("Unexpected Error"),
			},
			cred:  &credentials{},
			input: strings.NewReader(""),
			err:   errors.New("Failed reading password, error msg: Unexpected Error"),
		},
	}

	for _, c := range cases {
		cred, err := getCredentials(c.fpr, c.input)
		if !reflect.DeepEqual(c.err, err) {
			t.Errorf("Expected err to be %q but it was %q", c.err, err)
		}
		if !reflect.DeepEqual(c.cred, cred) {
			t.Errorf("Expected credentials to be %q but it was %q", c.cred, cred)
		}
	}
}

type mockAccountDB struct {
	repository.AccountRepository
}

func (mDB mockAccountDB) GetUsernames() []string {
	return []string{"username1", "username2"}
}

type mockAccountDBReturnEmpty struct {
	repository.AccountRepository
}

func (mDB mockAccountDBReturnEmpty) GetUsernames() []string {
	return []string{}
}

type FakePasswordReader struct {
	password []byte
	err      error
}

func (fpr FakePasswordReader) ReadPassword(fd int) ([]byte, error) {
	if fpr.err != nil {
		return nil, fpr.err
	}
	return fpr.password, nil
}
