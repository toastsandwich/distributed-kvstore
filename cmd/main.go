package cmd

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var App = cobra.Command{
	Use:   "kvstore",
	Short: "kvstore server",
	Long:  "kvstore is a distribute server",
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func Main() {
	App.AddCommand(&ServerCmd)
	if err := App.Execute(); err != nil {
		os.Exit(1)
	}
}
