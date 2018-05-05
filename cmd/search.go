package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
	"log"
)

var (
	searchCmd = &cobra.Command{
		Use:     "search",
		Short:   "Search notes given a keyword",
		Long:    "Keyword is searched against title and content of the note, if no keyword is given all notes will be printed",
		Example: "search myKeyword",
		Run:     searchWrapper,
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

func searchWrapper(cmd *cobra.Command, args []string) {
	keyword := ""
	if len(args) > 0 {
		keyword = args[0]
	}
	notes, err := search(keyword)
	if err != nil {
		log.Fatalln(err)
	}
	jNotes, err := transformNotes2JSONNotes(notes)
	if err != nil {
		log.Fatalln(err)
	}
	printNotes2Terminal(jNotes)
}

func search(keyword string) ([]*model.Note, error) {
	var notes []*model.Note
	var err error
	if len(keyword) == 0 {
		notes, err = NoteDB.GetNotes([]int64{})
	} else {
		notes, err = NoteDB.SearchNotesByKeyword(keyword)
	}
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Notes from DB, error msg: %v", err)
	}
	return notes, nil
}
