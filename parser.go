package poslog

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/percona/go-mysql/query"
	"github.com/winebarrel/poslog/utils"
)

var (
	rePrefix = regexp.MustCompile(`(?s)^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\s+[^:]+):([^:]*):([^:]*):([^:]*):([^:]*):(.*)`)
	reLog    = regexp.MustCompile(`(?s)^\s+(?:duration:\s+(\d+\.\d+)\s+ms\s+)?(?:statement|execute\s+[^:]+):(.*)`)
)

type Parser struct {
	Callback    func(block *LogBlock)
	Fingerprint bool
}

func (p *Parser) Parse(file io.Reader) error {
	reader := bufio.NewReader(file)

	var logBlk *LogBlock
	var stmtBldr *strings.Builder

	for {
		line, err := utils.ReadLine(reader)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if prefixMatches := rePrefix.FindStringSubmatch(line); prefixMatches != nil {
			if logBlk != nil {
				p.process(logBlk, stmtBldr)
			}

			logBlk, stmtBldr = nil, nil
			messageType := prefixMatches[5]

			if messageType != "LOG" {
				continue
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

			duration := ""
			stmt := ""

			if logMatches := reLog.FindStringSubmatch(prefixMatches[6]); logMatches != nil {
				duration = logMatches[1]
				stmt = logMatches[2]
			}

			logBlk = newLogBlock(
				prefixMatches[1], // timestamp
				host,
				port,
				user,
				database,
				prefixMatches[4], // pid
				messageType,
				duration,
			)

			stmtBldr = &strings.Builder{}
			stmtBldr.WriteString(stmt)
		} else if logBlk != nil {
			stmtBldr.WriteString(line)
		}
	}

	if logBlk != nil {
		p.process(logBlk, stmtBldr)
	}

	return nil
}

func (p *Parser) process(logBlk *LogBlock, stmtBldr *strings.Builder) {
	stmt := strings.TrimSpace(stmtBldr.String())
	logBlk.Statement = stmt

	if p.Fingerprint {
		logBlk.Fingerprint = query.Fingerprint(strings.ReplaceAll(stmt, `"`, ""))
	}

	p.Callback(logBlk)
}
