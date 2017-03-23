package mdtoc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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

	//var writeErr error
	//func writeOutput(b []byte) {
	//output.Write(b)
	//}

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

	if scanner.Err() != nil {
		return scanner.Err()
	}

	if wroteHeader {
		// TODO: HANDLE ERR, WRONG BYTES WRITTEN
		output.Write(headerEnd)
		output.Write([]byte("\n\n"))
	}

	// TODO: HANDLE ERR, WRONG BYTES WRITTEN
	_, err := output.Write(original.Bytes())
	return err
}

func GenerateFromFile(inputpath string, output io.Writer) error {
	file, err := os.Open(inputpath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("GenerateFromFile: error opening file: %s", err)
	}
	return Generate(file, output)
}

func GenerateInPlace(inputpath string) error {
	var output bytes.Buffer
	err := GenerateFromFile(inputpath, &output)
	if err != nil {
		return err
	}

	file, err := os.Create(inputpath)
	defer file.Close()
	if err != nil {
		// TODO: That is why we need a backup file for the original one :-)
		return fmt.Errorf("GenerateInPlace: unable to truncate file: %s", err)
	}

	expectedwrite := int64(output.Len())
	written, err := io.Copy(file, &output)
	if err != nil {
		// TODO: That is why we need a backup file for the original one :-)
		return fmt.Errorf("GenerateInPlace: unable to copy contents: %s", err)
	}
	if written != expectedwrite {
		// TODO: That is why we need a backup file for the original one :-)
		return fmt.Errorf(
			"GenerateInPlace: unable to copy contents: wrote %d expected %d",
			written,
			expectedwrite,
		)
	}
	return nil
}
