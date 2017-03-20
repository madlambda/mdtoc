package mdtoc_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/katcipis/mdtoc"
)

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
}

func TestTOC(t *testing.T) {

	testCases := []string{
		"noheaders",
		"empty",
		"onlyspaces",
		"atx/headerfirst",
		"atx/headerlast",
		"atx/trimrightspace",
		"atx/leftspacenotheader",
		"atx/invalidheader",
		"atx/onlyheader",
		"atx/headerwithhash",
		"atx/headerambiguity",
		"atx/onlyheaderandspaces",
		"atx/multiplelevels",
		"atx/specialchars",
		"atx/withspaces",
	}

	for _, name := range testCases {
		t.Run(name, func(t *testing.T) {

			input, err := ioutil.ReadFile(
				"testdata/" + name + "/input.md",
			)
			assertNoErr(t, err)
			wantRaw, err := ioutil.ReadFile(
				"testdata/" + name + "/output.md",
			)
			assertNoErr(t, err)
			want := string(wantRaw)

			var output bytes.Buffer
			err = mdtoc.Generate(
				bytes.NewBuffer(input),
				&output,
			)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := output.String()
			if want != got {
				t.Fatalf(
					"\nexpected:\n%q\ngot:\n%q\n\n",
					want,
					got,
				)
			}
		})
		t.Run(name+"/FromFile", func(t *testing.T) {
			inputfilepath := "testdata/" + name + "/input.md"
			wantRaw, err := ioutil.ReadFile(
				"testdata/" + name + "/output.md",
			)
			assertNoErr(t, err)
			want := string(wantRaw)

			var output bytes.Buffer
			err = mdtoc.GenerateFromFile(
				inputfilepath,
				&output,
			)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := output.String()
			if want != got {
				t.Fatalf(
					"\nexpected:\n%q\ngot:\n%q\n\n",
					want,
					got,
				)
			}
		})
	}
}

func TestGenerateFromInvalidFile(t *testing.T) {
	var output bytes.Buffer
	err := mdtoc.GenerateFromFile(
		"notvalid.haha.xt",
		&output,
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
