package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search notes given a keyword",
	Long:  "Keyword is searched against title and content of the note, if no keyword is given all notes will be printed",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := ""
		if len(args) > 0 {
			keyword = args[0]
		}
		var notes []*model.Note
		var err error
		if len(keyword) == 0 {
			notes, err = NoteDB.GetNotes([]int64{})
		} else {
			notes, err = NoteDB.SearchNotesByKeyword(keyword)
		}
		if err != nil {
			//TODO handle the error
			fmt.Printf("error msg: %v", err)
		}
		printNotes2Terminal(notes)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
