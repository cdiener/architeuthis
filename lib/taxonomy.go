package lib

import (
	"log"
	"os/exec"
	"strings"
)

type Lineage struct {
	Names  string
	Taxids string
}

func HasTaxonkit() (string, bool) {
	cmd := exec.Command("taxonkit", "version")
	var out strings.Builder
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", false
	}
	version := strings.Split(out.String(), " v")[1]

	return strings.Trim(version, "\r\n"), true
}

func AddLineage[K any](taxids map[string]K, data_dir string, format string) map[string]*Lineage {
	args := []string{"reformat", "--taxid-field", "1",
		"--show-lineage-taxids", "--add-prefix", "--trim"}
	args = append(args, "--format", format)
	if data_dir != "" {
		args = append(args, "--data-dir", data_dir)
	}
	keys := make([]string, len(taxids))
	i := 0
	for k := range taxids {
		keys[i] = k
		i++
	}

	cmd := exec.Command("taxonkit", args...)
	cmd.Stdin = strings.NewReader(strings.Join(keys, "\n"))

	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	results := make(map[string]*Lineage, len(taxids))
	for _, line := range strings.Split(strings.Trim(out.String(), "\r\n"), "\n") {
		entries := strings.Split(line, "\t")
		results[entries[0]] = &Lineage{entries[1], entries[2]}
	}

	return results
}
