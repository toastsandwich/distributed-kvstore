package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/toastsandwich/kvstore/internal/config"
	"github.com/toastsandwich/kvstore/internal/server"
)

var configPath string

var ServerCmd = &cobra.Command{
	Use:     "start",
	Short:   "start server with provided config",
	Long:    "use start command to start a server, it requires a configuration which must be provided by -C flag",
	Example: "kvstore start -C ~/kvstore/config.yaml",

	RunE: runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
	if configPath == "" {
		return fmt.Errorf("config path cannot be empty")
	}

	c, err := config.ReadFrom(configPath)
	if err != nil {
		return err
	}

	srv := server.New(c)
	srv.Start()
	return nil
}

func init() {
	ServerCmd.Flags().StringVarP(&configPath, "config", "c", "", "used to privide configuration for kvstore")
}
