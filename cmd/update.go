package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update existing note",
	Long: "Select a note to update by providing a valid id (required). \n" +
		"A tag can be removed by providing a '-' before the tag name, eg: \n" +
		"--tags tag1,-tag2 will insert tag1 and remove (if exist) tag2 to the note",
	Example: "update id -t title_1 --tags tag1,-tag2 -n notebook_1",
	Args:    cobra.ExactArgs(1),
	Run:     updateWrapper,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("title", "t", "", "Notes title.")
	updateCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	updateCmd.Flags().StringP("notebook", "n", "", "Notebook that this note belongs to")
}

func updateWrapper(cmd *cobra.Command, args []string) {
	title, _ := cmd.Flags().GetString("title")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	notebookTitle, _ := cmd.Flags().GetString("notebook")
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Panicf("ID could not be converted to integer")
	}
	update(id, title, tags, notebookTitle, viEditor)
}

func update(id int64, title string, tags []string, notebookTitle string, editor func(text string) string) {
	note, err := NoteDB.GetNote(id)
	if err != nil {
		log.Panicf("Error while retrieving Note from DB, error msg: %v", err)
	}
	memo := editor(note.Memo)
	err = constructUpdatedNote(note, title, notebookTitle, tags, memo)
	if err != nil {
		log.Panicf("Error while constructing updated note, error msg: %v", err)
	}
	err = NoteDB.UpdateNote(note)
	if err != nil {
		log.Panicf("Error while updating note, error msg: %v", err)
	}
}

func updateJSONNote(jNote *jsonNote) error {
	note, err := NoteDB.GetNote(jNote.ID)
	if err != nil {
		return fmt.Errorf("Error while retrieving Note from DB, error msg: %v", err)
	}
	err = constructUpdatedNote(note, jNote.Title, jNote.NotebookTitle, jNote.Tags, jNote.Memo)
	if err != nil {
		return fmt.Errorf("Error while constructing updated note, error msg: %v", err)
	}
	err = NoteDB.UpdateNote(note)
	if err != nil {
		return fmt.Errorf("Error while updating note, error msg: %v", err)
	}
	return nil
}

func constructUpdatedNote(note *model.Note, title, notebookTitle string, tags []string, memo string) error {
	if title != "" {
		note.UpdateTitle(title)
	}
	note.UpdateMemo(memo)

	toBeRemoved := []string{}
	toBeAdded := []string{}
	for _, tag := range tags {
		if strings.HasPrefix(tag, "-") {
			toBeRemoved = append(toBeRemoved, tag[1:])
		} else {
			toBeAdded = append(toBeAdded, tag)
		}
	}
	note.RemoveTags(toBeRemoved)
	note.AddTags(toBeAdded)
	if notebookTitle != "" {
		err := addNotebookToNote(note, notebookTitle)
		if err != nil {
			return err
		}
	}
	return nil
}
