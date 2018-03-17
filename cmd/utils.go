package cmd

import (
	"log"
	"fmt"
	"io"
	"io/ioutil"
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

func int2int64(input []int) []int64 {
	var result = make([]int64, 0, len(input))
	for _, tmp := range input {
		result = append(result, int64(tmp))
	}
	return result
}

func collectNotes(ids []int, notebookTitles, tags []string, printAll bool) map[int64]*model.Note {
	var notesMap = make(map[int64]*model.Note, 0)
	if printAll {
		allNotes, err := NoteDB.GetNotes([]int64{})
		if err != nil {
			fmt.Printf("Error while retrieving notes by id , error msg: %v", err)
			return nil
		} 
		for _, note := range allNotes {
			notesMap[note.ID] = note
		}
		return notesMap
	}
	//Add notes based on id
	if len(ids) > 0 {
		idNotes, err := NoteDB.GetNotes(int2int64(ids))
		if err != nil {
			fmt.Printf("Error while retrieving notes by id, error msg: %v", err)
			return nil
		}
		for _, note := range idNotes {
			notesMap[note.ID] = note
		}
		
	}
	if len(notebookTitles) > 0 {
		//Add notes based on notebook
		for _, notebookTitle := range notebookTitles {
			notebook, err := NotebookDB.GetNotebookByTitle(notebookTitle)
			if err != nil {
				fmt.Printf("Error while retrieving notebook for title: %v, error msg %v", notebookTitle, err)
				return nil
			} else if notebook != nil {
				for _, note := range notebook.Notes {
					notesMap[note.ID] = note
				}
			}
		}
	}
	if len(tags) > 0 {
		tagNotes, err := NoteDB.GetNotesByTag(tags)
		if err != nil {
			fmt.Print("Error while retrieving notes by tag, error msg: %v", err)
			return nil
		}
		for _, note := range tagNotes {
			notesMap[note.ID] = note
		}
	}
	return notesMap
}

func noteMap2Slice(m map[int64]*model.Note) []*model.Note {
	notes := make([]*model.Note,0, len(m))
	for _, note := range m {
		notes = append(notes, note)
	}
	return notes
}

func tagMap2Slice(m map[string]bool) []string {
	tags := make([]string,0,len(m))
	for tag := range m {
		tags = append(tags, tag)
	}
	return tags
}
