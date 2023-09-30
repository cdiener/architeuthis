/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/cdiener/architeuthis/lib"
	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summarize k-mer assignments for classified taxa.",
	Long: `Summarizes all individual k-mer assignments for each classified taxon
across reads. That is particularly helpful to check how unique you assignments are or
to identify isntances where one taxon can also be classified to another taxon.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := strings.Split(filepath.Base(args[0]), ".")[0]
		kmap, err := lib.Summarize(args[0])
		if err != nil {
			log.Fatal("Failed to build the mapping hash.")
		}

		out, _ := cmd.Flags().GetString("out")
		log.Printf("Saving map to %s.", out)
		lib.SaveMapping(kmap, out, id)
	},
}

func init() {
	mappingCmd.AddCommand(summaryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summaryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// summaryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	summaryCmd.Flags().String("out", "mapping_summary.csv", "The output file (CSV format).")
}
