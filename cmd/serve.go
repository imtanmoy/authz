package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/imtanmoy/authz/casbin"
	"github.com/imtanmoy/authz/db"
	"github.com/imtanmoy/authz/logger"
	"github.com/imtanmoy/authz/server"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Run: func(cmd *cobra.Command, args []string) {
		err := db.InitDB()
		if err != nil {
			logger.Fatalf("%s : %s", "Database Could not be initiated", err)
		}
		logger.Info("Database Initiated...")
		casbin.InitCasbin(db.DB)
		added, err := casbin.Enforcer.AddPolicy("eve", "data3", "read")
		if err != nil {
			panic(err)
		}
		fmt.Println(added)
		allRoles := casbin.Enforcer.GetAllRoles()
		fmt.Println(allRoles)
		policy := casbin.Enforcer.GetPolicy()
		fmt.Println(policy)
		removed, err := casbin.Enforcer.RemovePolicy("eve", "data3", "read")
		if err != nil {
			panic(err)
		}
		fmt.Println(removed)
		server, err := server.NewServer()
		if err != nil {
			logger.Fatalf("%s : %s", "Server could not be started", err)
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			oscall := <-c
			logger.Infof("system call:%+v", oscall)
			cancel()
		}()

		if err := server.Start(ctx); err != nil {
			logger.Infof("failed to serve:+%v\n", err)
		}
		close(c)
	},
}
