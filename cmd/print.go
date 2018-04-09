package cmd

import (
	"github.com/marcusolsson/tui-go"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"sort"
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
		jNotes, err := transformNotes2JSONNotes(noteMap2Slice(notes))
		if err != nil {
			log.Panicln(err)
		}
		printNotes2Terminal(jNotes)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma separated list of note ids.")
	printCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	printCmd.Flags().StringSliceP("notebook", "n", []string{}, "Comma separated list of notebook titles")
	printCmd.Flags().BoolP("all", "a", false, "Print all notes")
}

func printNotes2Terminal(jNotes []*jsonNote) {
	if len(jNotes) <= 0 {
		return
	}
	//Sort By date (descenting)
	sort.Slice(jNotes, func(i, j int) bool{
		return jNotes[i].LastUpdated.After(jNotes[j].LastUpdated)
	})

	theme := tui.NewTheme()
	selected := tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack}
	theme.SetStyle("table.cell.selected", selected)

	notesInfoHeader := tui.NewTable(0, 0)
	notesInfoHeader.SetColumnStretch(0, 1)
	notesInfoHeader.SetColumnStretch(1, 2)
	notesInfoHeader.SetColumnStretch(2, 3)
	notesInfoHeader.SetColumnStretch(3, 2)

	labelID:=tui.NewLabel("ID")
	labelID.SetStyleName("header")
	labelNT:= tui.NewLabel("Notebook Title")
	labelNT.SetStyleName("header")
	labelT:=tui.NewLabel("Note Title")
	labelT.SetStyleName("header")
	labelTags:=tui.NewLabel("Tags")
	labelTags.SetStyleName("header")
	notesInfoHeader.AppendRow(labelID, labelNT, labelT, labelTags)
	//unselect it
	notesInfoHeader.Select(-1)

	theme.SetStyle("label.header", tui.Style{Bold: tui.DecorationOn, Bg: tui.ColorDefault, Fg: tui.ColorBlue})


	notesInfo := tui.NewTable(0, 0)
	notesInfo.SetColumnStretch(0, 1)
	notesInfo.SetColumnStretch(1, 2)
	notesInfo.SetColumnStretch(2, 3)
	notesInfo.SetColumnStretch(3, 2)
	notesInfo.SetFocused(true)

	footer := tui.NewLabel("Press 'U' to update a note, 'D' to delete it or 'Esc' to exit.")
	footer.SetStyleName("footer")
	theme.SetStyle("label.footer", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})

	for _, note := range jNotes {
		notesInfo.AppendRow(
			//Note ID
			tui.NewLabel(strconv.Itoa(int(note.ID))),
			//Notebook title
			tui.NewLabel(note.NotebookTitle),
			//Note title
			tui.NewLabel(note.Title),
			//Note tags
			tui.NewLabel(strings.Join(note.Tags, ",")),
		)
	}
	memo := tui.NewLabel("")
	memo.SetSizePolicy(tui.Expanding, tui.Expanding)

	var created = tui.NewLabel("")
	var lastUpdated = tui.NewLabel("")

	dates := tui.NewTable(0, 0)
	dates.AppendRow(tui.NewLabel("Created:"), created, tui.NewLabel("Last Updated:"), lastUpdated)
	dates.SetSizePolicy(tui.Expanding, tui.Minimum)
	//unselect it
	dates.Select(-1)

	mainPart := tui.NewVBox(memo, dates)
	mainPart.SetSizePolicy(tui.Expanding, tui.Expanding)
	mainPart.SetBorder(true)
	mainPart.SetTitle("Note")

	notesInfo.OnSelectionChanged(func(t *tui.Table) {
		n := jNotes[t.Selected()]
		created.SetText(n.Created.Format("Jan 2 2006 15:04"))
		lastUpdated.SetText(n.LastUpdated.Format("Jan 2 2006 15:04"))
		memo.SetText(n.Memo)
	})
	notesInfo.Select(0)	
	
	scrollNotesInfo := tui.NewScrollArea(notesInfo)
	scrollNotesInfo.SetSizePolicy(tui.Maximum, tui.Minimum)

	root := tui.NewVBox(notesInfoHeader, scrollNotesInfo, mainPart, footer)
	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}
	ui.SetTheme(theme)
	//TODO: stop scrolling at the start/end
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+E", func() { 
		scrollNotesInfo.Scroll(0, -1) 
	})
	ui.SetKeybinding("Ctrl+Y", func() { 
		scrollNotesInfo.Scroll(0, 1) 
	})

	ui.SetKeybinding("Up", func() {
		if notesInfo.Selected() != 0 {
			notesInfo.Select(notesInfo.Selected() - 1) 
		}
	})
	ui.SetKeybinding("Down", func() { 
		if notesInfo.Selected() < len(jNotes) -1{
			notesInfo.Select(notesInfo.Selected() + 1) 
		}
	})

	ui.SetKeybinding("D", func() {
		toBeDeleted := jNotes[notesInfo.Selected()]
		notesInfo.RemoveRow(notesInfo.Selected()) 
		delete([]int64{toBeDeleted.ID})
	 })

	 ui.SetKeybinding("U", func() {
		toBeUpdated := jNotes[notesInfo.Selected()]
		ui.Quit()
		update(toBeUpdated.ID, toBeUpdated.Title, toBeUpdated.Tags, toBeUpdated.NotebookTitle, viEditor)
	 })


	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
