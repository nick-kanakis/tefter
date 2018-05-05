package cmd

import (
	"github.com/rivo/tview"
	"testing"
	"time"
)

func TestCreateUI(t *testing.T) {
	notes := mockJSONNotes()
	emptyNotes := []*jsonNote{}

	if createUI(notes) == nil {
		t.Error("UI should not be nil")
	}

	if createUI(emptyNotes) != nil {
		t.Error("Empty notes should return a nil ui")
	}

	if createUI(nil) != nil {
		t.Error("Nil input should return a nil ui")
	}
}

func TestConstructUpdateForm(t *testing.T) {

	jNote := &jsonNote{
		ID:            1,
		Title:         "Test Title",
		Tags:          []string{"tag1", "tag2"},
		NotebookTitle: "notebook",
	}

	form := constructUpdateForm(jNote)
	notebookTitle := form.GetFormItemByLabel("Notebook Title:").(*tview.InputField)
	noteTitle := form.GetFormItemByLabel("Note Title:").(*tview.InputField)
	tags := form.GetFormItemByLabel("Tags:").(*tview.InputField)

	if notebookTitle.GetText() != "notebook" {
		t.Errorf("Wrong notebook title: expected :%q got %q ", "notebook", notebookTitle.GetText())
	}

	if noteTitle.GetText() != "Test Title" {
		t.Errorf("Wrong note title: expected :%q got %q ", "Test Title", noteTitle.GetText())
	}

	if tags.GetText() != "tag1,tag2" {
		t.Errorf("Wrong tags: expected :%q got %q ", "tag1,tag2", tags.GetText())
	}
}

func TestConstructDatesRow(t *testing.T) {
	jNotes := []*jsonNote{
		&jsonNote{
			ID:            1,
			Title:         "Test Title",
			Tags:          []string{"tag1", "tag2"},
			NotebookTitle: "notebook",
		},
		&jsonNote{
			ID:            2,
			Title:         "Test Title2",
			Tags:          []string{},
			NotebookTitle: "notebook",
		},
	}
	datesRow := constructDatesRow(jNotes)

	if datesRow.GetColumnCount() != 4 {
		t.Errorf("Expected 4 columns in dates row got: %v", datesRow.GetColumnCount())
	}

	createdTxt := datesRow.GetCell(0, 1).Text
	updatedTxt := datesRow.GetCell(0, 3).Text

	if _, err := time.Parse("Jan 2 2006 15:04", createdTxt); err != nil {
		t.Errorf("Failed to convert created time")
	}
	if _, err := time.Parse("Jan 2 2006 15:04", updatedTxt); err != nil {
		t.Errorf("Failed to convert updated time")
	}
}

func TestConstructNotesTable(t *testing.T) {
	jNotes := []*jsonNote{
		&jsonNote{
			ID:            1,
			Title:         "Test Title",
			Tags:          []string{"tag1", "tag2"},
			NotebookTitle: "notebook",
		},
		&jsonNote{
			ID:            2,
			Title:         "Test Title2",
			Tags:          []string{},
			NotebookTitle: "notebook",
		},
	}

	notesTable := constructNotesTable(jNotes)

	if notesTable.GetRowCount() != 3 {
		t.Errorf("Expected 3 rows but got: %v", notesTable.GetRowCount())
	}

	if notesTable.GetCell(1, 0).Text != "1" {
		t.Errorf("Wrong note id: expected : 1 got %q ", notesTable.GetCell(1, 0).Text)
	}

	if notesTable.GetCell(1, 1).Text != "notebook" {
		t.Errorf("Wrong notebook title: expected :%q got %q ", "notebook", notesTable.GetCell(1, 1).Text)
	}

	if notesTable.GetCell(1, 2).Text != "Test Title" {
		t.Errorf("Wrong note title: expected :%q got %q ", "Test Title", notesTable.GetCell(1, 2).Text)
	}

	if notesTable.GetCell(1, 3).Text != "tag1,tag2" {
		t.Errorf("Wrong tags: expected :%q got %q ", "tag1,tag2", notesTable.GetCell(1, 3))
	}
}
