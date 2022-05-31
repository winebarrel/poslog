package poslog

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/percona/go-mysql/query"
	"github.com/winebarrel/poslog/utils"
)

var (
	reHeader         = regexp.MustCompile(`(?s)^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\s+[^:]+):([^:]*):([^:]*):([^:]*):([^:]*):(.*)`)
	reLog            = regexp.MustCompile(`(?s)^\s+(?:duration:\s+(\d+\.\d+)\s+ms\s+)?(?:statement|execute\s+[^:]+):(.*)`)
	reParams         = regexp.MustCompile(`(?s)^\s+parameters:\s+(.*)`)
	reParamsSplitter = regexp.MustCompile(`(?m), \$\d+ = `)
)

type Parser struct {
	Callback    func(block *LogBlock)
	Fingerprint bool
	FillParams  bool
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

		if hdrMatches := reHeader.FindStringSubmatch(line); hdrMatches != nil {
			messageType := hdrMatches[5]

			if logBlk != nil && messageType == "DETAIL" {
				if paramsMatches := reParams.FindStringSubmatch(hdrMatches[6]); paramsMatches != nil {
					logBlk.Params = parseParameters(paramsMatches[1])
				}
			}

			if logBlk != nil {
				p.process(logBlk, stmtBldr)
			}

			logBlk, stmtBldr = nil, nil

			if messageType != "LOG" {
				continue
			}

			host := hdrMatches[2]
			port := ""

			if strings.Contains(host, "(") {
				hostPort := strings.SplitN(host, "(", 2)
				host = hostPort[0]
				port = strings.TrimRight(hostPort[1], ")")
			}

			user := hdrMatches[3]
			database := ""

			if strings.Contains(user, "@") {
				userDatabase := strings.SplitN(user, "@", 2)
				user = userDatabase[0]
				database = userDatabase[1]
			}

			duration := ""
			stmt := ""

			if logMatches := reLog.FindStringSubmatch(hdrMatches[6]); logMatches != nil {
				duration = logMatches[1]
				stmt = logMatches[2]
			}

			logBlk, stmtBldr = newLogBlockAndStmtBuilder(
				hdrMatches[1], // timestamp
				host,
				port,
				user,
				database,
				hdrMatches[4], // pid
				messageType,
				duration,
				stmt,
			)
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

	if p.FillParams {
		for i := len(logBlk.Params); i > 0; i-- {
			placeholder := "$" + strconv.Itoa(i)
			v := logBlk.Params[i-1]
			logBlk.Statement = strings.ReplaceAll(logBlk.Statement, placeholder, v)
		}

		logBlk.Params = nil
	}

	if p.Fingerprint {
		logBlk.Fingerprint = query.Fingerprint(strings.ReplaceAll(stmt, `"`, ""))
	}

	p.Callback(logBlk)
}

func parseParameters(params string) []string {
	params = ", " + strings.TrimSpace(params)
	paramList := reParamsSplitter.Split(params, -1)

	if len(paramList) > 0 {
		paramList = paramList[1:]
	}

	return paramList
}
