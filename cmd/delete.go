package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deleteNoteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete one or more notes based on ID",
	Example: "delete -i 1,2,...",
	Run: func(cmd *cobra.Command, args []string) {

		ids, _ := cmd.Flags().GetIntSlice("ids")
		ids64 := int2int64(ids)
		err := NoteDB.DeleteNotes(ids64)

		if err != nil {
			//TODO handle the error
			fmt.Printf("error msg: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteNoteCmd)
	deleteNoteCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma-separated note ids")
	deleteNoteCmd.MarkFlagRequired("ids")
}
