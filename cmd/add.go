package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
)

const DEFAULT_NOTEBOOK_ID = 1

var addNoteCmd = &cobra.Command{
	Use:     "add",
	Short:   "Create a new note",
	Example: "add -t title_1 --tags tag1,tag2 -n notebook_1",
	Run: func(cmd *cobra.Command, args []string) {
		memo := viEditor("")

		title, _ := cmd.Flags().GetString("title")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		notebookTitle, _ := cmd.Flags().GetString("notebook")

		//By default all newNotes will be inserted to default notebook
		//In next steps the notebook may change
		note := model.NewNote(title, memo, DEFAULT_NOTEBOOK_ID, tags)
		err := addNotebookToNote(note, notebookTitle)
		if err != nil {
			//TODO handle the error
			fmt.Printf("error msg: %v", err)
		}

		_, err = NoteDB.SaveNote(note)
		if err != nil {
			//TODO handle the error
			fmt.Printf("error msg: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addNoteCmd)
	addNoteCmd.Flags().StringP("title", "t", "", "Notes title.")
	addNoteCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	addNoteCmd.Flags().StringP("notebook", "n", "", "Notebook that this note belongs to")
}

//If notebookTitle exists it will be inserted there.
//If notebookTitle is empty it will be inserted to the default notebook.
//If notebookTitle does not exists notebook will be created and note will be there.
func addNotebookToNote(note *model.Note, notebookTitle string) error {
	if notebookTitle == "" {
		note.NotebookID = DEFAULT_NOTEBOOK_ID
		return nil
	}

	notebook, err := NotebookDB.GetNotebookByTitle(notebookTitle)
	if err != nil {
		return err
	}

	if notebook == nil {
		newNotebook := model.NewNotebook(notebookTitle)
		id, err := NotebookDB.SaveNotebook(newNotebook)
		if err != nil {
			return err
		}
		note.UpdateNotebook(id)
	} else {
		note.UpdateNotebook(notebook.ID)
	}

	return nil
}
