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
	"path/filepath"
	"strings"

	"github.com/cdiener/architeuthis/lib"
	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summarize k-mer assignments for classified taxa on taxonomic ranks.",
	Long: `Summarizes all individual k-mer assignments for each classified taxon
on taxonomic ranks. This allows you to see whether assignments are consistent on
higher ranks. For instance, even though a taxon might have discordant species
assignments those might all be within the same family or genus.`,
	Run: func(cmd *cobra.Command, args []string) {
		datadir, err := cmd.Flags().GetString("data-dir")
		if err != nil {
			log.Fatal(err)
		}
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatal(err)
		}

		version, ok := lib.HasTaxonkit()
		if !ok {
			log.Fatal("no taxonkit installation could be found :(")
		} else {
			log.Printf("Found taxonkit=%s.", version)
		}
		filetype, named := lib.GetFormat(args[0])
		if filetype != "kraken2" {
			log.Fatal("mapping summaries require a Kraken2 file")
		}
		if named {
			log.Println("detected Kraken2 output with taxon names.")
		}

		id := strings.Split(filepath.Base(args[0]), ".")[0]
		kmap, err := lib.SummarizeKmers(args[0], named)
		if err != nil {
			log.Fatal("Failed to build the kmer mapping hash.")
		}
		collapsed := lib.CollapseRanks(kmap, datadir, format)

		out, _ := cmd.Flags().GetString("out")
		log.Printf("Saving map to %s.", out)
		lib.SaveMapping(collapsed, out, id)
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
	summaryCmd.Flags().String("data-dir", "", "The path to the taxonomy dumps.")
	summaryCmd.Flags().StringP("format", "f", "{K};{p};{c};{o};{f};{g};{s}", "The taxonomic ranks to connsider during scoring.")

}
