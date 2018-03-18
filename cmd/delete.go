package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var deleteNoteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete one or more notes based on ID",
	Example: "delete -i 1,2,...",
	Run:     deleteWrapper,
}

func init() {
	rootCmd.AddCommand(deleteNoteCmd)
	deleteNoteCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma-separated note ids")
	deleteNoteCmd.MarkFlagRequired("ids")
}

func deleteWrapper(cmd *cobra.Command, args []string) {
	ids, _ := cmd.Flags().GetIntSlice("ids")
	delete(ids, args)
}

func delete(ids []int, args []string) {
	ids64 := int2int64(ids)
	err := NoteDB.DeleteNotes(ids64)

	if err != nil {
		log.Panicf("Error while deleting notes, error msg: %v", err)
	}
}
