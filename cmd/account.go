package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
	"syscall"
)

var (
	getCredentials = credentials
	accountCmd     = &cobra.Command{
		Use:   "account",
		Short: "Add/Delete account",
	}
	addAccountCmd = &cobra.Command{
		Use:   "add username",
		Short: "Create new account",
		Run:   addAccount,
	}
	deleteAccountCmd = &cobra.Command{
		Use:   "delete username",
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

func init() {
	accountCmd.AddCommand(addAccountCmd)
	accountCmd.AddCommand(deleteAccountCmd)
	accountCmd.AddCommand(printUsernamesCmd)
	rootCmd.AddCommand(accountCmd)
}

func addAccount(cmd *cobra.Command, args []string) {
	username, password := getCredentials()
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Panicf("Failed hashing password, error msg: %v", err)
	}
	err = AccountDB.CreateAccount(username, hashedPassword)
	if err != nil {
		log.Panicf("Failed creating new account, error msg: %v", err)
	}
}

func deleteAccount(cmd *cobra.Command, args []string) {
	username, password := getCredentials()
	account, err := AccountDB.GetAccount(username)
	if err != nil {
		log.Panicf("Could not delete account for user: %v", username)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), password); err != nil {
		log.Panicln("Username and password don't match")
	}
	AccountDB.DeleteAccount(username)
	fmt.Printf("Account for user: %v deleted", username)
}

func getAccounts(cmd *cobra.Command, args []string) {
	usernames := AccountDB.GetUsernames()
	if len(usernames) == 0 {
		fmt.Println("DB is empty")
	}

	fmt.Println("Current Usernames in DB:")
	for _, username := range usernames {
		fmt.Println("> " + username)
	}
}

func credentials() (string, []byte) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Panicf("Failed reading password, error msg: %v", err)
	}
	password := string(bytePassword)
	if len(password) < 5 {
		log.Panicln("Password must be at least 5 chars")
	}
	if username == "" {
		log.Panicln("Empty username")
	}

	return strings.TrimSpace(username), bytePassword
}
