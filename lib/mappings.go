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

package lib

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Taxon struct {
	Lineage string
	Reads   int
	Classes map[string]int
}

type Mapping map[string]*Taxon

// Summarize combines
func SummarizeKmers(filepath string) (Mapping, error) {
	k2file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer k2file.Close()

	reads := 0
	k2map := make(Mapping)
	scanner := bufio.NewScanner(k2file)
	log.Printf("Reading k-mer assignments from %s.", filepath)
	for scanner.Scan() {
		err := ParseMapping(k2map, scanner.Text())
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		reads += 1
		if reads%1e6 == 0 {
			log.Printf("Processed %d reads...", reads)
		}
	}
	log.Printf("Processing %d reads - Done.", reads)
	return k2map, nil
}

func ParseMapping(k2map Mapping, line string) error {
	tokens := strings.Split(strings.Trim(line, " "), "\t")

	entry, ok := k2map[tokens[2]]
	if !ok {
		entry = &Taxon{Lineage: "", Reads: 0, Classes: make(map[string]int)}
		k2map[tokens[2]] = entry
	}
	entry.Reads += 1

	for _, s := range strings.Split(tokens[4], " ") {
		splits := strings.Split(s, ":")
		if splits[0] == "|" {
			continue
		}
		taxid := splits[0]
		count, err := strconv.Atoi(splits[1])
		if err != nil {
			log.Fatal("Could not parse the k-mer count!")
			return err
		}
		UpdateMapping(entry, taxid, count)
	}
	return nil
}

func UpdateMapping(entry *Taxon, kmer_taxid string, count int) {
	_, ok := entry.Classes[kmer_taxid]
	if !ok {
		entry.Classes[kmer_taxid] = count
	} else {
		entry.Classes[kmer_taxid] += count
	}
}

func SaveMapping(k2map Mapping, filepath string, sample_id string) error {
	var has_lineage bool
	for k := range k2map {
		has_lineage = (k2map[k].Lineage != "")
		break
	}
	mfile, err := os.Create(filepath)
	if err != nil {
		log.Fatal("Could not open file for writing!")
		return err
	}
	defer mfile.Close()
	writer := csv.NewWriter(mfile)
	if has_lineage {
		writer.Write([]string{"sample_id", "classification", "lineage", "n_reads", "taxid", "n_kmers"})
	} else {
		writer.Write([]string{"sample_id", "classification", "n_reads", "taxid", "n_kmers"})
	}
	var recs []string
	for class, v := range k2map {
		for taxid, n := range v.Classes {
			if has_lineage {
				recs = []string{sample_id, class, v.Lineage, strconv.Itoa(v.Reads), taxid, strconv.Itoa(n)}
			} else {
				recs = []string{sample_id, class, strconv.Itoa(v.Reads), taxid, strconv.Itoa(n)}
			}
			writer.Write(recs)
		}
	}
	writer.Flush()

	return nil
}

func CollapseRanks(k2map Mapping, data_dir string, format string) Mapping {
	taxa := make(map[string]bool, 100)
	for taxid, entry := range k2map {
		taxa[taxid] = true
		for k := range entry.Classes {
			taxa[k] = true
		}
	}
	lineage := AddLineage(taxa, data_dir, format)
	log.Printf("Got taxonomy for %d unique taxa. Collapsing on ranks.", len(k2map))

	collapsed := make(Mapping, 100)
	ntaxa := 0

	for taxid, entry := range k2map {
		ranks := &Taxon{
			Lineage: strings.Join(lineage[taxid].Names, ";"),
			Reads:   entry.Reads,
			Classes: make(map[string]int, 6)}
		for cl, cn := range entry.Classes {
			ref_lineage := lineage[taxid]
			kmer_lineage := lineage[cl]
			collapsed[taxid] = ranks
			matchRanks(ref_lineage, kmer_lineage, cn, ranks)
		}
		ntaxa += 1

		if ntaxa%1e3 == 0 {
			log.Printf("Processed %d taxa...", ntaxa)
		}
	}

	return collapsed
}

func matchRanks(ref_lineage *Lineage, kmer_lineage *Lineage, count int, entry *Taxon) int {
	matched_ranks := 0
	emptyTaxon, _ := regexp.Compile(`[a-z]__$`)
	for ik, k := range kmer_lineage.Taxids {
		if emptyTaxon.MatchString(kmer_lineage.Names[ik]) {
			break
		}
		i := slices.Index(ref_lineage.Taxids, k)
		if i >= 0 {
			ni := ref_lineage.Names[i]
			if !emptyTaxon.MatchString(ni) {
				matched_ranks += 1
				UpdateMapping(entry, ni, count)
			}
		} else {
			break
		}
	}

	return matched_ranks
}
