package mdtoc_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/katcipis/mdtoc"
)

func TestTOC(t *testing.T) {
	type testCase struct {
		input  string
		output string
	}

	testCases := map[string]testCase{
		"noHeaders": {
			input:  "lala\nbaba",
			output: "lala\nbaba",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			mdtoc.Generate(
				strings.NewReader(testCase.input),
				&output,
			)

			got := output.String()
			if testCase.output != got {
				t.Fatalf(
					"=== expected:\n%s\n=== got: %s\n",
					testCase.output,
					got,
				)
			}
		})
	}
}
