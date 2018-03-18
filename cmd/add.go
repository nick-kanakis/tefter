package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
	"log"
)

const DEFAULT_NOTEBOOK_ID = 1

var addNoteCmd = &cobra.Command{
	Use:     "add",
	Short:   "Create a new note",
	Example: "add -t title_1 --tags tag1,tag2 -n notebook_1",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(addNoteCmd)
	addNoteCmd.Flags().StringP("title", "t", "", "Notes title.")
	addNoteCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	addNoteCmd.Flags().StringP("notebook", "n", "", "Notebook that this note belongs to")
}

func addWrapper(cmd *cobra.Command, args []string) {
	title, _ := cmd.Flags().GetString("title")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	notebookTitle, _ := cmd.Flags().GetString("notebook")
	add(title, tags, notebookTitle, args, viEditor)
}

func add(title string, tags []string, notebookTitle string, args []string, editor func(text string) string) {
	memo := editor("")

	//All newNotes will be inserted to default notebook
	//In next steps the notebook may change see addNotebookToNote for more.
	note := model.NewNote(title, memo, DEFAULT_NOTEBOOK_ID, tags)
	err := addNotebookToNote(note, notebookTitle)
	if err != nil {
		log.Panicf("Error whole finding corresponding notebook for note, error msg: %v", err)
	}

	_, err = NoteDB.SaveNote(note)
	if err != nil {
		log.Panicf("Error while saving note, error msg: %v", err)
	}
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
