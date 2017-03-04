package mdtoc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const headerFormat = "- [%s](#%s)"

func writeHeader(output io.Writer, level int, header string) {
	// TODO: Handle special characters on header
	line := fmt.Sprintf(headerFormat, header, strings.ToLower(header))
	// TODO: Handle output writing errors
	// TODO: Handle when header has trailing spaces
	output.Write([]byte(line + "\n"))
}

func parseHeader(parsed []string) (int, string) {
	var level int
	for parsed[level] == "#" {
		level += 1
	}
	return level, strings.Join(parsed[level:], "")
}

func Generate(input io.Reader, output io.Writer) error {
	headerStart := []byte("<!-- mdtocstart -->")
	tocHeader := []byte("# Table of Contents")
	headerEnd := []byte("<!-- mdtocend -->")
	scanner := bufio.NewScanner(input)
	var original bytes.Buffer
	var wroteHeader bool

	for scanner.Scan() {
		line := scanner.Text()
		_, err := original.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
		parsed := strings.Split(line, "#")
		if len(parsed) == 1 {
			continue
		}
		// TODO: HANDLE WHEN # IS ON THE MIDDLE
		// TODO: Test when a line has only #
		if !wroteHeader {
			output.Write(headerStart)
			output.Write([]byte("\n"))
			output.Write(tocHeader)
			output.Write([]byte("\n\n"))
			wroteHeader = true
		}
		level, header := parseHeader(parsed)
		writeHeader(output, level, header)

	}
	// TODO: HANDLE SCAN ERR

	if wroteHeader {
		// TODO: HANDLE ERR, WRONG BYTES WRITTEN
		output.Write(headerEnd)
		output.Write([]byte("\n\n"))
	}

	// TODO: HANDLE ERR, WRONG BYTES WRITTEN
	output.Write(original.Bytes())
	return nil
}
