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
	ui := createUI(jNotes)
	if ui == nil {
		return
	}
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func createUI(jNotes []*jsonNote) tui.UI {
	if len(jNotes) == 0 {
		return nil
	}
	//Sort By date (descenting)
	sort.Slice(jNotes, func(i, j int) bool {
		return jNotes[i].LastUpdated.After(jNotes[j].LastUpdated)
	})
	//Create header & footer
	header, footer := constructHeaderFooter()

	//Create memo area
	memo := tui.NewLabel("")
	memo.SetSizePolicy(tui.Expanding, tui.Expanding)

	//Create dates row
	var created = tui.NewLabel("")
	var lastUpdated = tui.NewLabel("")
	datesRow := constructDatesRow(created, lastUpdated)

	//Main part consist of memo + dates row
	mainPart := tui.NewVBox(memo, datesRow)
	mainPart.SetSizePolicy(tui.Expanding, tui.Expanding)
	mainPart.SetBorder(true)
	mainPart.SetTitle("Note")

	//Create list of notes to select from + make the list scrollable
	notesInfoList := constructNotesList(jNotes, created, lastUpdated, memo)
	scrollNotesInfo := tui.NewScrollArea(notesInfoList)
	scrollNotesInfo.SetSizePolicy(tui.Maximum, tui.Minimum)

	//Join everything together and run the UI
	root := tui.NewVBox(header, scrollNotesInfo, mainPart, footer)
	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	//Set theme
	initializeTheme(ui)

	//FIXME: DefaultFocusChain not working so this is a workaround
	formItems := make([]tui.Widget, 0)

	//Create the key bindings
	initializeKeyBinding(ui, scrollNotesInfo, notesInfoList, jNotes, formItems, root)

	return ui
}

func initializeTheme(ui tui.UI) {
	theme := tui.NewTheme()
	selected := tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack}
	theme.SetStyle("table.cell.selected", selected)
	theme.SetStyle("button.focused", selected)
	header := tui.Style{Bold: tui.DecorationOn, Bg: tui.ColorDefault, Fg: tui.ColorBlue}
	theme.SetStyle("label.header", header)
	footer := tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed}
	theme.SetStyle("label.footer", footer)
	ui.SetTheme(theme)

}

func initializeKeyBinding(ui tui.UI, scrollNotesInfo *tui.ScrollArea, notesInfo *tui.Table, jNotes []*jsonNote, formItems []tui.Widget, root *tui.Box) {
	//Esc
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	//Up arrow
	ui.SetKeybinding("Up", func() {
		windowSize := scrollNotesInfo.SizeHint().Y
		index := notesInfo.Selected()
		if index > 0 {
			notesInfo.Select(index - 1)
			if len(jNotes) <= windowSize {
				return
			}
			if index > windowSize {
				scrollNotesInfo.Scroll(0, -1)
			} else {
				scrollNotesInfo.ScrollToTop()
			}
		}
	})

	//Down arrow
	ui.SetKeybinding("Down", func() {
		windowSize := scrollNotesInfo.SizeHint().Y
		index := notesInfo.Selected()
		if index < len(jNotes)-1 {
			notesInfo.Select(index + 1)
			if len(jNotes) <= windowSize {
				return
			}
			if index < len(jNotes)-windowSize {
				scrollNotesInfo.Scroll(0, 1)
			} else {
				scrollNotesInfo.ScrollToBottom()
			}
		}
	})

	//Ctrl + D
	ui.SetKeybinding("Ctrl+D", func() {
		index := notesInfo.Selected()
		toBeDeleted := jNotes[index]
		notesInfo.RemoveRow(index)
		jNotes = append(jNotes[:index], jNotes[index+1:]...)
		notesInfo.Select(index)

		delete([]int64{toBeDeleted.ID})
	})

	//Ctrl + U
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
}

func constructHeaderFooter() (*tui.Table, *tui.Label) {
	header := tui.NewTable(0, 0)
	header.SetColumnStretch(0, 1)
	header.SetColumnStretch(1, 2)
	header.SetColumnStretch(2, 3)
	header.SetColumnStretch(3, 2)

	labelID := tui.NewLabel("ID")
	labelID.SetStyleName("header")
	labelNT := tui.NewLabel("Notebook Title")
	labelNT.SetStyleName("header")
	labelT := tui.NewLabel("Note Title")
	labelT.SetStyleName("header")
	labelTags := tui.NewLabel("Tags")
	labelTags.SetStyleName("header")
	header.AppendRow(labelID, labelNT, labelT, labelTags)
	//unselect it
	header.Select(-1)

	footer := tui.NewLabel("Press 'Ctrl+U' to update a note, 'Ctrl+D' to delete it or 'Esc' to exit.")
	footer.SetStyleName("footer")

	return header, footer
}

func constructNotesList(jNotes []*jsonNote, created, lastUpdated, memo *tui.Label) *tui.Table {
	notesInfoList := tui.NewTable(0, 0)
	notesInfoList.SetColumnStretch(0, 1)
	notesInfoList.SetColumnStretch(1, 2)
	notesInfoList.SetColumnStretch(2, 3)
	notesInfoList.SetColumnStretch(3, 2)
	notesInfoList.SetFocused(true)

	for _, note := range jNotes {
		notesInfoList.AppendRow(
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

	notesInfoList.OnSelectionChanged(func(t *tui.Table) {
		n := jNotes[t.Selected()]
		created.SetText(n.Created.Format("Jan 2 2006 15:04"))
		lastUpdated.SetText(n.LastUpdated.Format("Jan 2 2006 15:04"))
		memo.SetText(n.Memo)
	})
	notesInfoList.Select(0)

	return notesInfoList
}

func constructDatesRow(created, lastUpdated *tui.Label) *tui.Table {
	dates := tui.NewTable(0, 0)
	dates.AppendRow(tui.NewLabel("Created:"), created, tui.NewLabel("Last Updated:"), lastUpdated)
	dates.SetSizePolicy(tui.Expanding, tui.Minimum)
	//unselect it
	dates.Select(-1)
	return dates
}
