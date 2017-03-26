package mdtoc_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
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
		"removeheaders",
		"empty",
		"onlyspaces",
		"trimheaders",
		"atx/ignorecode",
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
		"atx/updateheader",
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
			assertNoErr(t, err)

			got := output.String()
			if want != got {
				t.Fatalf(
					"\nexpected:\n%q\ngot:\n%q\n\n",
					want,
					got,
				)
			}
		})
		t.Run(name+"/InPlace", func(t *testing.T) {
			inputdata, err := ioutil.ReadFile(
				"testdata/" + name + "/input.md",
			)
			assertNoErr(t, err)

			tmpfile, err := ioutil.TempFile("", "mdtoc.inplace.test")
			assertNoErr(t, err)
			defer func() {
				err := os.Remove(tmpfile.Name())
				assertNoErr(t, err)
			}()

			n, err := tmpfile.Write(inputdata)

			assertNoErr(t, err)
			assertNoErr(t, tmpfile.Close())

			if n != len(inputdata) {
				t.Fatalf(
					"expected to write %d but wrote %d",
					n,
					inputdata,
				)
			}

			wantRaw, err := ioutil.ReadFile(
				"testdata/" + name + "/output.md",
			)

			assertNoErr(t, err)
			want := string(wantRaw)

			err = mdtoc.GenerateInPlace(tmpfile.Name())
			assertNoErr(t, err)

			gotRaw, err := ioutil.ReadFile(tmpfile.Name())
			assertNoErr(t, err)
			got := string(gotRaw)

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

type fakeReadWriter struct {
	t                  *testing.T
	writeCalls         int
	explodeSecondWrite bool
}

func (f *fakeReadWriter) Read(b []byte) (int, error) {
	return 0, errors.New("injected error")
}

func (f *fakeReadWriter) Write(p []byte) (n int, err error) {
	if f.writeCalls > 0 {
		if f.explodeSecondWrite {
			f.t.Fatal("should not call write after error")
		}
	}
	f.writeCalls += 1
	return 0, errors.New("injected error")
}

func TestInputIoError(t *testing.T) {
	var output bytes.Buffer
	err := mdtoc.Generate(&fakeReadWriter{}, &output)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if output.Len() != 0 {
		t.Fatalf("should have not written to output, wrote[%d]", output.Len())
	}
}

func TestOnCorruptedHeaderFails(t *testing.T) {
	const corruptedheader = `
		<!-- mdtocstart -->
		# Table of Contents
		- [Header](#header)
		# Header
	`
	input := strings.NewReader(corruptedheader)
	var output bytes.Buffer
	err := mdtoc.Generate(input, &output)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if output.Len() != 0 {
		t.Fatalf("expected empty output, got: %q", output)
	}
}

func TestOutputIoError(t *testing.T) {
	input := strings.NewReader("whatever")
	err := mdtoc.Generate(input, &fakeReadWriter{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestStopWritingOnOutputError(t *testing.T) {
	input := strings.NewReader("# Header")
	err := mdtoc.Generate(input, &fakeReadWriter{
		t:                  t,
		explodeSecondWrite: true,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGenerateFromInvalidFile(t *testing.T) {
	var output bytes.Buffer
	err := mdtoc.GenerateFromFile(
		"notvalid.haha.xt.777",
		&output,
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGenerateInPlaceInvalidFile(t *testing.T) {
	invalidFile := "notvalid.haha.xt.666"
	defer os.Remove(invalidFile)

	err := mdtoc.GenerateInPlace(invalidFile)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
