package poslog

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

func newLogBlock(timestamp, host, port, user, database, pid, messageType, duration string) *LogBlock {
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

	return logBlk
}
