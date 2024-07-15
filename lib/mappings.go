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
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Taxon struct {
	Lineage string
	Reads   int
	Classes map[string]int
}

type ReadScore struct {
	TaxonName    string
	TaxonID      uint32
	ID           string
	Kmers        uint32
	Consistency  float64
	Confidence   float64
	Multiplicity uint32
	Entropy      float64
}

type Mapping map[string]*Taxon

// Summarize combines
func SummarizeKmers(filepath string, named bool) (Mapping, error) {
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
		err := ParseMapping(k2map, scanner.Text(), named)
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

func Entropy(abundances map[string]uint32) float64 {
	total := 0.0
	ent := 0.0
	for _, cn := range abundances {
		total += float64(cn)
	}
	for _, cn := range abundances {
		p := float64(cn) / total
		ent -= p * math.Log(p)
	}
	return ent
}

func Multiplicity(abundances map[string]uint32) uint32 {
	return uint32(len(abundances))
}

func Confidence(abundances map[string]uint32, classification string) float64 {
	total := 0.0
	for _, cn := range abundances {
		total += float64(cn)
	}
	var cn uint32
	cn, ok := abundances[classification]
	if !ok {
		cn = 0
	}

	return float64(cn) / total
}

func TaxID(token string, named bool) string {
	if named {
		split := strings.SplitN(token, "(taxid ", 2)
		if len(split) != 2 {
			log.Fatalf("Could not parse taxon name %s.", token)
		}
		return strings.Trim(split[1], " )")
	}

	return token
}

func ScoreRead(line string, taxondb map[string]*Lineage, named bool) *ReadScore {
	tokens := strings.Split(strings.Trim(line, " "), "\t")
	if tokens[0] != "C" {
		return nil
	}
	taxid := TaxID(tokens[2], named)
	taxid_int, err := strconv.Atoi(taxid)
	if err != nil {
		log.Fatalf("Uh-oh I thought taxon ID %s was numeric.", taxid)
	}
	lin := taxondb[taxid]
	ridx, leaf := GetLeaf(lin)
	if ridx == -1 {
		return nil
	}

	// Get classifications
	abundances := make(map[string]uint32)
	consistent := 0
	classified := 0
	var splits []string
	for _, s := range strings.Split(tokens[4], " ") {
		splits = strings.SplitN(s, ":", 2)
		if splits[0] == "|" || splits[0] == "A" {
			continue
		}
		tid, err := strconv.Atoi(splits[0])
		cn, err2 := strconv.Atoi(splits[1])
		if err != nil || err2 != nil {
			log.Fatalf("Could not parse taxon ID %s:%s.", splits[0], splits[1])
		}

		if tid > 1 {
			kmer_lin := taxondb[splits[0]]
			idx, name := GetLeaf(kmer_lin)
			if idx == -1 {
				continue
			}
			if idx > ridx {
				name = kmer_lin.Names[ridx]
			}
			if idx >= ridx {
				_, ok := abundances[name]
				if ok {
					abundances[name] += uint32(cn)
				} else {
					abundances[name] = uint32(cn)
				}
			}
			classified += cn
			if slices.Contains(lin.Names, name) {
				consistent += cn
			}
		}
	}

	score := ReadScore{
		ID:           tokens[1],
		TaxonID:      uint32(taxid_int),
		TaxonName:    leaf,
		Entropy:      Entropy(abundances),
		Multiplicity: Multiplicity(abundances),
		Kmers:        uint32(classified),
		Consistency:  float64(consistent) / float64(classified),
		Confidence:   Confidence(abundances, leaf),
	}

	return &score
}

func ScoreReadsToFile(k2path string, out string, data_dir string, format string, named bool) error {
	// Set up Kraken reader
	sample_id := strings.Split(k2path, ".")[0]
	k2file, err := os.Open(k2path)
	if err != nil {
		log.Fatalf("Could not open %s. Does this file exist?", k2path)
	}
	defer k2file.Close()

	// Set up output
	sfile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer sfile.Close()
	writer := csv.NewWriter(sfile)
	header := []string{
		"sample_id", "read_id", "taxid", "name", "rank", "n_kmers",
		"consistency", "confidence", "multiplicity", "entropy"}
	writer.Write(header)

	log.Println("Pass 1: Building the taxa database...")
	taxondb, _ := TaxonDB(k2path, data_dir, format, named)

	reads := 0
	scanner := bufio.NewScanner(k2file)

	log.Println("Pass 2: Score individuals reads...")
	log.Printf("Reading k-mer assignments from %s an dwriting to %s.", k2path, out)
	for scanner.Scan() {
		s := ScoreRead(scanner.Text(), taxondb, named)
		if s == nil {
			continue
		}
		record := []string{
			sample_id, s.ID, strconv.Itoa(int(s.TaxonID)), s.TaxonName,
			strings.Split(s.TaxonName, "__")[0], strconv.Itoa(int(s.Kmers)),
			fmt.Sprint(s.Consistency), fmt.Sprint(s.Confidence),
			strconv.Itoa(int(s.Multiplicity)), fmt.Sprint(s.Entropy),
		}
		writer.Write(record)

		reads += 1
		if reads%1e6 == 0 {
			log.Printf("Processed %d reads...", reads)
		}
	}
	log.Printf("Processing %d reads - Done.", reads)
	writer.Flush()

	return nil
}

func FilterReads(k2path string, out string, data_dir string,
	format string, named bool, min_consistency float64, max_entropy float64, max_multiplicity uint32) error {
	// Set up Kraken reader
	k2file, err := os.Open(k2path)
	if err != nil {
		log.Fatalf("Could not open %s. Does this file exist?", k2path)
	}
	defer k2file.Close()

	// Set up output
	sfile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer sfile.Close()
	writer := bufio.NewWriter(sfile)

	log.Println("Pass 1: Building the taxa database...")
	taxondb, _ := TaxonDB(k2path, data_dir, format, named)

	reads := 0
	passed := 0

	scanner := bufio.NewScanner(k2file)

	log.Println("Pass 2: Score individuals reads...")
	log.Printf("Reading k-mer assignments from %s and writing to %s.", k2path, out)
	for scanner.Scan() {
		s := ScoreRead(scanner.Text(), taxondb, named)
		reads += 1
		if s == nil || s.Consistency < min_consistency ||
			s.Entropy > max_entropy || s.Multiplicity > max_multiplicity {
			continue
		}
		writer.Write(scanner.Bytes())
		writer.WriteRune('\n')

		passed += 1
		if reads%1e6 == 0 {
			log.Printf("Processed %d reads...", reads)
		}
	}
	log.Printf("Processing %d reads - Done. %d/%d reads passed the filter.",
		reads, passed, reads)
	writer.Flush()

	return nil
}

func TaxonDB(filepath string, data_dir string, format string, named bool) (map[string]*Lineage, int) {
	k2file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Could not open %s. Does this file exist?", filepath)
	}
	defer k2file.Close()

	reads := 0
	scanner := bufio.NewScanner(k2file)
	taxids := make(map[string]bool, 1e4)

	log.Printf("Reading k-mer assignments from %s.", filepath)
	for scanner.Scan() {
		tokens := strings.Split(strings.Trim(scanner.Text(), " "), "\t")
		if tokens[0] != "C" {
			continue
		}
		tid := TaxID(tokens[2], named)
		taxids[tid] = true

		for _, s := range strings.Split(tokens[4], " ") {
			splits := strings.Split(s, ":")
			if splits[0] == "|" || splits[0] == "A" {
				continue
			}
			tid := splits[0]
			if tid != "0" && tid != "1" {
				taxids[tid] = true
			}
		}
		reads += 1
		if reads%1e6 == 0 {
			log.Printf("Processed %d reads...", reads)
		}
	}
	log.Printf("Processing %d reads - Done.", reads)

	lineages := AddLineage(taxids, data_dir, format)

	return lineages, reads
}

