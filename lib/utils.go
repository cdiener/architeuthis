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
	"bytes"
	"io"
	"log"
	"os"
)

// Count the number of lines in a file
func CountLines(path string) (int, error) {
	buf := make([]byte, 64*1024)
	sep := []byte{'\n'}
	reader, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not open the file %s.", path)
	}

	count := 0
	for {
		nbytes, err := reader.Read(buf)
		count += bytes.Count(buf[:nbytes], sep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
