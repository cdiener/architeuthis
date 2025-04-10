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
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cdiener/architeuthis/lib"
	"github.com/spf13/cobra"
)

var TaxidIndex = map[string]string{
	"bracken":        "taxonomy_id",
	"bracken-merged": "taxonomy_id",
	"mapping":        "classification",
}

// lineageCmd represents the lineage command
var lineageCmd = &cobra.Command{
	Use:   "lineage",
	Short: "Add lineage information to Bracken output.",
	Long: `Sometimes you would like to annotate taxonomy IDs with their full
canonical lineage. This command helps with this.

The 'lineage' command does require taxonkit and the databases to be installed to
work.
`,
	Run: func(cmd *cobra.Command, args []string) {
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatal(err)
		}
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatal(err)
		}
		datadir, err := cmd.Flags().GetString("data-dir")
		if err != nil {
			log.Fatal(err)
		}
		filetype, lineage := lib.GetFormat(args[0])
		if lineage {
			log.Fatalf("file %s already contains lineage information", args[0])
		}
		if filetype != "bracken" && filetype != "mapping" && filetype != "bracken-merged" {
			log.Fatalf("file %s is not bracken or mapping summary format", args[0])
		}
		err = FoldInLineage(args[0], filetype, format, out, datadir)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lineageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lineageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	lineageCmd.Flags().String("data-dir", "", "The path to the taxonomy dumps.")
	lineageCmd.Flags().StringP("format", "f", "{K};{p};{c};{o};{f};{g};{s}", "The taxonomic ranks to connsider during scoring.")
	lineageCmd.Flags().StringP("out", "o", "annotated.csv", "The filename of the output CSV.")
}

func FoldInLineage(filename string, filetype string, format string, out string, data_dir string) error {
	version, ok := lib.HasTaxonkit()
	if !ok {
		return errors.New("no taxonkit installation could be found :(")
	} else {
		log.Printf("Found taxonkit=%s.", version)
	}

	log.Printf("Mapping taxonomy IDs from %s.", filename)

	infile, err := os.Open(filename)
	if err != nil {
		return err
	}

	reader := csv.NewReader(infile)
	if filetype == "bracken" {
		reader.Comma = '\t'
	}
	header, err := reader.Read()
	if err != nil {
		return err
	}

	idx := -1
	for i, v := range header {
		if v == TaxidIndex[filetype] {
			idx = i
			break
		}
	}
	if idx < 0 {
		return errors.New(
			"this file does not seem to be of the specified type." +
				"Are you sure the your '--type' is the correct one?")
	}

	taxids := make(map[string]bool, 100)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		taxids[record[idx]] = true
	}

	log.Printf("Will map %d unique taxids with taxonkit.", len(taxids))
	lineages := lib.AddLineage(taxids, data_dir, format)

	log.Printf("Writing annotated data to %s.", out)

	infile.Seek(0, io.SeekStart)
	outfile, err := os.Create(out)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(outfile)
	header, err = reader.Read()
	if err != nil {
		return err
	}
	header = append(header, "lineage", "taxid_lineage")
	err = writer.Write(header)
	if err != nil {
		return err
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		l := *lineages[record[idx]]
		record = append(record, strings.Join(l.Names, ";"), strings.Join(l.Taxids, ";"))
		writer.Write(record)
	}
	writer.Flush()

	return nil
}
