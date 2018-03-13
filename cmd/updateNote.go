package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "updateNote",
	Short:   "Update existing note",
	Example: "updateNote -i id -t title_1 -a tag1,tag2 -n notebook_1",
	Run: func(cmd *cobra.Command, args []string) {
		id,_ := cmd.Flags().GetInt64("id")
		title, _ := cmd.Flags().GetString("title")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		notebookTitle, _ := cmd.Flags().GetString("notebook")

		note, err := NoteDB.GetNote(id)
		if err != nil {
			//TODO handle the error
			fmt.Printf("error msg: %v", err)
		}
		constructUpdatedNote(note, title, notebookTitle, tags)
		NoteDB.UpdateNote(note)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Int64P("id", "i", 0, "Id of note to be updated.")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().StringP("title", "t", "", "Notes title.")
	updateCmd.Flags().StringSliceP("tags", "a", []string{}, "Comma-separated tags of note.")
	updateCmd.Flags().StringP("notebook", "n", "", "Notebook that this note belongs to")
}

func constructUpdatedNote(note *model.Note, title, notebookTitle string, tags []string){
		memo,err := openEditor(note.Memo)
		if err != nil {
			//TODO handle the error
			fmt.Printf("error msg: %v", err)
		}
		if title != ""{
			note.UpdateTitle(title)
		}
		note.UpdateMemo(memo)
		note.UpdateTags(tags)
		if notebookTitle != ""{
			err = addNotebookToNote(note, notebookTitle)
			if err != nil {
				//TODO handle the error
				fmt.Printf("error msg: %v", err)
			}
		}	
}
