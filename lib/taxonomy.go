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
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type Lineage struct {
	Names  []string
	Taxids []string
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
		"--show-lineage-taxids", "--add-prefix"}
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
		results[entries[0]] = &Lineage{
			strings.Split(entries[1], ";"),
			strings.Split(entries[2], ";"),
		}
	}

	return results
}

func GetRanks(format string) []string {
	re := regexp.MustCompile(`{(\w)}`)
	var r []string
	for _, match := range re.FindAllStringSubmatch(format, -1) {
		if match[1] == "" {
			log.Fatalf("Incorrect format term %s.", match[0])
		}
		r = append(r, match[1])
	}
	return r
}

func GetLeaf(lin *Lineage) (int, string) {
	leaf := ""
	idx := -1
	for i, s := range lin.Names {
		if len(s) <= 3 {
			break
		}
		idx = i
		leaf = s
	}
	return idx, leaf
}
