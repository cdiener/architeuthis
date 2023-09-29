package cmd

import (
	"bufio"
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
	for scanner.Scan() {
		err := ParseMapping(k2map, scanner.Text())
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		reads += 1
		if reads%1e6 == 0 {
			fmt.Printf("Processed %dM reads...\n", reads/1e6)
		}
	}
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
