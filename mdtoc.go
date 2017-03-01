package mdtoc

import "io"

func Generate(input io.Reader, output io.Writer) {
	io.Copy(output, input)
}
