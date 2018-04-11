package cmd

import (
	"github.com/marcusolsson/tui-go"
	"github.com/spf13/cobra"
	"log"
	"sort"
	"strconv"
	"strings"
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
	sort.Slice(jNotes, func(i, j int) bool {
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

	labelID := tui.NewLabel("ID")
	labelID.SetStyleName("header")
	labelNT := tui.NewLabel("Notebook Title")
	labelNT.SetStyleName("header")
	labelT := tui.NewLabel("Note Title")
	labelT.SetStyleName("header")
	labelTags := tui.NewLabel("Tags")
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

	footer := tui.NewLabel("Press 'Ctrl+U' to update a note, 'Ctrl+D' to delete it or 'Esc' to exit.")
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
	//DefaultFocusChain not working so this is a work around
	formItems := make([]tui.Widget, 0)

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	ui.SetKeybinding("Up", func() {
		windowSize := scrollNotesInfo.SizeHint().Y
		index := notesInfo.Selected()
		if index > 0 {
			notesInfo.Select(index - 1)
			if index > windowSize {
				scrollNotesInfo.Scroll(0, -1)
			} else {
				scrollNotesInfo.ScrollToTop()
			}
		}
	})
	ui.SetKeybinding("Down", func() {
		windowSize := scrollNotesInfo.SizeHint().Y
		index := notesInfo.Selected()
		if index < len(jNotes)-1 {
			notesInfo.Select(index + 1)
			if index < len(jNotes)-windowSize {
				scrollNotesInfo.Scroll(0, 1)
			} else {
				scrollNotesInfo.ScrollToBottom()
			}
		}
	})

	ui.SetKeybinding("Ctrl+D", func() {
		toBeDeleted := jNotes[notesInfo.Selected()]
		notesInfo.RemoveRow(notesInfo.Selected())
		delete([]int64{toBeDeleted.ID})
	})

	ui.SetKeybinding("Tab", func() {
		var current int
		var currentWidget tui.Widget
		if formItems == nil || len(formItems) == 0 {
			return
		}
		for i, w := range formItems {
			if w != nil && w.IsFocused() {
				current = i
				currentWidget = w
			}
		}
		if current != len(formItems)-1 {
			formItems[current+1].SetFocused(true)
		} else {
			formItems[0].SetFocused(true)
		}
		currentWidget.SetFocused(false)
	})

	ui.SetKeybinding("Ctrl+U", func() {
		toBeUpdated := jNotes[notesInfo.Selected()]

		title := tui.NewEntry()
		title.SetFocused(true)
		title.SetText(toBeUpdated.Title)
		formItems = append(formItems, title)

		notebookTitle := tui.NewEntry()
		notebookTitle.SetText(toBeUpdated.NotebookTitle)
		formItems = append(formItems, notebookTitle)

		form := tui.NewGrid(0, 0)
		form.AppendRow(tui.NewLabel("Note title:"), title)
		form.AppendRow(tui.NewLabel("Notebook title:"), notebookTitle)
		form.SetSizePolicy(tui.Maximum, tui.Maximum)

		//TODO: focused button must change color
		continueBtn := tui.NewButton("[Continue]")
		formItems = append(formItems, continueBtn)
		continueBtn.OnActivated(func(b *tui.Button) {
			ui.Quit()
			update(toBeUpdated.ID, title.Text(), toBeUpdated.Tags, notebookTitle.Text(), viEditor)
		})

		cancelBtn := tui.NewButton("[Cancel]")
		formItems = append(formItems, cancelBtn)
		cancelBtn.OnActivated(func(b *tui.Button) {
			ui.SetWidget(root)
		})

		buttons := tui.NewHBox(
			tui.NewSpacer(),
			tui.NewPadder(5, 3, continueBtn),
			tui.NewPadder(5, 3, cancelBtn),
		)

		window := tui.NewVBox(form, tui.NewSpacer(), buttons)
		window.SetBorder(true)
		window.SetSizePolicy(tui.Expanding, tui.Expanding)
		window.SetSizePolicy(tui.Expanding, tui.Expanding)
		ui.SetWidget(window)

	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
