package lib

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Taxon struct {
	Reads   int
	Classes map[string]int
}

type Mapping map[string]*Taxon

// Summarize combines
func Summarize(filepath string) (Mapping, error) {
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
			fmt.Printf("\rProcessing reads: %d", reads)
		}
	}
	fmt.Printf("\rProcessing reads: %d - Done.\n", reads)
	return k2map, nil
}

func ParseMapping(k2map Mapping, line string) error {
	tokens := strings.Split(strings.Trim(line, " "), "\t")

	entry, ok := k2map[tokens[2]]
	if !ok {
		entry = &Taxon{Reads: 0, Classes: make(map[string]int)}
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
	mfile, err := os.Create(filepath)
	if err != nil {
		log.Fatal("Could not open file for writing!")
		return err
	}
	defer mfile.Close()
	writer := csv.NewWriter(mfile)
	writer.Write([]string{"sample_id", "classification", "n_reads", "taxid", "n_kmers"})
	for class, v := range k2map {
		for taxid, n := range v.Classes {
			recs := []string{sample_id, class, strconv.Itoa(v.Reads), taxid, strconv.Itoa(n)}
			writer.Write(recs)
		}
	}
	writer.Flush()

	return nil
}
