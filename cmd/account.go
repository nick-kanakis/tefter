package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
)

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Add/Delete account",
	}
	addAccountCmd = &cobra.Command{
		Use:   "add",
		Short: "Create new account",
		Run:   addAccount,
	}
	deleteAccountCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete existing account",
		Args:  cobra.NoArgs,
		Run:   deleteAccount,
	}
	printUsernamesCmd = &cobra.Command{
		Use:   "print",
		Short: "Show all usernames",
		Run:   getAccounts,
	}
)

type credentials struct {
	username string
	password []byte
}

func init() {
	accountCmd.AddCommand(addAccountCmd)
	accountCmd.AddCommand(deleteAccountCmd)
	accountCmd.AddCommand(printUsernamesCmd)
	rootCmd.AddCommand(accountCmd)
}

func addAccount(cmd *cobra.Command, args []string) {
	pr := terminalPasswordReader{}
	credentials, err := getCredentials(pr, os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(credentials.password, 10)
	if err != nil {
		log.Fatalf("Failed hashing password, error msg: %v", err)
	}
	err = AccountDB.CreateAccount(credentials.username, hashedPassword)
	if err != nil {
		log.Fatalf("Failed creating new account, error msg: %v", err)
	}
}

func deleteAccount(cmd *cobra.Command, args []string) {
	pr := terminalPasswordReader{}
	credentials, err := getCredentials(pr, os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	account, err := AccountDB.GetAccount(credentials.username)
	if err != nil {
		log.Fatalf("Could not delete account for user: %v", credentials.username)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), credentials.password); err != nil {
		log.Fatalln("Username and password don't match")
	}
	AccountDB.DeleteAccount(credentials.username)
	fmt.Printf("Account for user: %v deleted", credentials.username)
}

func getAccounts(cmd *cobra.Command, args []string) {
	usernames := AccountDB.GetUsernames()
	if len(usernames) == 0 {
		fmt.Println("DB is empty")
		return
	}

	fmt.Println("Current Usernames in DB:")
	for _, username := range usernames {
		fmt.Println("> " + username)
	}
}

type passwordReader interface {
	ReadPassword(fd int) ([]byte, error)
}

type terminalPasswordReader struct{}

func (pr terminalPasswordReader) ReadPassword(fd int) ([]byte, error) {
	return terminal.ReadPassword(fd)
}

func getCredentials(pr passwordReader, input io.Reader) (*credentials, error) {
	reader := bufio.NewReader(input)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if len(username) == 0 {
		return &credentials{}, errors.New("Empty username")
	}
	
	fmt.Print("Enter Password: ")
	bytePassword, err := pr.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return &credentials{}, fmt.Errorf("Failed reading password, error msg: %v", err)
	}

	password := string(bytePassword)
	if len(password) < 5 {
		return &credentials{}, errors.New("Password must be at least 5 chars")
	}
	return &credentials{username, bytePassword}, nil
}
