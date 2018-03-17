package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
	"log"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update existing note",
	Example: "update-i id -t title_1 -tags tag1,tag2 -n notebook_1",
	Run:     updateWrapper,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Int64P("id", "i", 0, "Id of note to be updated.")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().StringP("title", "t", "", "Notes title.")
	updateCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	updateCmd.Flags().StringP("notebook", "n", "", "Notebook that this note belongs to")
}

func updateWrapper(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetInt64("id")
	title, _ := cmd.Flags().GetString("title")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	notebookTitle, _ := cmd.Flags().GetString("notebook")
	update(id, title, tags, notebookTitle, viEditor)
}

func update(id int64, title string, tags []string, notebookTitle string, editor func(text string) string) {
	note, err := NoteDB.GetNote(id)
	if err != nil {
		log.Panicf("Error while retrieving Note from DB, error msg: %v", err)
	}
	err = constructUpdatedNote(note, title, notebookTitle, tags, editor)
	if err != nil {
		log.Panicf("Error while constructing updated note, error msg: %v", err)
	}
	err = NoteDB.UpdateNote(note)
	if err != nil {
		log.Panicf("Error while updating note, error msg: %v", err)
	}
}

func constructUpdatedNote(note *model.Note, title, notebookTitle string, tags []string, editor func(text string) string) error {
	memo := editor(note.Memo)
	if title != "" {
		note.UpdateTitle(title)
	}
	note.UpdateMemo(memo)
	note.UpdateTags(tags)
	if notebookTitle != "" {
		err := addNotebookToNote(note, notebookTitle)
		if err != nil {
			return err
		}
	}
	return nil
}
