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
	"golang.org/x/exp/slices"
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

Supported formats are Bracken outout and mapping summaries.`,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatal("Error in reading the output filename.")
		}
		if len(args) < 1 {
			log.Fatal("Need at least 2 files to merge.")
		}
		formats := make([]string, len(args))
		for i, fname := range args {
			f, _ := lib.GetFormat(fname)
			if f == "" {
				log.Fatalf("file %s is not recognized as a supported type", fname)
			}
			formats[i] = f
		}
		formats = slices.Compact(formats)
		if len(formats) > 1 {
			log.Fatalf("arguments have differing formats, found the following: %v", formats)
		}

		if formats[0] == "kraken2" || formats[0] == "report" {
			log.Fatalf("merging kraken2 output or report files is not supported")
		}

		if formats[0] == "mapping" || formats[0] == "bracken-merged" {
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
	mergeCmd.Flags().StringP("out", "o", "merged.csv", "The output filename.")
}