func ParseMapping(k2map Mapping, line string, named bool) error {
	tokens := strings.Split(strings.Trim(line, " "), "\t")
	tid := TaxID(tokens[2], named)
	entry, ok := k2map[tid]
	if !ok {
		entry = &Taxon{Lineage: "", Reads: 0, Classes: make(map[string]int)}

		k2map[tid] = entry
	}
	entry.Reads += 1

	for _, s := range strings.Split(tokens[4], " ") {
		splits := strings.SplitN(s, ":", 2)
		if splits[0] == "|" || splits[0] == "A" {
			continue
		}
		taxid := splits[0]
		count, err := strconv.Atoi(splits[1])
		if err != nil {
			log.Fatal("Could not parse the k-mer count!")
			return err
		}
		if taxid != "0" && taxid != "1" {
			UpdateMapping(entry, taxid, count)
		}
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
		writer.Write([]string{
			"sample_id", "classification", "lineage", "total_reads",
			"name", "rank", "kmers", "in_lineage"})
	} else {
		writer.Write([]string{
			"sample_id", "classification", "total_reads",
			"taxid", "kmers"})
	}
	var recs []string
	for class, v := range k2map {
		for taxid, n := range v.Classes {
			if has_lineage {
				match := 0
				if strings.Contains(v.Lineage, taxid) {
					match = 1
				}
				rank := strings.Split(taxid, "__")[0]
				recs = []string{
					sample_id, class, v.Lineage, strconv.Itoa(v.Reads),
					taxid, rank, strconv.Itoa(n), strconv.Itoa(match)}
			} else {
				recs = []string{
					sample_id, class, strconv.Itoa(v.Reads),
					taxid, strconv.Itoa(n)}
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
			matchRanks(ref_lineage, kmer_lineage, cn, ranks)
		}
		collapsed[taxid] = ranks
		ntaxa += 1

		if ntaxa%1e3 == 0 {
			log.Printf("Processed %d taxa...", ntaxa)
		}
	}

	return collapsed
}

func matchRanks(ref_lineage *Lineage, kmer_lineage *Lineage, count int, entry *Taxon) int {
	matched_ranks := 0
	for i, kn := range kmer_lineage.Names {
		rn := ref_lineage.Names[i]
		if len(kn) < 4 || len(rn) < 4 {
			break
		}
		matched_ranks += 1
		UpdateMapping(entry, kn, count)
	}

	return matched_ranks
}
