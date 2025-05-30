package poslog_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/poslog"
)

var httpLog = strings.Join([]string{
	"* Host rdsauthproxy:1108 was resolved.\n",
	"* IPv6: (none)\n",
	"* IPv4: 10.0.0.0\n",
	"*   Trying 10.0.0.0:1108...\n",
	"* Connected to rdsauthproxy (10.0.0.0) port 1108\n",
	"* using HTTP/1.x\n",
	"> POST /authenticateRequest HTTP/1.1\r\n",
	"Host: rdsauthproxy:1108\r\n",
	"Accept: */*\r\n",
	"Content-Length: 1234\r\n",
	"Content-Type: multipart/form-data; boundary=------------------------xxx\r\n",
	"\r\n",
	"* upload completely sent off: 1234 bytes\n",
	"< HTTP/1.1 200 OK\r\n",
	"< Content-Type: text/html;charset=utf-8\r\n",
	"< Content-Length: 0\r\n",
	"< \r\n",
}, "")

func TestParse(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var buf bytes.Buffer
	fmt.Fprintln(&buf, `2022-05-30 04:59:41 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: select now();
2022-05-30 04:59:46 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: begin;
2022-05-30 04:59:48 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: insert into hello values (1);
2022-05-30 04:59:50 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: commit;`)

	logs := []*poslog.LogBlock{}

	proc := func(logBlk *poslog.LogBlock) {
		logs = append(logs, logBlk)
	}

	p := &poslog.Parser{
		Callback:    proc,
		Fingerprint: false,
		FillParams:  false,
	}

	err := p.Parse(&buf)

	require.NoError(err)
	assert.Equal([]*poslog.LogBlock{
		{
			Timestamp:       "2022-05-30 04:59:41 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "select now();",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
		{
			Timestamp:       "2022-05-30 04:59:46 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "begin;",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
		{
			Timestamp:       "2022-05-30 04:59:48 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "insert into hello values (1);",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
		{
			Timestamp:       "2022-05-30 04:59:50 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "commit;",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
	}, logs)
}

func TestParseWithFingerprint(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var buf bytes.Buffer
	fmt.Fprintln(&buf, `2022-05-30 04:59:41 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: select 1;
2022-05-30 04:59:46 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: begin;
2022-05-30 04:59:48 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: insert into hello values (1);
2022-05-30 04:59:50 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: commit;`)

	logs := []*poslog.LogBlock{}

	proc := func(logBlk *poslog.LogBlock) {
		logs = append(logs, logBlk)
	}

	p := &poslog.Parser{
		Callback:    proc,
		Fingerprint: true,
		FillParams:  false,
	}

	err := p.Parse(&buf)

	require.NoError(err)
	assert.Equal([]*poslog.LogBlock{
		{
			Timestamp:       "2022-05-30 04:59:41 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "select 1;",
			Params:          nil,
			Fingerprint:     "select ?;",
			FingerprintSHA1: "c85932052f4e365b62b58e3b785ac4938e2afe44",
		},
		{
			Timestamp:       "2022-05-30 04:59:46 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "begin;",
			Params:          nil,
			Fingerprint:     "begin;",
			FingerprintSHA1: "61e1fe7963eac26a601afac53cf6b3e63ab73842",
		},
		{
			Timestamp:       "2022-05-30 04:59:48 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "insert into hello values (1);",
			Params:          nil,
			Fingerprint:     "insert into hello values(?+);",
			FingerprintSHA1: "09a59e4f68251a63c367ca5502f5d5959a45dc04",
		},
		{
			Timestamp:       "2022-05-30 04:59:50 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "commit;",
			Params:          nil,
			Fingerprint:     "commit;",
			FingerprintSHA1: "ba9df5ac4cba4bec768f732948b4ce99b57ddaa3",
		},
	}, logs)
}

func TestParseSkipHTTP(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var buf bytes.Buffer
	fmt.Fprintln(&buf, strings.Join([]string{
		"2022-05-30 04:59:41 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: select \n\tnow();\n",
		httpLog,
		"2022-05-30 04:59:46 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: begin;\n",
		httpLog,
		"2022-05-30 04:59:48 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: insert into \n\thello values (1);\n",
		httpLog,
		"2022-05-30 04:59:50 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: commit;\n",
	}, ""))

	logs := []*poslog.LogBlock{}

	proc := func(logBlk *poslog.LogBlock) {
		logs = append(logs, logBlk)
	}

	p := &poslog.Parser{
		Callback:    proc,
		Fingerprint: false,
		FillParams:  false,
	}

	err := p.Parse(&buf)

	require.NoError(err)
	assert.Equal([]*poslog.LogBlock{
		{
			Timestamp:       "2022-05-30 04:59:41 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "select \n\tnow();",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
		{
			Timestamp:       "2022-05-30 04:59:46 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "begin;",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
		{
			Timestamp:       "2022-05-30 04:59:48 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "insert into \n\thello values (1);",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
		{
			Timestamp:       "2022-05-30 04:59:50 UTC",
			Host:            "10.0.3.147",
			Port:            "57382",
			User:            "postgres",
			Database:        "postgres",
			Pid:             "[12768]",
			MessageType:     "LOG",
			Duration:        "",
			Statement:       "commit;",
			Params:          nil,
			Fingerprint:     "",
			FingerprintSHA1: "",
		},
	}, logs)
}
