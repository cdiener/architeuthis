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

// kmersCmd represents the kmers command
var kmersCmd = &cobra.Command{
	Use:   "kmers",
	Short: "Summarize k-mer assignments for classified taxa.",
	Long: `Summarizes all individual k-mer assignments for each classified taxon
across reads. That is particularly helpful to check how unique you assignments are or
to identify instances where one taxon can also be classified as another taxon.`,
	Run: func(cmd *cobra.Command, args []string) {
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
			log.Fatal("Failed to build the mapping hash.")
		}

		out, _ := cmd.Flags().GetString("out")
		log.Printf("Saving map to %s.", out)
		lib.SaveMapping(kmap, out, id)
	},
}

func init() {
	mappingCmd.AddCommand(kmersCmd)

	kmersCmd.Flags().String("out", "mapping_kmers.csv", "The output file (CSV format).")
}
