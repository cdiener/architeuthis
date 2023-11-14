package lib

import (
	"bytes"
	"log"
	"path/filepath"
	"testing"
)

func TestKmers(t *testing.T) {
	filename := filepath.Join("..", "testdata", "test.k2")
	k2map, err := SummarizeKmers(filename)
	if err != nil {
		t.Fatal("Error when running summary.")
	}
	counts := k2map["816"].Reads
	if counts != 93 {
		t.Errorf("Expected %q but got %q.", 93, counts)
	}
}

func TestCollapse(t *testing.T) {
	filename := filepath.Join("..", "testdata", "test.k2")
	k2map, err := SummarizeKmers(filename)
	if err != nil {
		t.Fatal("Error when running summary.")
	}
	collapsed := CollapseRanks(k2map, "", "{k};{p};{c};{o};{f};{g};{s}")

	c := collapsed["816"]
	if c.Reads != 93 {
		t.Errorf("Expected %q but got %q.", 93, c.Reads)
	}

	bac, ok := c.Classes["k__Bacteria"]
	if !ok {
		t.Errorf("Expected %q in ranks got %q.", "k__Bacteria", c.Classes)
	} else if bac <= 0 {
		t.Errorf("Expected positive bacteria counts got %q.", bac)
	}
}

func BenchmarkKmers(b *testing.B) {
	var str bytes.Buffer
	log.SetOutput(&str)

	filename := filepath.Join("..", "testdata", "test.k2")
	for n := 0; n < b.N; n++ {
		SummarizeKmers(filename)
	}
}

func BenchmarkCollapse(b *testing.B) {
	var str bytes.Buffer
	log.SetOutput(&str)

	filename := filepath.Join("..", "testdata", "test.k2")
	k2map, _ := SummarizeKmers(filename)
	for n := 0; n < b.N; n++ {
		CollapseRanks(k2map, "", "{k};{p};{c};{o};{f};{g};{s}")
	}
}
