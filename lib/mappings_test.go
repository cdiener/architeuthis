package lib

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var taxondb map[string]*Lineage
var lines []string

func init() {
	filename := filepath.Join("..", "testdata", "test.k2")
	taxondb, _ = TaxonDB(filename, "", "{K};{p};{c};{o};{f};{g};{s}", false)

	lines = make([]string, 100)
	k2file, _ := os.Open(filename)
	scanner := bufio.NewScanner(k2file)
	for i := 0; i < 100; i++ {
		scanner.Scan()
		lines[i] = scanner.Text()
	}
	k2file.Close()
}

func TestKmers(t *testing.T) {
	filename := filepath.Join("..", "testdata", "test.k2")
	k2map, err := SummarizeKmers(filename, false)
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
	k2map, err := SummarizeKmers(filename, false)
	if err != nil {
		t.Fatal("Error when running summary.")
	}
	collapsed := CollapseRanks(k2map, "", "{K};{p};{c};{o};{f};{g};{s}")

	c := collapsed["816"]
	if c.Reads != 93 {
		t.Errorf("Expected %q but got %q.", 93, c.Reads)
	}

	bac, ok := c.Classes["g__Bacteroides"]
	if !ok {
		t.Errorf("Expected %q in ranks got %q.", "g__Bacteroides", c.Classes)
	} else if bac <= 0 {
		t.Errorf("Expected positive bacteria counts got %q.", bac)
	}
}

func TestScoring(t *testing.T) {
	if len(lines) < 10 {
		t.Error("Not initialized.")
	}
	for i := 0; i < 10; i++ {
		score := ScoreRead(lines[i], taxondb, false)
		if score.Consistency > 1 || score.Consistency < 0 {
			t.Errorf("Got invalid consistency score: %f", score.Consistency)
		}

		if score.Entropy < 0 {
			t.Errorf("Got invalid entropy score: %f", score.Entropy)
		}

		if score.Confidence > 1 || score.Confidence < 0.1 {
			t.Errorf("Got invalid confidence score: %f", score.Consistency)
		}
	}
}

func BenchmarkScoring(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ScoreRead(lines[n%100], taxondb, false)
	}
}

func BenchmarkKmers(b *testing.B) {
	var str bytes.Buffer
	log.SetOutput(&str)

	filename := filepath.Join("..", "testdata", "test.k2")
	for n := 0; n < b.N; n++ {
		SummarizeKmers(filename, false)
	}
}

func BenchmarkCollapse(b *testing.B) {
	var str bytes.Buffer
	log.SetOutput(&str)

	filename := filepath.Join("..", "testdata", "test.k2")
	k2map, _ := SummarizeKmers(filename, false)
	for n := 0; n < b.N; n++ {
		CollapseRanks(k2map, "", "{k};{p};{c};{o};{f};{g};{s}")
	}
}
