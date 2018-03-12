package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/repository"
	"github.com/spf13/cobra"
	"os"
)

var (
	NoteDB     repository.NoteRepository
	NotebookDB repository.NotebookRepository

	rootCmd = &cobra.Command{
		Use:   "tefter",
		Short: "Tefter is a simple memo book application",
	}
)

//Execute add all commands to root.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
