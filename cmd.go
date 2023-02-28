package main

import (
	"fmt"
	// "os"

	"github.com/spf13/cobra"
)

var (
	InputUp     = "up"
	InputInit   = "init"
	InputDown   = "down"
	InputRemove = "remove"
)

var rootCmd = &cobra.Command{
	Use:   "sdm",
	Short: "Software Deployment Manager",
	// Long: `Software Deployment Manager. This software help you manger muplite container templates.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SDM",
	// Long:  `All software has versions. This is sdm's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SDM Software Deployment Manager v0.0.1 -- HEAD")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialization configuration.",
	// Long:  `Initialization configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Init !!!!")
		initConfigFile()
		initDockerComposefile()
		initSettingsfile()
		cmdInput = InputInit
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers.",
	// Long:  `Create and start containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdInput = InputUp
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers, networks.",
	// Long:  `Create and start containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdInput = InputDown
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove all applications.",
	// Long:  `Create and start containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdInput = InputRemove
	},
}
