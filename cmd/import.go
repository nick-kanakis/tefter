package cmd

import (
	"encoding/json"
	"fmt"
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
	fsr := fileSystemReader{}
	if err := importNotes(fsr, args[0]); err != nil {
		log.Fatalln(err)
	}
}

type fileReader interface {
	ReadFile(string) ([]byte, error)
}

type fileSystemReader struct{}

func (fsr fileSystemReader) ReadFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func importNotes(fr fileReader, path string) error {
	raw, err := fr.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error while reading file, error msg: %v", err)
	}

	var jsonNotes []jsonNote
	if err = json.Unmarshal(raw, &jsonNotes); err != nil {
		return fmt.Errorf("Could not unmarshal file at path: %v, error msg: %v", path, err)
	}

	for _, jsonNote := range jsonNotes {
		note := model.NewNote(jsonNote.Title, jsonNote.Memo, repository.DEFAULT_NOTEBOOK_ID, jsonNote.Tags)
		err = addNotebookToNote(note, jsonNote.NotebookTitle)
		if err != nil {
			return err
		}
		_, err = NoteDB.SaveNote(note)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(importCmd)
}
