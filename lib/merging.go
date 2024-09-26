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
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SimpleAppend(files []string, out string, header bool) error {
	merged, err := os.Create(out)
	if err != nil {
		return err
	}
	defer merged.Close()
	writer := bufio.NewWriter(merged)

	for i, file := range files {
		fi, err := os.Open(file)
		if err != nil {
			return err
		}
		reader := bufio.NewScanner(fi)
		lines := 0

		if header {
			if i > 0 {
				reader.Scan()
			} else {
				lines = -1
			}
		}
		for reader.Scan() {
			_, err := writer.Write(reader.Bytes())
			writer.WriteByte('\n')
			if err != nil {
				return err
			}
			lines++
		}
		fi.Close()
		log.Printf("Wrote %d records from %s.", lines, file)
	}
	writer.Flush()

	return nil
}

func SampleAppend(files []string, out string, sep rune) error {
	merged, err := os.Create(out)
	if err != nil {
		return err
	}
	defer merged.Close()
	writer := csv.NewWriter(merged)

	var records []string
	field := make([]string, 1)
	var n_elems int
	for i, file := range files {
		sample_id := strings.Split(filepath.Base(file), ".")[0]
		fi, err := os.Open(file)
		if err != nil {
			return err
		}
		reader := csv.NewReader(fi)
		reader.Comma = sep
		header, err := reader.Read()
		if err != nil {
			fi.Close()
			return err
		}
		if i == 0 {
			n_elems = len(header)
			field[0] = "sample_id"
			header = append(field, header...)
			err = writer.Write(header)
			if err != nil {
				return err
			}
		}
		if (i > 0) && (len(header) != n_elems) {
			return fmt.Errorf(
				"file %s has a different format than previous files. "+
					"Are you sure all files have the same format?",
				file,
			)
		}

		lines := 0
		field[0] = sample_id
		for {
			records, err = reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			records = append(field, records...)
			err = writer.Write(records)
			if err != nil {
				return err
			}
			lines += 1
		}
		fi.Close()
		log.Printf("Wrote %d records from %s.", lines, file)
	}
	writer.Flush()

	return nil
}
