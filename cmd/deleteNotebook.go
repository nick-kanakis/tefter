package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var deleteNotebooksCmd = &cobra.Command{
	Use:     "deleteNotebooks",
	Short:   "Delete one or more notebooks based on their title",
	Example: "deleteNotebooks notebook notebook2...",
	Args:    cobra.MinimumNArgs(1),
	Run:     deleteNotebooksWrapper,
}

func deleteNotebooksWrapper(cmd *cobra.Command, args []string) {
	err := deleteNotebooks(args)
	if err != nil {
		log.Panicln(err)
	}
}

func deleteNotebooks(titles []string) error {
	if len(titles) <= 0 {
		return fmt.Errorf("No argument passed, at least one notebook title should be provided")
	}

	for _, notebookTitle := range titles {
		notebook, err := NotebookDB.GetNotebookByTitle(notebookTitle)
		if err != nil || notebook == nil {
			return fmt.Errorf("Could note get notebook for title: %v error msg: %v", notebookTitle, err)
		}
		NotebookDB.DeleteNotebook(notebook.ID)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deleteNotebooksCmd)
}
