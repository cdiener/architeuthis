package lib

import (
	"os"
	"path/filepath"
	"testing"
)

var simple = filepath.Join("..", "testdata", "test.k2")
var with_sample = filepath.Join("..", "testdata", "S_positive_1.b2")

func TestSimple(t *testing.T) {
	lines, _ := CountLines(simple)
	out, err := os.CreateTemp("", "merged.*.b2")
	if err != nil {
		t.Fatalf("Could not create temporary file. %q", err)
	}
	defer os.Remove(out.Name())

	err = SimpleAppend([]string{simple, simple}, out.Name(), false)
	if err != nil {
		t.Fatal("Simple merge failed.")
	}
	merged_lines, _ := CountLines(out.Name())
	if merged_lines != 2*lines {
		t.Errorf("Input files had %d lines but merged file had %d.", lines, merged_lines)
	}
}

func TestComplex(t *testing.T) {
	lines, _ := CountLines(with_sample)
	out, err := os.CreateTemp("", "merged.*.b2")
	if err != nil {
		t.Fatal("Could not create temporary file.")
	}
	defer os.Remove(out.Name())

	err = SampleAppend([]string{with_sample, with_sample}, out.Name(), '\t')
	if err != nil {
		t.Fatal("Complex merge failed.")
	}
	merged_lines, _ := CountLines(out.Name())
	if merged_lines != 2*lines-1 {
		t.Errorf("Input files had %d lines but merged file had %d.", lines, merged_lines)
	}
}
