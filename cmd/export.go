package cmd

import (
	"encoding/json"
	"github.com/nicolasmanic/tefter/model"
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
	Long: "There are 4 ways to export a set of notes" +
		" 1) Give a comma separated list of note ids" +
		" 2) Give a comma separated list of notebook titles" +
		" 3) Give a comma separated list of tags," +
		" 4) If -a or --all flag is set all notes will be printed",
	Example: "export -i 1,2,... -n notebook1,notebook2,... -t tag1,tag2,... ",
	Run:     export,
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
	notes := collectNotes(ids, notebookTitles, tags, all)
	export2JSON(noteMap2Slice(notes))
}

func export2JSON(notes []*model.Note) {
	notebookTitlesMap, err := NotebookDB.GetAllNotebooksTitle()
	if err != nil {
		log.Panicf("Error while retrieving Notebooks titles, error msg: %v", err)
	}
	var exportedNotes []*jsonNote
	for _, note := range notes {
		exportedNote := &jsonNote{
			ID:            note.ID,
			Title:         note.Title,
			Memo:          note.Memo,
			Created:       note.Created,
			LastUpdated:   note.LastUpdated,
			Tags:          tagMap2Slice(note.Tags),
			NotebookTitle: notebookTitlesMap[note.NotebookID],
		}
		exportedNotes = append(exportedNotes, exportedNote)
	}
	jsonNotes, err := json.Marshal(exportedNotes)
	if err != nil {
		log.Panicf("Error while marshalling Notes, error msg: %v", err)
	}
	ioutil.WriteFile("notes.json", jsonNotes, 0644)
}
