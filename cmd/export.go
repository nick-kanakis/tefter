package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"time"
)

type jsonNote struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Memo          string    `json:"memo"`
	Created       time.Time `json:"created"`
	LastUpdated   time.Time `json:"updated"`
	Tags          []string  `json:"tags"`
	NotebookTitle string    `json:"notebook_title"`
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports notes to json format",
	Long: "There are 4 ways to export a set of notes\n" +
		" 1) Give a comma separated list of note ids\n" +
		" 2) Give a comma separated list of notebook titles\n" +
		" 3) Give a comma separated list of tags,\n" +
		" 4) If -a or --all flag is set all notes will be printed\n",
	Example: "export -i 1,2,... -n notebook1,notebook2,... -t tag1,tag2,...\n " +
		"export -a",
	Run: export,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma separated list of note ids.")
	exportCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	exportCmd.Flags().StringSliceP("notebook", "n", []string{}, "Comma separated list of notebook titles")
	exportCmd.Flags().BoolP("all", "a", false, "Export all notes")
}

func export(cmd *cobra.Command, args []string) {
	ids, _ := cmd.Flags().GetIntSlice("ids")
	notebookTitles, _ := cmd.Flags().GetStringSlice("notebook")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	all, _ := cmd.Flags().GetBool("all")
	jsonNotes, err := retrieveJSONNotes(ids, notebookTitles, tags, all)
	if err != nil {
		log.Panicln(err)
	}
	export2File(jsonNotes)
}

func export2File(jsonNotes []*jsonNote) {
	marshalledNotes, err := json.Marshal(jsonNotes)
	if err != nil {
		log.Panicf("Error while marshalling Notes, error msg: %v", err)
	}
	ioutil.WriteFile("notes.json", marshalledNotes, 0644)
}

func retrieveJSONNotes(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
	notes := collectNotesFromDB(ids, notebookTitles, tags, getAll)
	return transformNotes2JSONNotes(noteMap2Slice(notes))
}
