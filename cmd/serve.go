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
		casbin.Init(db.DB)

		// _, _ = casbin.Enforcer.AddPolicy("ROLE_1", "perm_1", "read")
		// _, _ = casbin.Enforcer.AddPolicy("ROLE_1", "perm_2", "read")

		// _, _ = casbin.Enforcer.AddPolicy("tanmoy", "perm_2", "view")

		//To add rule to an user or vice-versa
		// _, _ = casbin.Enforcer.AddGroupingPolicy("alice", "data_admin")

		// _, _ = casbin.Enforcer.AddGroupingPolicy("bob", "ROLE_2")

		// res, _ := casbin.Enforcer.Enforce("alice", "perm_1", "read")
		// fmt.Println(res)
		// res, _ = casbin.Enforcer.Enforce("bob", "perm_2", "view")
		// res, _ = casbin.Enforcer.Enforce("tanmoy", "perm_2", "view")
		// res, _ = casbin.Enforcer.Enforce("tanmoy", "perm_2", "read")

		// Get all roles
		// allRoles := casbin.Enforcer.GetAllRoles()
		// fmt.Println(allRoles)

		// Get all permission may be
		// allRoles := casbin.Enforcer.GetAllObjects()
		// fmt.Println(allRoles)

		// Get all subjects of policy
		// allRoles := casbin.Enforcer.GetAllSubjects()
		// fmt.Println(allRoles)

		// Get all roles for a user
		// allRoles, _ := casbin.Enforcer.GetRolesForUser("alice")
		// fmt.Println(allRoles)

		// Get all permission may be
		allRoles, _ := casbin.Enforcer.DeleteUser("alice")
		fmt.Println(allRoles)

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
