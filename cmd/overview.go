package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var overviewCmd = &cobra.Command{
	Use:     "overview",
	Short:   "Take a quick glance at the available notebooks and notes",
	Example: "overview -d",
	Run:     overviewWrapper,
}

func init() {
	rootCmd.AddCommand(overviewCmd)
	overviewCmd.Flags().BoolP("deep", "d", false, "Deep overview of notebooks")
}

func overviewWrapper(cmd *cobra.Command, args []string) {
	deep, _ := cmd.Flags().GetBool("deep")
	notebooks, err := NotebookDB.GetNotebooks([]int64{})
	if err != nil {
		log.Panicln(err)
	}
	printOverview(notebooks, deep)
}

func printOverview(notebooks []*model.Notebook, deep bool) {
	fmt.Println("> Notebooks:")
	for _, notebook := range notebooks {
		fmt.Println(" > " + notebook.Title)
		if deep {
			fmt.Println("  > notes:")
			for _, note := range notebook.Notes {
				fmt.Println("  - " + strconv.FormatInt(note.ID, 10) + " " + note.Title)
			}
		} else {
			fmt.Println("  > number of notes: " + strconv.Itoa(len(notebook.Notes)))
		}
	}
}
