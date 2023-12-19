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

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter Kraken output based on read quality.",
	Long: `Not all read classifications are the same and this comand will help you to
filter reads by several advanced scoring metrics on individual read kmer
patterns. In general it uses the following two scoring schemes:

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
abundance of kmers into account as well).`,
	Run: func(cmd *cobra.Command, args []string) {
		datadir, err := cmd.Flags().GetString("data-dir")
		if err != nil {
			log.Fatal(err)
		}
		max_entropy, err := cmd.Flags().GetFloat64("max-entropy")
		if err != nil {
			log.Fatal(err)
		}
		min_consistency, err := cmd.Flags().GetFloat64("min-consistency")
		if err != nil {
			log.Fatal(err)
		}
		max_multiplicity, err := cmd.Flags().GetUint32("max-multiplicity")
		if err != nil {
			log.Fatal(err)
		}

		version, ok := lib.HasTaxonkit()
		if !ok {
			log.Fatal("no taxonkit installation could be found :(")
		} else {
			log.Printf("Found taxonkit=%s.", version)
		}
		filetype, _ := lib.GetFormat(args[0])
		if filetype != "kraken2" {
			log.Fatal("read scoring requires a Kraken2 file.")
		}

		out, _ := cmd.Flags().GetString("out")

		err = lib.FilterReads(args[0], out, datadir, "{k};{p};{c};{o};{f};{g};{s}",
			min_consistency, max_entropy, max_multiplicity)

		if err != nil {
			log.Fatalf("filtering failed with error: %v.", err)
		}
	},
}

func init() {
	mappingCmd.AddCommand(filterCmd)

	filterCmd.Flags().String("data-dir", "", "The path to the taxonomy dumps.")
	filterCmd.Flags().String("out", "filtered.k2", "The output file (Kraken format).")
	filterCmd.Flags().Float64("max-entropy", 0.1, "Maximum entropy for kmer classifications at classified rank.")
	filterCmd.Flags().Float64("min-consistency", 0.9, "Minimum consistency of the read classification.")
	filterCmd.Flags().Uint32("max-multiplicity", 2, "Maximum number of alternative classifications on the classified rank.")
}
