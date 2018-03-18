package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var deleteNotebookCmd = &cobra.Command{
	Use:     "deleteNotebook",
	Short:   "Delete one or more notebooks based on their title",
	Example: "delete notebook notebook2...",
	Args:    cobra.MinimumNArgs(1),
	Run:     deleteNotebookWrapper,
}

func deleteNotebookWrapper(cmd *cobra.Command, args []string) {
	deleteNotebook(args)
}

func deleteNotebook(args []string) {
	if len(args) <= 0 {
		log.Panicf("No argument passed, at least one notebook title should be provided")
	}

	for _, notebookTitle := range args {
		notebook, err := NotebookDB.GetNotebookByTitle(notebookTitle)
		if err != nil || notebook == nil {
			log.Panicf("Could note get notebook for title: %v error msg: %v", notebookTitle, err)
		}
		NotebookDB.DeleteNotebook(notebook.ID)
	}
}

func init() {
	rootCmd.AddCommand(deleteNotebookCmd)
}
