package cmd

import (
	"fmt"
	"github.com/imtanmoy/authz/config"
	"github.com/imtanmoy/authz/logger"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(config.InitConfig)
	err := logger.InitLogger()
	if err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "Root",
	Short: "Cramstack Authz",
	Long:  "Cramstack Authz service for authorization",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
