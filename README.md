# poslog

Parser to extract SQL from postgresql.log

[![build](https://github.com/winebarrel/poslog/actions/workflows/build.yml/badge.svg)](https://github.com/winebarrel/poslog/actions/workflows/build.yml)

## Installation

```sh
brew install winebarrel/poslog/poslog
```

## Usage

```sh
$ poslog -h
Usage of poslog:
  -fingerprint
    	Add SQL fingerprint
  -version
    	Print version and exit

$ cat postgresql.log
2022-05-30 04:59:41 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: select now();
2022-05-30 04:59:46 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: begin;
2022-05-30 04:59:48 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: insert into hello values (1);
2022-05-30 04:59:50 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: commit;
...

$ poslog postgresql.log # or `cat postgresql.log | poslog`
{"Timestamp":"2022-05-30 04:59:41 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" select now();"}
{"Timestamp":"2022-05-30 04:59:46 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" begin;"}
{"Timestamp":"2022-05-30 04:59:48 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" insert into hello values (1);"}
{"Timestamp":"2022-05-30 04:59:50 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" commit;"}
...

$ poslog -fingerprint postgresql.log | jq -r .Fingerprint
select now();
begin;
insert into hello values(?+);
commit;
...
```
