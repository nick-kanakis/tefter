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
	Use:   "import",
	Short: "Import notes from json file",
	Long: "Provide a path of a .json file be imported.\n" +
		"[{\n\t'title':'',\n\t'memo':' ',\n\t'created':'2018-03-19T18:58:29.5553579+02:00',\n\t'updated':'2018-03-19T18:58:29.5553579+02:00',\n\t'tags':[tag1, tag2],\n\t'notebook_title':''\n}]",
	Args:    cobra.ExactArgs(1),
	Example: "import /c/documents/notes.json ",
	Run:     importNotesWrapper,
}

func importNotesWrapper(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Panicln("No argument passed, at least one json file path should be provided")
	}
	importNotes(args[0])
}

func importNotes(path string) {

	raw, err := ioutil.ReadFile(path)
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
