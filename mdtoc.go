package mdtoc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const headerFormat = "- [%s](#%s)"
const atxHeader = "#"
const headerIdent = "    "

func isValidHeaderRune(r rune) bool {
	return unicode.IsNumber(r) || unicode.IsLetter(r) || unicode.IsSpace(r)
}

func normalizeHeader(header string) string {
	lowerNoHash := strings.TrimLeft(strings.ToLower(header), "#")
	noInvalidChars := []rune{}

	for _, r := range lowerNoHash {
		if isValidHeaderRune(r) {
			noInvalidChars = append(noInvalidChars, r)
		}
	}

	return strings.Replace(string(noInvalidChars), " ", "-", -1)
}

func writeHeader(
	output io.Writer,
	level int,
	header string,
	headersCount map[string]int,
) {
	normalizedHeader := normalizeHeader(header)
	count := headersCount[normalizedHeader]
	headersCount[normalizedHeader] = count + 1
	if count > 0 {
		normalizedHeader = fmt.Sprintf("%s-%d", normalizedHeader, count)
	}
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

func parseHeader(line string) (int, string, bool) {
	if !startsWithAtxHeader(line) {
		return 0, "", false
	}
	spaceTrimmed := strings.TrimRight(line, " ")
	parsed := strings.Split(spaceTrimmed, " ")
	if len(parsed) == 1 {
		return 0, "", false
	}
	headerlevel := len(parsed[0])
	header := parsed[1:]
	return headerlevel, strings.Join(header, " "), true
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
		level, header, ok := parseHeader(line)
		if !ok {
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
