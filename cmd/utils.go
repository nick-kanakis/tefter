package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nicolasmanic/tefter/model"
)

//TODO: unit test
func viEditor(text string) string {
	vi := "vim"
	fpath := os.TempDir() + "/tmpMemo.txt"
	f, err := os.Create(fpath)
	defer os.Remove(fpath)
	if err != nil {
		log.Fatalf("Could not create tmp file for memo, error msg %v", err)
	}
	_, err = io.Copy(f, strings.NewReader(text))
	if err != nil {
		fmt.Printf("Failed copying memo to tmp file, error msg: %v\n", err)
	}
	f.Close()
	path, err := exec.LookPath(vi)
	if err != nil {
		log.Fatalf("Could not open VI, error msg %v", err)
	}

	cmd := exec.Command(path, fpath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Could not start VI, error msg %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Error while using editing memo, error msg %v", err)
	}

	memo, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatalf("Could not read tmp file, error msg %v", err)
	}
	return string(memo)
}

func int64Slice(input []int) []int64 {
	var result = make([]int64, 0, len(input))
	for _, tmp := range input {
		result = append(result, int2int64(tmp))
	}
	return result
}

func int2int64(input int) int64 {
	return int64(input)
}

func collectNotesFromDB(ids []int, notebookTitles, tags []string, getAll bool) map[int64]*model.Note {
	var notesMap = make(map[int64]*model.Note, 0)
	if getAll {
		//Get all notes in the DB
		allNotes, err := NoteDB.GetNotes([]int64{})
		if err != nil {
			log.Panicf("Error while retrieving notes by id , error msg: %v", err)
		}
		for _, note := range allNotes {
			notesMap[note.ID] = note
		}
		return notesMap
	}
	if len(ids) > 0 {
		//Add notes based on ids
		idNotes, err := NoteDB.GetNotes(int64Slice(ids))
		if err != nil {
			log.Panicf("Error while retrieving notes by id, error msg: %v", err)
		}
		for _, note := range idNotes {
			notesMap[note.ID] = note
		}
	}
	if len(notebookTitles) > 0 {
		//Get notes based on notebook titles
		for _, notebookTitle := range notebookTitles {
			notebook, err := NotebookDB.GetNotebookByTitle(notebookTitle)
			if err != nil {
				log.Panicf("Error while retrieving notebook for title: %v, error msg %v", notebookTitle, err)
			} else if notebook != nil {
				for _, note := range notebook.Notes {
					notesMap[note.ID] = note
				}
			}
		}
	}
	if len(tags) > 0 {
		//Get notes based on tags
		tagNotes, err := NoteDB.GetNotesByTag(tags)
		if err != nil {
			log.Panicf("Error while retrieving notes by tag, error msg: %v", err)
		}
		for _, note := range tagNotes {
			notesMap[note.ID] = note
		}
	}
	return notesMap
}

func noteMap2Slice(m map[int64]*model.Note) []*model.Note {
	notes := make([]*model.Note, 0, len(m))
	for _, note := range m {
		notes = append(notes, note)
	}
	return notes
}

func tagMap2Slice(m map[string]bool) []string {
	tags := make([]string, 0, len(m))
	for tag := range m {
		tags = append(tags, tag)
	}
	return tags
}

//TODO: Unit test
func transformNotes2JSONNotes(notes []*model.Note) ([]*jsonNote, error) {
	var jNotes []*jsonNote
	notebookTitlesMap, err := NotebookDB.GetAllNotebooksTitle()
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving Notebooks titles, error msg: %v", err)
	}

	for _, note := range notes {
		jNote := &jsonNote{
			ID:            note.ID,
			Title:         note.Title,
			Memo:          note.Memo,
			Created:       note.Created,
			LastUpdated:   note.LastUpdated,
			Tags:          tagMap2Slice(note.Tags),
			NotebookTitle: notebookTitlesMap[note.NotebookID],
		}
		jNotes = append(jNotes, jNote)
	}
	return jNotes, nil
}
