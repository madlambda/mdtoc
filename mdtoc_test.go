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
		"atx/headerfirst",
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
			mdtoc.Generate(
				bytes.NewBuffer(input),
				&output,
			)

			got := output.String()
			if want != got {
				t.Fatalf(
					"=== expected:\n%s\n=== got: %s\n",
					want,
					got,
				)
			}
		})
	}
}
