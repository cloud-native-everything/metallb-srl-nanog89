/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/cloud-native-everything/srklab/pkg"
	"github.com/spf13/cobra"
)


// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Add metrics to lab srklab",
	Long: `Add metrics to lab srklab, requires YAML config file with the specs.
ensure to clean the lab first.`,
	Run: func(cmd *cobra.Command, args []string) {

		srklab.Metrics(configFile)
	},
}

func init() {
	
	metricsCmd.Flags().StringVarP(&configFile,"conf", "c", "", "configuration YAML file")

	if err := metricsCmd.MarkFlagRequired("conf"); err != nil {
		fmt.Println(err)
	}
	
	rootCmd.AddCommand(metricsCmd)



	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metricsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metricsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
