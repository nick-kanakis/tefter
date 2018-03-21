package cmd

import (
	"log"
	"strconv"
	"strings"
	"github.com/marcusolsson/tui-go"
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
		" 4) If -a or --all flag is set all notes will be printed" +
		"Press Esc to exit print mode",
	Example: "print -i 1,2,... -n notebook1,notebook2,... -t tag1,tag2,... ",
	Run: func(cmd *cobra.Command, args []string) {

		ids, _ := cmd.Flags().GetIntSlice("ids")
		notebookTitles, _ := cmd.Flags().GetStringSlice("notebook")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		printAll, _ := cmd.Flags().GetBool("all")

		notes := collectNotesFromDB(ids, notebookTitles, tags, printAll)
		printNotes2Terminal(noteMap2Slice(notes))
		//print2()
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma separated list of note ids.")
	printCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	printCmd.Flags().StringSliceP("notebook", "n", []string{}, "Comma separated list of notebook titles")
	printCmd.Flags().BoolP("all", "a", false, "Print all notes")
}

func printNotes2Terminal(notes []*model.Note) {
	if len(notes) <= 0 {
		return
	}
	notebookTitlesMap, err := NotebookDB.GetAllNotebooksTitle()
	if err != nil {
		log.Panicf("Error while retrieving notebook by title, error msg: %v", err)
	}

	notesInfoHeader := tui.NewTable(0, 0)
	notesInfoHeader.SetColumnStretch(0, 1)
	notesInfoHeader.SetColumnStretch(1, 2)
	notesInfoHeader.SetColumnStretch(2, 3)
	notesInfoHeader.SetColumnStretch(3, 2)
	notesInfoHeader.AppendRow(tui.NewLabel("ID"), tui.NewLabel("Notebook Title"),tui.NewLabel("Note Title"),tui.NewLabel("Tags"))

	notesInfo := tui.NewTable(0, 0)
	notesInfo.SetColumnStretch(0, 1)
	notesInfo.SetColumnStretch(1, 2)
	notesInfo.SetColumnStretch(2, 3)
	notesInfo.SetColumnStretch(3, 2)
	notesInfo.SetFocused(true)

	for _, note := range notes {
		notesInfo.AppendRow(
			//Note ID
			tui.NewLabel(strconv.Itoa(int(note.ID))),
			//Notebook title
			tui.NewLabel(notebookTitlesMap[note.NotebookID]),
			//Note title
			tui.NewLabel(note.Title),
			//Note tags
			tui.NewLabel(strings.Join(tagMap2Slice(note.Tags), ",")),
		)
	}
	memo := tui.NewLabel("")
	memo.SetSizePolicy(tui.Expanding, tui.Expanding)

	var created = tui.NewLabel("")
	var lastUpdated = tui.NewLabel("")

	dates := tui.NewTable(0, 0)
	dates.AppendRow(tui.NewLabel("Created:"), created, tui.NewLabel("Last Updated:"), lastUpdated)
	dates.SetSizePolicy(tui.Expanding, tui.Minimum)

	mainPart := tui.NewVBox(memo, dates)
	mainPart.SetSizePolicy(tui.Expanding, tui.Expanding)

	notesInfo.OnSelectionChanged(func(t *tui.Table) {
		n := notes[t.Selected()]
		created.SetText(n.Created.Format("Jan 2 2006 15:04"))
		lastUpdated.SetText(n.LastUpdated.Format("Jan 2 2006 15:04"))
		memo.SetText(n.Memo)
	})
	notesInfo.Select(0)

	root := tui.NewVBox(notesInfoHeader, notesInfo, tui.NewLabel(""), mainPart)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
