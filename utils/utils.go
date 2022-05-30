package utils

import (
	"bufio"
)

const (
	readLineBufSize = 4096
)

func ReadLine(reader *bufio.Reader) (string, error) {
	buf := make([]byte, 0, readLineBufSize)
	var err error

	for {
		line, isPrefix, e := reader.ReadLine()
		err = e

		if len(line) > 0 {
			buf = append(buf, line...)
		}

		if !isPrefix || err != nil {
			break
		}
	}

	return string(buf) + "\n", err
}
