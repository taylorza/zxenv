/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taylorza/zxenv/internal/engine"
)

var imageUrls = map[string]struct{}{
	"32mb":  {},
	"128mb": {},
	"512mb": {},
	"2gb":   {},
	"4gb":   {},
	"8gb":   {},
	"16gb":  {},
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the development environment in the current folder",
	Long:  `Initializes the development environment in the current folder`,
	Run: func(cmd *cobra.Command, args []string) {
		size, _ := cmd.Flags().GetString("size")
		size = strings.ToLower(size)

		if _, ok := imageUrls[size]; !ok {
			fmt.Println("Invalid size, valid sizes are: 32mb, 128mb, 512mb, 2gb, 4gb, 8gb, 16gb")
			return
		}

		env := &engine.Environment{
			BasePath: ".",
			SDSize:   size,
		}

		engine.SetupDevelopment(env)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("size", "s", "512mb", "Image size of the SD card (32mb, 128mb, 512mb, 2gb, 4gb, 8gb, 16gb)")
}
