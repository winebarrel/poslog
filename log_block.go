package poslog

import (
	"strings"
)

type LogBlock struct {
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

func newLogBlockAndStmtBuilder(timestamp, host, port, user, database, pid, messageType, duration, stmt string) (*LogBlock, *strings.Builder) {
	logBlk := &LogBlock{
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

	return logBlk, stmtBldr
}
