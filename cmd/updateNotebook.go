package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var updateNotebookCmd = &cobra.Command{
	Use:     "updateNotebook",
	Short:   "Set new title to an existing notebook",
	Long:    "Update requires 2 arguments first the old notebook title, and the new title",
	Example: "update 'Old Notebook Title' 'New Notebook Title'",
	Args:    cobra.ExactArgs(2),
	Run:     updateNotebookWrapper,
}

func init() {
	rootCmd.AddCommand(updateNotebookCmd)
}

func updateNotebookWrapper(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		log.Panicf("Incorrect number of arguments passed, you must pass an old a new notebook title")
	}
	if err := updateNotebook(args[0], args[1]); err != nil {
		log.Panicln(err)
	}
}

func updateNotebook(oldTitle, newTitle string) error {
	notebook, err := NotebookDB.GetNotebookByTitle(oldTitle)
	if err != nil {
		return fmt.Errorf("Error while retrieving notebook by title, error msg: %v", err)
	} else if notebook != nil {
		notebook.Title = newTitle
		err = NotebookDB.UpdateNotebook(notebook)
		if err != nil {
			return fmt.Errorf("Error while updating notebook, error msg: %v", err)
		}
	}
	return nil
}
