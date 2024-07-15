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

// scoreCmd represents the score command
var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Scores and evaluates reads.",
	Long: `Not all read classifications are the same and this comand will help you to
quantify this by applyting several advanced scoring metrics on individual read kmer
patterns and saving them to a CSWV file. In general it uses the following two scoring schemes:

Consistency
-----------
What proportion of individual kmer classifications are within the taxonomy of the
final read classification. Kmers can be classified into various taxonomic ranks
but if all of those appear in the read classification we would call this 100% consistent.

Ambiguity
---------
What is the variability of kmer classification on all taxonomic ranks. For instance, a
read may map to 100 different species but only one genus. in that case the species
assignment would be very ambiguous but the genus assignment would not be. We report
the number of different classifications (multiplicity) and the shannon index (taking
abundance of kmers into account as well).
`,
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
			log.Fatal("read scoring requires a Kraken2 file.")
		}
		if named {
			log.Println("detected Kraken2 output with taxon names.")
		}

		out, _ := cmd.Flags().GetString("out")

		err = lib.ScoreReadsToFile(args[0], out, datadir, format, named)
		if err != nil {
			log.Fatalf("Saving file failed with error: %v", err)
		}
	},
}

func init() {
	mappingCmd.AddCommand(scoreCmd)

	scoreCmd.Flags().String("out", "mapping_scores.csv", "The output file (CSV format).")
	scoreCmd.Flags().String("data-dir", "", "The path to the taxonomy dumps.")
	scoreCmd.Flags().StringP("format", "f", "{k};{p};{c};{o};{f};{g};{s}", "The taxonomic ranks to connsider during scoring.")
}
