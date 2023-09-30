/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// mappingCmd represents the mapping command
var mappingCmd = &cobra.Command{
	Use:   "mapping",
	Short: "Analyze read and k-mer level mapping.",
	Long: `The mapping command helps with analyzing the read- and k-mer level
taxonomic assignments of your Kraken output.`,
}

func init() {
	rootCmd.AddCommand(mappingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mappingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mappingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
