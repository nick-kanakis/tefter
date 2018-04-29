package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var deleteNotebooksCmd = &cobra.Command{
	Use:     "deleteNotebook",
	Short:   "Delete one or more notebooks based on their title",
	Example: "deleteNotebook notebook notebook2...",
	Args:    cobra.MinimumNArgs(1),
	Run:     deleteNotebooksWrapper,
}

func deleteNotebooksWrapper(cmd *cobra.Command, args []string) {
	if err := deleteNotebooks(args); err != nil {
		log.Fatalln(err)
	}
}

func deleteNotebooks(titles []string) error {
	if len(titles) <= 0 {
		return errors.New("No argument passed, at least one notebook title should be provided")
	}

	for _, notebookTitle := range titles {
		notebook, err := NotebookDB.GetNotebookByTitle(notebookTitle)
		if err != nil || notebook == nil {
			return fmt.Errorf("Could not retrieve notebook for title: %v error msg: %v", notebookTitle, err)
		}
		NotebookDB.DeleteNotebook(notebook.ID)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deleteNotebooksCmd)
}
