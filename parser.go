package poslog

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/percona/go-mysql/query"
	"github.com/winebarrel/poslog/utils"
)

var rePrefix = regexp.MustCompile(`(?s)^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\s+[^:]+):([^:]*):([^:]*):([^:]*):([^:]*):(.*)`)
var reLog = regexp.MustCompile(`(?s)^\s+(?:duration:\s+(\d+\.\d+)\s+ms\s+)?(?:statement|execute\s+[^:]+):(.*)`)

type Block struct {
	Timestamp   string
	Host        string
	Port        string
	User        string
	Database    string
	Pid         string
	MessageType string
	Duration    string
	Statement   string
	Fingerprint string `json:",omitempty"`
}

func newBlock(timestamp, host, port, user, database, pid, messageType, duration, stmt string) (*Block, *strings.Builder) {
	block := &Block{
		Timestamp:   timestamp,
		Host:        host,
		Port:        port,
		User:        user,
		Database:    database,
		Pid:         pid,
		MessageType: messageType,
		Duration:    duration,
	}

	stmtBldr := &strings.Builder{}
	stmtBldr.WriteString(stmt)

	return block, stmtBldr
}

func callBack(block *Block, stmtBldr *strings.Builder, fingerprint bool, cb func(block *Block)) {
	stmt := strings.TrimSpace(stmtBldr.String())
	block.Statement = stmt

	if fingerprint {
		block.Fingerprint = query.Fingerprint(strings.ReplaceAll(stmt, `"`, ""))
	}

	cb(block)
}

func Parse(file io.Reader, fingerprint bool, cb func(block *Block)) error {
	reader := bufio.NewReader(file)

	var block *Block
	var stmtBldr *strings.Builder

	for {
		rawLine, err := utils.ReadLine(reader)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		line := string(rawLine)
		line += "\n"

		if prefixMatches := rePrefix.FindStringSubmatch(line); prefixMatches != nil {
			if block != nil {
				callBack(block, stmtBldr, fingerprint, cb)
			}

			messageType := prefixMatches[5]

			if messageType != "LOG" {
				block, stmtBldr = nil, nil
				continue
			}

			duration := ""
			stmt := ""

			if logMatches := reLog.FindStringSubmatch(prefixMatches[6]); logMatches != nil {
				duration = logMatches[1]
				stmt = logMatches[2]
			}

			host := prefixMatches[2]
			port := ""

			if strings.Contains(host, "(") {
				hostPort := strings.SplitN(host, "(", 2)
				host = hostPort[0]
				port = strings.TrimRight(hostPort[1], ")")
			}

			user := prefixMatches[3]
			database := ""

			if strings.Contains(user, "@") {
				userDatabase := strings.SplitN(user, "@", 2)
				user = userDatabase[0]
				database = userDatabase[1]
			}

			block, stmtBldr = newBlock(
				prefixMatches[1], // timestamp
				host,
				port,
				user,
				database,
				prefixMatches[4], // pid
				messageType,
				duration,
				stmt,
			)
		} else if block != nil {
			stmtBldr.WriteString(line)
		}
	}

	if block != nil {
		callBack(block, stmtBldr, fingerprint, cb)
	}

	return nil
}
