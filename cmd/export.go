package cmd

import (
	"encoding/json"
	"fmt"
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
	Run: exportWrapper,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma separated list of note ids.")
	exportCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	exportCmd.Flags().StringSliceP("notebook", "n", []string{}, "Comma separated list of notebook titles")
	exportCmd.Flags().BoolP("all", "a", false, "Export all notes")
}

func exportWrapper(cmd *cobra.Command, args []string) {
	ids, _ := cmd.Flags().GetIntSlice("ids")
	notebookTitles, _ := cmd.Flags().GetStringSlice("notebook")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	all, _ := cmd.Flags().GetBool("all")
	if err := export(ids, notebookTitles, tags, all); err != nil {
		log.Fatalln(err)
	}
}

func export(ids []int, notebookTitles, tags []string, getAll bool) error {
	jNotes, err := retrieveJSONNotes(ids, notebookTitles, tags, getAll)
	if err != nil {
		return err
	}
	return writeNotes(jNotes)
}

func retrieveJSONNotes(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
	notes, err := collectNotesFromDB(ids, notebookTitles, tags, getAll)
	if err != nil {
		return nil, err
	}
	jNotes, err := transformNotes2JSONNotes(noteMap2Slice(notes))
	if err != nil {
		return nil, err
	}
	return jNotes, nil
}

func writeNotes(jsonNotes []*jsonNote) error {
	marshalledNotes, err := json.Marshal(jsonNotes)
	if err != nil {
		return fmt.Errorf("Error while marshalling Notes, error msg: %v", err)
	}
	return ioutil.WriteFile("notes.json", marshalledNotes, 0644)
}
