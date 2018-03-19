package cmd

import (
	"encoding/json"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "Import notes from json file",
	Args:    cobra.ExactArgs(1),
	Example: "import /c/documents/notes.json ",
	Run:     importNotesWrapper,
}

func importNotesWrapper(cmd *cobra.Command, args []string) {
	importNotes(args)
}

func importNotes(args []string) {
	if len(args) <= 0 {
		log.Panicln("No argument passed, at least one json file path should be provided")
	}
	raw, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Panicf("Error while reading file, error msg: %v", err)
	}

	var jsonNotes []jsonNote
	json.Unmarshal(raw, &jsonNotes)

	for _, jsonNote := range jsonNotes {
		note := model.NewNote(jsonNote.Title, jsonNote.Memo, repository.DEFAULT_NOTEBOOK_ID, jsonNote.Tags)
		addNotebookToNote(note, jsonNote.NotebookTitle)
		NoteDB.SaveNote(note)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
}
