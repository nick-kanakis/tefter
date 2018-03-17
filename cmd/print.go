package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print notes based on given ids",
	Long: "There are 4 ways to print a set of notes" +
		" 1) Give a comma separated list of note ids" +
		" 2) Give a comma separated list of notebook titles" +
		" 3) Give a comma separated list of tags," +
		" 4) If -a or --all flag is set all notes will be printed",
	Example: "print -i 1,2,... -n notebook1,notebook2,... -t tag1,tag2,... ",
	Run: func(cmd *cobra.Command, args []string) {

		ids, _ := cmd.Flags().GetIntSlice("ids")
		notebookTitles, _ := cmd.Flags().GetStringSlice("notebook")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		printAll, _ := cmd.Flags().GetBool("all")

		notes := collectNotes(ids, notebookTitles, tags, printAll)

		printNotes2Terminal(noteMap2Slice(notes))
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma separated list of note ids.")
	printCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	printCmd.Flags().StringSliceP("notebook", "n", []string{}, "Comma separated list of notebook titles")
	printCmd.Flags().BoolP("all", "a", false, "Print all notes")
}

//TODO: better visualization of data
func printNotes2Terminal(notes []*model.Note) {
	notebookTitlesMap, err := NotebookDB.GetAllNotebooksTitle()
	if err != nil {
		//TODO handle the error
		fmt.Printf("error msg: %v", err)
	}
	for _, note := range notes {
		fmt.Println("------------------------------")
		fmt.Printf("ID: %v	title: %v	Notebook: %v \n", note.ID, note.Title, notebookTitlesMap[note.NotebookID])
		fmt.Printf("tags: %v \n", note.Tags)
		fmt.Println(note.Memo)
		fmt.Printf("Created: %v	Updated: %v\n", note.Created, note.LastUpdated)
	}
}
