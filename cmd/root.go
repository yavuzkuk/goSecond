/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go2Second/request"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go2Second",
	Short: "A brief description of your application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// request.GetProduct(product, min, max, output, descFilter, show)
		request.PageNumberDolap(product, min, max, output, descFilter, show)
		request.PageSahibinden(product, min, max, descFilter, show, output)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var product string
var min int
var max int
var output string
var descFilter string
var show bool

func init() {

	rootCmd.Flags().BoolVarP(&show, "show", "s", false, "Disable output")
	rootCmd.Flags().IntVarP(&min, "min", "", -1, "Minimum price")
	rootCmd.Flags().IntVarP(&max, "max", "", -1, "Maximum price")
	rootCmd.Flags().StringVarP(&product, "product", "p", "", "Product name")
	rootCmd.Flags().StringVarP(&descFilter, "desc", "d", "", "Description filter")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output name")
}
