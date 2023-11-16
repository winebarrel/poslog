package utils

import (
	"bufio"
	"bytes"
)

func ReadLine(r *bufio.Reader) (string, error) {
	var buf bytes.Buffer

	for {
		line, isPrefix, err := r.ReadLine()
		n := len(line)

		if n > 0 {
			buf.Write(line)
		}

		if !isPrefix || err != nil {
			return buf.String() + "\n", err
		}
	}
}
