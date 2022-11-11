package linkshortener

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

const (
	webPort = "8000"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the Link Shortener",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd, args)
	},
}

func serve(cmd *cobra.Command, args []string) {
	http.HandleFunc("/", getLink)
	http.HandleFunc("/create", createLink)

	_ = http.ListenAndServe(fmt.Sprintf(":%s", webPort), nil)
}
