package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var( 
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Add/Delete account",
		}
	addAccountCmd = &cobra.Command{
		Use:   "add username",
		Short: "Create new account",
		Args:  cobra.ExactArgs(1),
		Run: addAccount,
		}
	deleteAccountCmd =  &cobra.Command{
		Use:   "delete username",
		Short: "Delete existing account",
		Args:  cobra.NoArgs,
		Run: deleteAccount,
		}
	printUsernamesCmd =  &cobra.Command{
		Use:   "print",
		Short: "Show all usernames",
		Args:  cobra.ExactArgs(1),
		Run: getAccounts,
		}
)

func init() {
	accountCmd.AddCommand(addAccountCmd)
	accountCmd.AddCommand(deleteAccountCmd)
	accountCmd.AddCommand(printUsernamesCmd)
	rootCmd.AddCommand(accountCmd)
}

func addAccount(cmd *cobra.Command, args []string) {
	panic("Not implemented")
}

func deleteAccount(cmd *cobra.Command, args []string) {
	panic("Not implemented")
}

func getAccounts(cmd *cobra.Command, args []string) {
	usernames := AccountDB.GetUsernames()
	if len(usernames) == 0 {
		fmt.Println("DB is empty")
	}
	
	fmt.Println("Current Usernames in DB:")
	for _, username:=range usernames{
		fmt.Println("> "+username)
	}
}

