package mappings

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type taxon struct {
	reads   int
	classes map[string]int
}

type mapping map[string]taxon

func Summarize(filepath string) (mapping, error) {
	k2file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	reads := 0
	k2map := make(mapping)
	scanner := bufio.NewScanner(k2file)
	for scanner.Scan() {
		ParseMapping(&k2map, scanner.Text())
		reads += 1
	}
}

func ParseMapping(k2map *mapping, line string) error {
	tokens := strings.Split(line, "\t")

	for _, s := range strings.Split(tokens[4], " ") {
		splits := strings.Split(s, ":")
		if splits[0] == "|" {
			continue
		}
		taxid := splits[0]
		count, err := strconv.Atoi(splits[1])
		if err != nil {
			log.Fatal("Could not parse the k-mer count!")
		}

	}
}

func UpdateMapping(k2map *mapping, read_class string, kmer_taxid string, count int) {

}
