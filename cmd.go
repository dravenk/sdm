package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "sdm",
	Short: "Software Deployment Manager",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SDM",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("SDM Software Deployment Manager v0.0.1 -- HEAD")
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete configuration files.",
	Run: func(cmd *cobra.Command, args []string) {
		logln("Execute: Remove ", tplSettings)
		if err := os.Remove(tplSettings); err != nil {
			log.Fatal(err)
		}
		logln("Execute: Remove ", tplCompose)
		if err := os.Remove(tplCompose); err != nil {
			log.Fatal(err)
		}
		logln("Execute: Remove ", valuesFile)
		if err := os.Remove(valuesFile); err != nil {
			log.Fatal(err)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialization configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		cmdInput = InputInit

		initConfigFile()
		initDockerComposefile()
		initSettingsfile()

		minport, _ := cmd.Flags().GetInt("minport")
		if minport != 0 {
			Conf.Minport = minport
		}

		maxport, _ := cmd.Flags().GetInt("maxport")
		if maxport != 0 {
			Conf.Maxport = maxport
		}

		apps, _ := cmd.Flags().GetStringSlice("app")
		if len(apps) > 0 {
			Conf.AppsName = apps
		}

		dir, _ := cmd.Flags().GetString("dir")
		if dir != "" {
			Conf.Appsdir = dir
		}
		image, _ := cmd.Flags().GetString("image")
		if image != "" {
			Conf.Image = image
		}

		initConfigFile()
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers.",
	Run: func(cmd *cobra.Command, args []string) {
		cmdInput = InputUp
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers, networks.",
	Run: func(cmd *cobra.Command, args []string) {
		cmdInput = InputDown
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove all applications.",
	Run: func(cmd *cobra.Command, args []string) {
		cmdInput = InputRemove
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	downCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")
	initCmd.PersistentFlags().StringSliceP("app", "a", []string{"dp1", "dp2"}, "All the applications name.")
	initCmd.PersistentFlags().String("dir", "apps", "The directory where you store your projects.")
	initCmd.PersistentFlags().String("image", "dravenk/dp:10-fpm", "The directory where you store your projects.")
	initCmd.PersistentFlags().Int("minport", 8000, "The minimum port can be use.")
	initCmd.PersistentFlags().Int("maxport", 9000, "The maximum port can be use.")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(cleanCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// // Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// // Search config in home directory with name ".config" (without extension).
		// viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
