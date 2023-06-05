/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/cloud-native-everything/srklab/pkg"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start srklab",
	Long: `Start srklab, requires YAML config file with the specs.
ensure to clean the lab first.`,
	Run: func(cmd *cobra.Command, args []string) {

		srklab.Start(configFile)
	},
}

func init() {
	
	startCmd.Flags().StringVarP(&configFile,"conf", "c", "", "configuration YAML file")

	if err := startCmd.MarkFlagRequired("conf"); err != nil {
		fmt.Println(err)
	}
	
	rootCmd.AddCommand(startCmd)



	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
