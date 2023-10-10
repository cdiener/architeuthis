package lib

import (
	"bytes"
	"log"
	"path/filepath"
	"testing"
)

func TestSummary(t *testing.T) {
	filename := filepath.Join("..", "testdata", "test.k2")
	k2map, err := Summarize(filename)
	if err != nil {
		t.Fatal("Error when running summary.")
	}
	counts := k2map["816"].Reads
	if counts != 93 {
		t.Errorf("Expected %q but got %q.", 93, counts)
	}
}

func BenchmarkSummary(b *testing.B) {
	var str bytes.Buffer
	log.SetOutput(&str)

	filename := filepath.Join("..", "testdata", "test.k2")
	for n := 0; n < b.N; n++ {
		Summarize(filename)
	}
}
