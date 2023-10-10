/*
Copyright © 2023 Christian Diener <mail(a)cdiener.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"log"

	"github.com/cdiener/architeuthis/lib"
	"github.com/spf13/cobra"
)

var SimpleFormats = map[string]bool{
	"bracken": false,
	"mapping": true,
	"report":  false,
}

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge various output files related to Kraken.",
	Long: `This quickly merges Kraken output files across several samples.

This can optionally add in the full lineage if desired with the '--with-lineage'
option. However this will require a 'taxonkit' installation.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatal("Error reading the format argument.")
		}
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatal("Error in reading the output filename.")
		}
		is_simple, ok := SimpleFormats[format]
		if !ok {
			log.Fatalf(
				"%s is not a valid format. "+
					"Please choose from 'bracken', 'mapping, or 'report'.",
				format)
		}
		if len(args) < 1 {
			log.Fatal("Need at least 2 files to merge.")
		}
		if is_simple {
			err = lib.SimpleAppend(args, out)
		} else {
			err = lib.SampleAppend(args, out)
		}
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	mergeCmd.Flags().StringP("format", "f", "bracken", "The format of the files to merge.")
	mergeCmd.Flags().StringP("out", "o", "merged.csv", "The output filename.")
}
