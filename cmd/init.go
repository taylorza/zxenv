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
	"core3-128mb": {},
	"2gb":         {},
	"4gb":         {},
	"8gb":         {},
	"16gb":        {},
}

var emulators = map[string]struct{}{
	"cspect":  {},
	"cspect2": {},
	"zesarux": {},
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the development environment in the current folder",
	Long:  `Initializes the development environment in the current folder`,
	Run: func(cmd *cobra.Command, args []string) {
		size, _ := cmd.Flags().GetString("size")
		size = strings.ToLower(size)

		emulator, _ := cmd.Flags().GetString("emulator")
		emulator = strings.ToLower(emulator)

		if _, ok := imageUrls[size]; !ok {
			fmt.Println("Invalid size, valid sizes are: Core3-128mb, 2gb, 4gb, 8gb, 16gb")
			return
		}

		if _, ok := emulators[emulator]; !ok {
			fmt.Println("Invalid emulator, valid emulators are: cspect, zesarux")
			return
		}

		env := &engine.Environment{
			BasePath: ".",
			SDSize:   size,
			Emulator: emulator,
		}

		engine.SetupDevelopment(env)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("size", "s", "Core3-128mb", "Image size of the SD card (Core3-128mb, 2gb, 4gb, 8gb, 16gb)")
	initCmd.Flags().StringP("emulator", "e", "cspect", "Emulator to use (cspect, zesarux)")
}
