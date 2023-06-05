/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/cloud-native-everything/srklab/pkg"
	"github.com/spf13/cobra"
)


// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear srklab",
	Long: `Clear srklab, requires YAML config file with the specs.
sorry to see you leave.`,
	Run: func(cmd *cobra.Command, args []string) {

		srklab.Clear(configFile)
	},
}

func init() {
	
	clearCmd.Flags().StringVarP(&configFile,"conf", "c", "", "configuration YAML file")

	if err := clearCmd.MarkFlagRequired("conf"); err != nil {
		fmt.Println(err)
	}
	
	rootCmd.AddCommand(clearCmd)



	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
