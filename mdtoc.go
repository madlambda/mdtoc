package mdtoc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const headerFormat = "- [%s](#%s)"
const atxHeader = "#"
const headerIdent = "    "

func writeHeader(
	output io.Writer,
	level int,
	header string,
	headersCount map[string]int,
) {
	normalizedHeader := strings.TrimLeft(strings.ToLower(header), "#")
	count := headersCount[normalizedHeader]
	headersCount[normalizedHeader] = count + 1
	if count > 0 {
		normalizedHeader = fmt.Sprintf("%s-%d", normalizedHeader, count)
	}
	// TODO: Handle special characters on header
	line := fmt.Sprintf(
		headerFormat,
		header,
		normalizedHeader,
	)
	// TODO: Handle output writing errors
	for i := 1; i < level; i++ {
		output.Write([]byte(headerIdent))
	}
	output.Write([]byte(line + "\n"))
}

func parseHeader(parsed []string) (int, string) {
	headerlevel := len(parsed[0])
	header := parsed[1:]
	return headerlevel, strings.Trim(strings.Join(header, ""), " ")
}

func startsWithAtxHeader(line string) bool {
	return strings.Index(line, atxHeader) == 0
}

func Generate(input io.Reader, output io.Writer) error {
	headerStart := []byte("<!-- mdtocstart -->")
	tocHeader := []byte("# Table of Contents")
	headerEnd := []byte("<!-- mdtocend -->")
	scanner := bufio.NewScanner(input)
	headersCount := map[string]int{}

	var original bytes.Buffer
	var wroteHeader bool

	for scanner.Scan() {
		line := scanner.Text()
		_, err := original.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
		if !startsWithAtxHeader(line) {
			continue
		}
		// markdown atx headers MUST have space after #
		parsed := strings.Split(line, " ")
		if len(parsed) == 1 {
			continue
		}
		level, header := parseHeader(parsed)
		if header == "" {
			continue
		}
		if !wroteHeader {
			output.Write(headerStart)
			output.Write([]byte("\n"))
			output.Write(tocHeader)
			output.Write([]byte("\n\n"))
			wroteHeader = true
		}
		writeHeader(output, level, header, headersCount)
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
