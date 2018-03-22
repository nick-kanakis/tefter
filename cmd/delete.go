package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var deleteNoteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete one or more notes based on ID(s)",
	Args:    cobra.MinimumNArgs(1),
	Example: "delete 1,2,...",
	Run:     deleteWrapper,
}

func init() {
	rootCmd.AddCommand(deleteNoteCmd)
}

func deleteWrapper(cmd *cobra.Command, args []string) {
	var ids = make([]int64, 0, len(args))
	for _, argument := range args {
		id, err := strconv.ParseInt(argument, 10, 64)
		if err != nil {
			log.Panicf("Could note transform input to id for argument %v", argument)
		}
		ids = append(ids, id)
	}
	if err := delete(ids); err != nil {
		log.Panicln(err)
	}
}

func delete(ids []int64) error {
	err := NoteDB.DeleteNotes(ids)
	if err != nil {
		return fmt.Errorf("Error while deleting notes, error msg: %v", err)
	}
	return nil
}
