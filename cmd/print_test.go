package cmd

import "testing"

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
