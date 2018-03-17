package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var updateNotebookCmd = &cobra.Command{
	Use:     "updateNotebook",
	Short:   "Set new title to an existing notebook",
	Example: "update 'Old Notebook Title' 'New Notebook Title'",
	Args:    cobra.ExactArgs(2),
	Run:     updateNotebookWrapper,
}

func init() {
	rootCmd.AddCommand(updateNotebookCmd)
}

func updateNotebookWrapper(cmd *cobra.Command, args []string) {
	updateNotebook(args)
}

func updateNotebook(args []string) {
	if len(args) < 2 {
		log.Panicf("No argument passed, you must pass an old a new notebook title")
	}
	notebook, err := NotebookDB.GetNotebookByTitle(args[0])
	if err != nil {
		log.Panicf("error msg: %v", err)
	} else if notebook != nil {
		notebook.Title = args[1]
		err = NotebookDB.UpdateNotebook(notebook)
		if err != nil {
			log.Panicf("error msg: %v", err)
		}
	}
}
