package lib

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strings"
)

var bracken_header = []string{"name", "taxonomy_id", "taxonomy_lvl",
	"kraken_assigned_reads", "added_reads", "new_est_reads",
	"fraction_total_reads"}

var mappings_header = []string{"sample_id", "classification",
	"n_reads", "taxid", "n_kmers"}

func GetFormat(filename string) (string, bool) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not open file %s: %s", filename, err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		log.Fatalf("file %s does not contain a single line", filename)
	}
	if err != nil {
		log.Fatalf("could not read from file %s: %s", filename, err)
	}

	tsv := strings.Split(scanner.Text(), "\t")
	csv := strings.Split(scanner.Text(), ",")

	has_lineage := slices.Contains(csv, "lineage") && slices.Contains(csv, "taxid_lineage")
	if ((tsv[0] == "C") || (tsv[0] == "U")) && len(tsv) == 5 {
		return "kraken2", has_lineage
	}

	if (len(tsv) == 6) && (tsv[3] == "U") && (tsv[5] == "unclassified") {
		return "report", has_lineage
	}

	if len(csv) >= 7 {
		if slices.Compare(csv[0:7], bracken_header) == 0 {
			return "bracken", has_lineage
		}
		if slices.Compare(csv[1:7], bracken_header[0:6]) == 0 && csv[0] == "sample_id" {
			return "bracken-merged", has_lineage
		}
	}
	if len(csv) >= 5 {
		if slices.Compare(csv[0:5], mappings_header) == 0 {
			return "mapping", has_lineage
		}
	}

	return "", has_lineage
}
