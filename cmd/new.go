/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taylorza/zxenv/internal/engine"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new project",
	Long:  `Creates a new project in the development environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return
		}
		env, err := engine.LoadEnvironment(".")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = engine.CreateProject(env, args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
