package cmd

import "github.com/spf13/cobra"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Initiate rest API interface",
	Long: "Run a http server for managing notes/notebooks via REST calls\n" +
		"If no -p flag is not set the default port will be 8080\n" +
		"Available endpoints:\n" +
		"POST /addNote \n" +
		"PUT /updateNote \n" +
		"GET /getNotesByID/{ids} (comma separated IDs) \n" +
		"GET /getNotesByNotebookTitle/{notebookTitles} (comma separated notebook titles) \n" +
		"GET /getNotesByTags/{Tags} (comma separated tags) \n" +
		"GET /getAllNotes \n" +
		"DELETE /deleteNotes/{ids} (comma separated IDs)\n" +
		"GET /searchBy/{keyword} \n" +
		"PUT /updateNotebook/{oldTitle}/{newTitle} \n" +
		"DELETE /deleteNotebooks/{notebookTitles} (comma separated notebook titles)\n",
	Example: "serve -p 7000",
	Run:     serve,
}

func serve(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")
	server := NewServer()
	server.Run(port)
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "Server port")
}
